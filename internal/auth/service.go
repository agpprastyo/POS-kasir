package auth

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AthService struct {
	Repo           repository.Store
	Log            logger.ILogger
	Token          utils.Manager
	AvatarRepo     IAthRepo
	ActivityLogger activitylog.IActivityService
}

func NewAuthService(repo repository.Store, log logger.ILogger, tokenManager utils.Manager, avaRepo IAthRepo, actLog activitylog.IActivityService) IAuthService {
	return &AthService{
		Repo:           repo,
		Log:            log,
		Token:          tokenManager,
		AvatarRepo:     avaRepo,
		ActivityLogger: actLog,
	}
}


type IAuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.ProfileResponse, error)
	Profile(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponse, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*dto.ProfileResponse, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, req dto.UpdatePasswordRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
}

func (s *AthService) Profile(ctx context.Context, userID uuid.UUID) (*dto.ProfileResponse, error) {
	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, common.ErrNotFound
	}

	if user.Avatar != nil {
		avatarURL, err := s.AvatarRepo.AvatarLink(ctx, user.ID, *user.Avatar)
		if err != nil {
			switch {
			case errors.Is(err, common.ErrAvatarNotFound):
				s.Log.Errorf("Profile | Avatar not found for user: %v", userID)
				avatarURL = ""
			case errors.Is(err, common.ErrAvatarLink):
				s.Log.Errorf("Profile | Failed to generate avatar link for user: %v", userID)
				return nil, common.ErrAvatarLink
			default:
				s.Log.Errorf("Profile | Unexpected error while generating avatar link for user: %v: %v", userID, err)
			}
		}
		user.Avatar = &avatarURL
	} else {
		user.Avatar = nil
	}

	response := dto.ProfileResponse{
		ID:        user.ID,
		IsActive:  user.IsActive,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Avatar:    user.Avatar,
		Role:      user.Role,
	}

	return &response, nil

}

func (s *AthService) UpdatePassword(ctx context.Context, userID uuid.UUID, req dto.UpdatePasswordRequest) error {
	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		s.Log.Errorf("UpdatePassword | Failed to find user by ID: %v", userID)
		return common.ErrNotFound
	}

	if !utils.CheckPassword(user.PasswordHash, req.OldPassword) {
		s.Log.Errorf("UpdatePassword | Incorrect old password for user: %v", userID)
		return common.ErrInvalidCredentials
	}

	newPassHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		s.Log.Errorf("UpdatePassword | Failed to hash new password: %v", err)
		return common.ErrInvalidCredentials
	}

	params := repository.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: newPassHash,
	}

	if err := s.Repo.UpdateUserPassword(ctx, params); err != nil {
		s.Log.Errorf("UpdatePassword | Failed to update user password: %v", err)
		return common.ErrInternal
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.Log.Warnf("UpdatePassword | Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"updated_username": user.Username,
		"updated_email":    user.Email,
		"updated_role":     user.Role,
	}

	s.ActivityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATEPASSWORD,
		repository.LogEntityTypeUSER,
		user.ID.String(),
		logDetails,
	)

	return nil
}

func (s *AthService) UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*dto.ProfileResponse, error) {

	const maxSize = 3 * 1024 * 1024
	if len(data) > maxSize {
		s.Log.Errorf("UploadAvatar | File size exceeds limit: %d bytes", len(data))
		return nil, common.ErrFileTooLarge
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		s.Log.Errorf("UploadAvatar | Failed to decode image: %v", err)
		return nil, common.ErrFileTypeNotSupported
	}
	bounds := img.Bounds()
	if bounds.Dx() != bounds.Dy() {
		s.Log.Errorf("UploadAvatar | Image is not square: %dx%d", bounds.Dx(), bounds.Dy())
		return nil, common.ErrImageNotSquare
	}

	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	if err != nil {
		s.Log.Errorf("UploadAvatar | Failed to encode image: %v", err)
		return nil, common.ErrImageProcessingFailed
	}

	filename := "avatars/" + userID.String() + ".jpg"

	url, err := s.AvatarRepo.UploadAvatar(ctx, filename, buf.Bytes())
	if err != nil {
		s.Log.Errorf("UploadAvatar | Failed to upload avatar: %v", err)
		return nil, common.ErrUploadFailed
	}

	params := repository.UpdateUserParams{
		ID:     userID,
		Avatar: &filename,
	}
	profile, err := s.Repo.UpdateUser(ctx, params)
	if err != nil {
		s.Log.Errorf("UploadAvatar | Failed to update user avatar: %v", err)
		return nil, common.ErrUploadFailed
	}

	s.Log.Infof("UploadAvatar | profile updated successfully: %v", profile.Username)

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.Log.Warnf("UploadAvatar | Actor user ID not found in context for activity logging")
	}
	logDetails := map[string]interface{}{
		"updated_username": profile.Username,
		"updated_email":    profile.Email,
		"updated_role":     profile.Role,
		"updated_avatar":   url,
	}
	s.ActivityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATEAVATAR,
		repository.LogEntityTypeUSER,
		profile.ID.String(),
		logDetails,
	)

	return &dto.ProfileResponse{
		ID:        profile.ID,
		IsActive:  profile.IsActive,
		Username:  profile.Username,
		Email:     profile.Email,
		Avatar:    &url,
		CreatedAt: profile.CreatedAt.Time,
		UpdatedAt: profile.UpdatedAt.Time,
		Role:      profile.Role,
	}, nil
}

func (s *AthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.ProfileResponse, error) {

	if _, err := s.Repo.GetUserByEmail(ctx, req.Email); err == nil {
		return nil, common.ErrEmailExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed check email existence: %w", err)
	}

	if _, err := s.Repo.GetUserByUsername(ctx, req.Username); err == nil {
		return nil, common.ErrUsernameExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("failed check username existence: %w", err)
	}

	userUUID, err := uuid.NewV7()
	if err != nil {
		s.Log.Errorf("Register | Failed to create user UUID: %v", err)
		return nil, err
	}

	passHash, err := utils.HashPassword(req.Password)
	if err != nil {
		s.Log.Errorf("Register | Failed to hash password: %v", err)
		return nil, common.ErrInvalidInput
	}

	params := repository.CreateUserParams{
		ID:           userUUID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passHash,
		Avatar:       nil,
		Role:         req.Role,
	}

	user, err := s.Repo.CreateUser(ctx, params)
	if err != nil {
		s.Log.Errorf("Register | Failed to create user: %v", err)
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, common.ErrNotFound
		case errors.Is(err, common.ErrUserExists):
			return nil, common.ErrUserExists
		default:
			return nil, common.ErrInternal
		}
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.Log.Warnf("Register | Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"created_username": user.Username,
		"created_email":    user.Email,
		"created_role":     user.Role,
	}

	s.ActivityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeCREATE,
		repository.LogEntityTypeUSER,
		user.ID.String(),
		logDetails,
	)

	return &dto.ProfileResponse{
		ID:        user.ID,
		IsActive:  user.IsActive,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Avatar:    user.Avatar,
		Role:      user.Role,
	}, nil
}

func (s *AthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.Repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.Log.Errorf("Login | Failed to find user by email 1: %v", req.Email)
			return nil, common.ErrNotFound
		default:
			s.Log.Errorf("Login | Failed to find user by email 2: %v", req.Email)
			return nil, common.ErrInvalidCredentials
		}
	}

	logDetails := map[string]interface{}{
		"login_username": user.Username,
		"login_email":    user.Email,
		"login_role":     user.Role,
	}

	pass := utils.CheckPassword(user.PasswordHash, req.Password)
	if !pass {
		s.Log.Errorf("Login | Failed to find user by email: %v", req.Email)
		s.ActivityLogger.Log(
			ctx,
			user.ID,
			repository.LogActionTypeLOGINFAILED,
			repository.LogEntityTypeUSER,
			user.ID.String(),
			logDetails,
		)
		return nil, common.ErrInvalidCredentials
	}

	token, expiredAt, err := s.Token.GenerateToken(user.Username, user.Email, user.ID, user.Role)
	if err != nil {
		s.Log.Errorf("Login | Failed to generate Token: %v", err)
		s.ActivityLogger.Log(
			ctx,
			user.ID,
			repository.LogActionTypeLOGINFAILED,
			repository.LogEntityTypeUSER,
			user.ID.String(),
			logDetails,
		)
		return nil, common.ErrInvalidCredentials
	}

	refreshToken, _, err := s.Token.GenerateRefreshToken(user.Username, user.Email, user.ID, user.Role)
	if err != nil {
		s.Log.Errorf("Login | Failed to generate Refresh Token: %v", err)
		return nil, common.ErrInternal
	}

	// Update refresh token in database (Single Session Enforcement)
	if err := s.Repo.UpdateRefreshToken(ctx, repository.UpdateRefreshTokenParams{
		ID:           user.ID,
		RefreshToken: &refreshToken,
	}); err != nil {
		s.Log.Errorf("Login | Failed to update refresh token: %v", err)
		return nil, common.ErrInternal
	}

	s.ActivityLogger.Log(
		ctx,
		user.ID,
		repository.LogActionTypeLOGINSUCCESS,
		repository.LogEntityTypeUSER,
		user.ID.String(),
		logDetails,
	)

	return &dto.LoginResponse{
		ExpiredAt:    expiredAt,
		Token:        token,
		RefreshToken: refreshToken,
		Profile: dto.ProfileResponse{
			ID:        user.ID,
			IsActive:  user.IsActive,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			Role:      user.Role,
		},
	}, nil
}

func (s *AthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	// 1. Verify token signature
	claims, err := s.Token.VerifyToken(refreshToken)
	if err != nil {
		s.Log.Errorf("RefreshToken | Invalid token: %v", err)
		return nil, common.ErrUnauthorized
	}

	if claims.Type != "refresh" {
		s.Log.Errorf("RefreshToken | Invalid token type: %v", claims.Type)
		return nil, common.ErrUnauthorized
	}

	// 2. Check token in database (Single Session Enforcement)
	user, err := s.Repo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		s.Log.Errorf("RefreshToken | User not found: %v", claims.UserID)
		return nil, common.ErrUnauthorized
	}

	if user.RefreshToken == nil || *user.RefreshToken != refreshToken {
		s.Log.Errorf("RefreshToken | Token mismatch or revoked for user: %v", claims.UserID)
		return nil, common.ErrUnauthorized // Force logout
	}

	// 3. Generate new tokens (Rotation)
	newAccessToken, newExpiredAt, err := s.Token.GenerateToken(user.Username, user.Email, user.ID, user.Role)
	if err != nil {
		s.Log.Errorf("RefreshToken | Failed to generate access token: %v", err)
		return nil, common.ErrInternal
	}

	newRefreshToken, _, err := s.Token.GenerateRefreshToken(user.Username, user.Email, user.ID, user.Role)
	if err != nil {
		s.Log.Errorf("RefreshToken | Failed to generate refresh token: %v", err)
		return nil, common.ErrInternal
	}

	// 4. Update database with new refresh token
	if err := s.Repo.UpdateRefreshToken(ctx, repository.UpdateRefreshTokenParams{
		ID:           user.ID,
		RefreshToken: &newRefreshToken,
	}); err != nil {
		s.Log.Errorf("RefreshToken | Failed to update refresh token in DB: %v", err)
		return nil, common.ErrInternal
	}

	return &dto.LoginResponse{
		ExpiredAt:    newExpiredAt,
		Token:        newAccessToken,
		RefreshToken: newRefreshToken,
		Profile: dto.ProfileResponse{
			ID:        user.ID,
			IsActive:  user.IsActive,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			Role:      user.Role,
		},
	}, nil
}
