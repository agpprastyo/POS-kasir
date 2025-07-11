package auth

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
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

// AthService is a concrete implementation of IAuthService.
type AthService struct {
	repo           repository.Querier
	log            *logger.Logger
	token          utils.Manager
	avatarRepo     IAthRepo
	activityLogger activitylog.Service
}

func (s *AthService) Profile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, common.ErrNotFound
	}

	if user.Avatar != nil {
		// Fetch avatar URL if it exists
		avatarURL, err := s.avatarRepo.AvatarLink(ctx, user.ID, *user.Avatar)
		if err != nil {
			s.log.Errorf("User Service | Failed to get avatar link: %v", err)
			return nil, fmt.Errorf("failed to get avatar link: %w", err)
		}
		user.Avatar = &avatarURL
	} else {
		user.Avatar = nil
	}

	response := ProfileResponse{
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

func NewAuthService(repo repository.Querier, log *logger.Logger, tokenManager utils.Manager, avaRepo IAthRepo, actLog activitylog.Service) IAuthService {
	return &AthService{
		repo:           repo,
		log:            log,
		token:          tokenManager,
		avatarRepo:     avaRepo,
		activityLogger: actLog,
	}
}

// IAuthService defines authentication service methods.
type IAuthService interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*ProfileResponse, error)
	Profile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*ProfileResponse, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, req UpdatePasswordRequest) error
}

func (s *AthService) UpdatePassword(ctx context.Context, userID uuid.UUID, req UpdatePasswordRequest) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		s.log.Errorf("User Service | Failed to find user by ID: %v", userID)
		return common.ErrNotFound
	}

	if !utils.CheckPassword(user.PasswordHash, req.OldPassword) {
		s.log.Errorf("User Service | Incorrect old password for user: %v", userID)
		return common.ErrInvalidCredentials
	}

	newPassHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		s.log.Errorf("User Service | Failed to hash new password: %v", err)
		return err
	}

	params := repository.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: newPassHash,
	}

	if err := s.repo.UpdateUserPassword(ctx, params); err != nil {
		s.log.Errorf("User Service | Failed to update user password: %v", err)
		return fmt.Errorf("failed to update password in database")
	}

	// Log activity
	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"updated_username": user.Username,
		"updated_email":    user.Email,
		"updated_role":     user.Role,
	}

	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATEPASSWORD,
		repository.LogEntityTypeUSER,
		user.ID.String(),
		logDetails,
	)

	s.log.Infof("User Service | Password updated successfully for user: %v", user.Username)
	return nil
}

func (s *AthService) UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*ProfileResponse, error) {
	fmt.Println("UploadAvatar 1")
	// Validate file size
	const maxSize = 3 * 1024 * 1024 // 3MB
	if len(data) > maxSize {
		return nil, fmt.Errorf("avatar file too large, max 3MB allowed")
	}

	fmt.Println("UploadAvatar 2")

	// Validate image format and aspect ratio
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid image format: %w", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != bounds.Dy() {
		return nil, fmt.Errorf("avatar image must have a 1:1 aspect ratio")
	}

	fmt.Println("UploadAvatar 3")

	// Compress image
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	fmt.Println("UploadAvatar 4")

	// Generate filename
	filename := "avatars/" + userID.String() + ".jpg"

	fmt.Println("UploadAvatar 5")

	// Upload file via repository
	url, err := s.avatarRepo.UploadAvatar(ctx, filename, buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to upload avatar: %w", err)
	}

	fmt.Println("UploadAvatar 6")

	// Update database
	params := repository.UpdateUserParams{
		ID:     userID,
		Avatar: &filename,
	}
	profile, err := s.repo.UpdateUser(ctx, params)
	if err != nil {
		s.log.Errorf("User Service | Failed to update user avatar: %v", err)
		return nil, fmt.Errorf("failed to update avatar in database")
	}

	s.log.Infof("User Service | profile updated successfully: %v", profile.Username)

	// Log activity
	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for activity logging")
	}
	logDetails := map[string]interface{}{
		"updated_username": profile.Username,
		"updated_email":    profile.Email,
		"updated_role":     profile.Role,
		"updated_avatar":   url,
	}
	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeUPDATEAVATAR,
		repository.LogEntityTypeUSER,
		profile.ID.String(),
		logDetails,
	)

	// Prepare response
	return &ProfileResponse{
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

type checkResult struct {
	exists bool
	err    error
}

func (s *AthService) Register(ctx context.Context, req RegisterRequest) (*ProfileResponse, error) {

	emailCh := make(chan checkResult, 1)
	usernameCh := make(chan checkResult, 1)

	go func() {
		_, err := s.repo.GetUserByEmail(ctx, req.Email)
		if err == nil {
			emailCh <- checkResult{exists: true}
		} else if !errors.Is(err, pgx.ErrNoRows) {
			emailCh <- checkResult{err: err}
		} else {
			emailCh <- checkResult{}
		}
	}()

	go func() {
		_, err := s.repo.GetUserByUsername(ctx, req.Username)
		if err == nil {
			usernameCh <- checkResult{exists: true}
		} else if !errors.Is(err, pgx.ErrNoRows) {
			usernameCh <- checkResult{err: err}
		} else {
			usernameCh <- checkResult{}
		}
	}()

	select {
	case res := <-emailCh:
		if res.err != nil {
			return nil, res.err
		}
		if res.exists {
			return nil, common.ErrUserExists
		}
	case res := <-usernameCh:
		if res.err != nil {
			return nil, res.err
		}
		if res.exists {
			return nil, common.ErrUserExists
		}
	}

	userUUID, err := uuid.NewV7()
	if err != nil {
		s.log.Errorf("User Service | Failed to create user UUID: %v", err)
		return nil, err
	}

	passHash, err := utils.HashPassword(req.Password)
	if err != nil {
		s.log.Errorf("User Service | Failed to hash password: %v", err)
		return nil, err
	}

	params := repository.CreateUserParams{
		ID:           userUUID,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passHash,
		Avatar:       nil,
		Role:         req.Role,
	}

	user, err := s.repo.CreateUser(ctx, params)
	if err != nil {
		s.log.Errorf("User Service | Failed to create user: %v", err)
		return nil, err
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"created_username": user.Username,
		"created_email":    user.Email,
		"created_role":     user.Role,
	}

	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeCREATE,
		repository.LogEntityTypeUSER,
		user.ID.String(),
		logDetails,
	)

	return &ProfileResponse{
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

func (s *AthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Errorf("User Service | Failed to find user by email 1: %v", req.Email)
			return nil, common.ErrNotFound
		default:
			s.log.Errorf("User Service | Failed to find user by email 2: %v", req.Email)
			return nil, common.ErrInvalidCredentials
		}
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"login_username": user.Username,
		"login_email":    user.Email,
		"login_role":     user.Role,
	}

	pass := utils.CheckPassword(user.PasswordHash, req.Password)
	if !pass {
		s.log.Errorf("User Service | Failed to find user by email: %v", req.Email)
		s.activityLogger.Log(
			ctx,
			actorID,
			repository.LogActionTypeLOGINFAILED,
			repository.LogEntityTypeUSER,
			user.ID.String(),
			logDetails,
		)
		return nil, common.ErrInvalidCredentials
	}

	token, expiredAt, err := s.token.GenerateToken(user.Username, user.Email, user.ID, user.Role)
	if err != nil {
		s.log.Errorf("User Service | Failed to generate token: %v", err)
		s.activityLogger.Log(
			ctx,
			actorID,
			repository.LogActionTypeLOGINFAILED,
			repository.LogEntityTypeUSER,
			user.ID.String(),
			logDetails,
		)
		return nil, common.ErrInvalidCredentials
	}

	s.activityLogger.Log(
		ctx,
		actorID,
		repository.LogActionTypeLOGINSUCCESS,
		repository.LogEntityTypeUSER,
		user.ID.String(),
		logDetails,
	)

	return &LoginResponse{
		ExpiredAt: expiredAt,
		Token:     token,
		Profile: ProfileResponse{
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
