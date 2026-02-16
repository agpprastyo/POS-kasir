package user_test

import (
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/internal/common"
	"POS-kasir/internal/user"
	user_repo "POS-kasir/internal/user/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/utils"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login_Extended(t *testing.T) {
	hashedPassword, _ := utils.HashPassword("password123")
	userObj := user_repo.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         user_repo.UserRoleCashier,
		IsActive:     true,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}
	req := user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	t.Run("GenerateTokenError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(userObj, nil)
		mockToken.EXPECT().GenerateToken(userObj.Username, userObj.Email, userObj.ID, string(userObj.Role)).Return("", time.Time{}, errors.New("token error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		mockActivity.EXPECT().Log(ctx, userObj.ID, activitylog_repo.LogActionTypeLOGINFAILED, activitylog_repo.LogEntityTypeUSER, userObj.ID.String(), gomock.Any()).Times(1)

		resp, err := service.Login(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrInvalidCredentials, err)
	})

	t.Run("GenerateRefreshTokenError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(userObj, nil)
		mockToken.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("access_token", time.Now().Add(time.Hour), nil)
		mockToken.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", time.Time{}, errors.New("refresh token error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.Login(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrInternal, err)
	})

	t.Run("UpdateRefreshTokenError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(userObj, nil)
		mockToken.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("access_token", time.Now().Add(time.Hour), nil)
		mockToken.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("refresh_token", time.Now().Add(time.Hour*24), nil)

		mockRepo.EXPECT().UpdateRefreshToken(ctx, gomock.Any()).Return(errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.Login(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrInternal, err)
	})
}

func TestAuthService_Profile_Extended(t *testing.T) {
	userID := uuid.New()
	avatarFile := "avatar.jpg"
	userObj := user_repo.User{
		ID:        userID,
		Username:  "testuser",
		Avatar:    &avatarFile,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	t.Run("AvatarNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userObj, nil)
		mockAvatar.EXPECT().AvatarLink(ctx, userID, avatarFile).Return("", common.ErrAvatarNotFound)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.Profile(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Avatar)
		assert.Equal(t, "", *resp.Avatar)
	})

	t.Run("AvatarLinkError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userObj, nil)
		mockAvatar.EXPECT().AvatarLink(ctx, userID, avatarFile).Return("", common.ErrAvatarLink)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.Profile(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrAvatarLink, err)
	})

	t.Run("UnexpectedAvatarError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userObj, nil)
		mockAvatar.EXPECT().AvatarLink(ctx, userID, avatarFile).Return("", errors.New("unexpected"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.Profile(ctx, userID)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotNil(t, resp.Avatar)
		assert.Equal(t, "", *resp.Avatar)
	})
}

func TestAuthService_Register_Extended(t *testing.T) {
	req := user.RegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password",
		Role:     "cashier",
	}

	t.Run("CheckEmailError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(user_repo.User{}, errors.New("db error"))

		resp, err := service.Register(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed check email existence")
	})

	t.Run("CheckUsernameError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(user_repo.User{}, pgx.ErrNoRows)
		mockRepo.EXPECT().GetUserByUsername(ctx, req.Username).Return(user_repo.User{}, errors.New("db error"))

		resp, err := service.Register(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "failed check username existence")
	})
}

func TestAuthService_RefreshToken_Extended(t *testing.T) {
	refreshToken := "valid_refresh_token"
	userID := uuid.New()
	claims := utils.JWTClaims{
		UserID: userID,
		Type:   "refresh",
	}
	userObj := user_repo.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		Role:         user_repo.UserRoleCashier,
		RefreshToken: &refreshToken,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	t.Run("VerifyTokenError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockToken.EXPECT().VerifyToken(refreshToken).Return(utils.JWTClaims{}, errors.New("invalid token"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.RefreshToken(ctx, refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrUnauthorized, err)
	})

	t.Run("InvalidTokenType", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		accessClaims := claims
		accessClaims.Type = "access"
		mockToken.EXPECT().VerifyToken(refreshToken).Return(accessClaims, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.RefreshToken(ctx, refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrUnauthorized, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockToken.EXPECT().VerifyToken(refreshToken).Return(claims, nil)
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.RefreshToken(ctx, refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrUnauthorized, err)
	})

	t.Run("TokenMismatch", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockToken.EXPECT().VerifyToken(refreshToken).Return(claims, nil)
		userWithDifferentToken := userObj
		diffToken := "different_token"
		userWithDifferentToken.RefreshToken = &diffToken
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userWithDifferentToken, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.RefreshToken(ctx, refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrUnauthorized, err)
	})

	t.Run("GenerateAccessTokenError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockToken.EXPECT().VerifyToken(refreshToken).Return(claims, nil)
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userObj, nil)
		mockToken.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", time.Time{}, errors.New("token error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.RefreshToken(ctx, refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrInternal, err)
	})

	t.Run("GenerateRefreshTokenError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockToken.EXPECT().VerifyToken(refreshToken).Return(claims, nil)
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userObj, nil)
		mockToken.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("new_access", time.Now(), nil)
		mockToken.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("", time.Time{}, errors.New("refresh error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.RefreshToken(ctx, refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrInternal, err)
	})

	t.Run("UpdateRefreshTokenError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockToken := mocks.NewMockManager(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
		ctx := context.Background()

		mockToken.EXPECT().VerifyToken(refreshToken).Return(claims, nil)
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userObj, nil)
		mockToken.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("new_access", time.Now(), nil)
		mockToken.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("new_refresh", time.Now().Add(time.Hour), nil)

		mockRepo.EXPECT().UpdateRefreshToken(ctx, gomock.Any()).Return(errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.RefreshToken(ctx, refreshToken)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrInternal, err)
	})
}
