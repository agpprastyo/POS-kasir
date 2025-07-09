package auth

import (
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
	repo       repository.Queries
	log        *logger.Logger
	token      utils.Manager
	avatarRepo IAthRepo
}

func (s *AthService) Profile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, common.ErrNotFound
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

func NewAuthService(repo repository.Queries, log *logger.Logger, tokenManager utils.Manager, avaRepo IAthRepo) IAuthService {
	return &AthService{
		repo:       repo,
		log:        log,
		token:      tokenManager,
		avatarRepo: avaRepo,
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

	pass := utils.CheckPassword(user.PasswordHash, req.Password)
	if !pass {
		s.log.Errorf("User Service | Failed to find user by email: %v", req.Email)
		return nil, common.ErrInvalidCredentials
	}

	token, expiredAt, err := s.token.GenerateToken(user.Username, user.Email, user.ID, user.Role)
	if err != nil {
		s.log.Errorf("User Service | Failed to generate token: %v", err)
		return nil, common.ErrInvalidCredentials
	}

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
