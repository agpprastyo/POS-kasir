package auth

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AthService is a concrete implementation of IAuthService.
type AthService struct {
	repo  repository.Queries
	log   *logger.Logger
	token utils.Manager
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
