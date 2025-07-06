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
	avatarRepo AthRepo
}

func (s *AthService) Profile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, common.ErrNotFound
	}

	response := ProfileResponse{
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time,
		UpdatedAt: user.UpdatedAt.Time,
		Avatar:    user.Avatar,
		Role:      user.Role,
	}

	return &response, nil

}

func NewAuthService(repo repository.Queries, log *logger.Logger, tokenManager utils.Manager) IAuthService {
	return &AthService{
		repo:  repo,
		log:   log,
		token: tokenManager,
	}
}

// IAuthService defines authentication service methods.
type IAuthService interface {
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
	Register(ctx context.Context, req RegisterRequest) (*ProfileResponse, error)
	Profile(ctx context.Context, userID uuid.UUID) (*ProfileResponse, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*ProfileResponse, error)
}

func (s *AthService) UploadAvatar(ctx context.Context, userID uuid.UUID, data []byte) (*ProfileResponse, error) {
	// Validate file size
	const maxSize = 3 * 1024 * 1024 // 3MB
	if len(data) > maxSize {
		return nil, fmt.Errorf("avatar file too large, max 3MB allowed")
	}

	// Validate image format and aspect ratio
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("invalid image format: %w", err)
	}
	bounds := img.Bounds()
	if bounds.Dx() != bounds.Dy() {
		return nil, fmt.Errorf("avatar image must have a 1:1 aspect ratio")
	}

	// Compress image
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return nil, fmt.Errorf("failed to process image: %w", err)
	}

	// Generate filename
	filename := "avatars/" + userID.String() + ".jpg"

	// Upload file via repository
	url, err := s.avatarRepo.UploadAvatar(ctx, filename, buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("failed to upload avatar: %w", err)
	}

	// Update database
	params := repository.UpdateAvatarParams{
		ID:     userID,
		Avatar: &filename,
	}
	err = s.repo.UpdateAvatar(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update avatar in database: %w", err)
	}

	// Fetch updated profile
	profile, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user profile: %w", err)
	}

	// Prepare response
	return &ProfileResponse{
		Username:  profile.Username,
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
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Time,
			UpdatedAt: user.UpdatedAt.Time,
			Role:      user.Role,
		},
	}, nil

}
