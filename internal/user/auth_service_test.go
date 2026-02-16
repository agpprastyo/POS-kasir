package user_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/user"
	user_repo "POS-kasir/internal/user/repository"
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/utils"
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockToken := mocks.NewMockManager(ctrl)
	mockAvatar := mocks.NewMockIAthRepo(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)

	service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
	ctx := context.Background()

	// Setup data
	password := "password123"
	ctx = context.Background()

	// Hash password for mock DB return
	hashedPassword, _ := utils.HashPassword(password)

	testUser := user_repo.User{
		ID:           uuid.New(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedPassword,
		Role:         user_repo.UserRoleCashier,
		IsActive:     true,
	}

	req := user.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(testUser, nil)

		mockToken.EXPECT().GenerateToken(testUser.Username, testUser.Email, testUser.ID, user_repo.UserRole(testUser.Role)).Return("access_token", time.Now().Add(1*time.Hour), nil)
		mockToken.EXPECT().GenerateRefreshToken(testUser.Username, testUser.Email, testUser.ID, user_repo.UserRole(testUser.Role)).Return("refresh_token", time.Now().Add(24*time.Hour), nil)

		mockRepo.EXPECT().UpdateRefreshToken(ctx, gomock.Any()).Return(nil)

		mockActivity.EXPECT().Log(ctx, testUser.ID, activitylog_repo.LogActionTypeLOGINSUCCESS, activitylog_repo.LogEntityTypeUSER, testUser.ID.String(), gomock.Any())

		resp, err := service.Login(ctx, req)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.Token != "access_token" {
			t.Errorf("expected access token 'access_token', got %v", resp.Token)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(user_repo.User{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		_, err := service.Login(ctx, req)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(testUser, nil)

		reqWrongPass := req
		reqWrongPass.Password = "wrongpassword"

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
		mockActivity.EXPECT().Log(ctx, testUser.ID, activitylog_repo.LogActionTypeLOGINFAILED, activitylog_repo.LogEntityTypeUSER, testUser.ID.String(), gomock.Any())

		_, err := service.Login(ctx, reqWrongPass)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAuthService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockToken := mocks.NewMockManager(ctrl)
	mockAvatar := mocks.NewMockIAthRepo(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)

	service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
	ctx := context.Background()

	req := user.RegisterRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
		Role:     user_repo.UserRoleCashier,
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(user_repo.User{}, pgx.ErrNoRows)
		mockRepo.EXPECT().GetUserByUsername(ctx, req.Username).Return(user_repo.User{}, pgx.ErrNoRows)

		newUser := user_repo.User{
			ID:        uuid.New(),
			Username:  req.Username,
			Email:     req.Email,
			Role:      user_repo.UserRole(req.Role),
			IsActive:  true,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}
		mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(newUser, nil)

		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeCREATE, activitylog_repo.LogEntityTypeUSER, gomock.Any(), gomock.Any())

		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.Register(ctx, req)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Email != req.Email {
			t.Errorf("expected email %v, got %v", req.Email, resp.Email)
		}
	})

	t.Run("EmailExists", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(user_repo.User{ID: uuid.New()}, nil) // found user

		_, err := service.Register(ctx, req)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("UsernameExists", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByEmail(ctx, req.Email).Return(user_repo.User{}, pgx.ErrNoRows)           // Email OK
		mockRepo.EXPECT().GetUserByUsername(ctx, req.Username).Return(user_repo.User{ID: uuid.New()}, nil) // Username found

		_, err := service.Register(ctx, req)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAuthService_Profile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockToken := mocks.NewMockManager(ctrl)
	mockAvatar := mocks.NewMockIAthRepo(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)

	service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
	ctx := context.Background()

	userID := uuid.New()
	baseUser := user_repo.User{
		ID:        userID,
		Username:  "user_profile",
		Email:     "user@profile.com",
		Role:      user_repo.UserRoleCashier,
		IsActive:  true,
		CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	t.Run("Success_WithAvatar", func(t *testing.T) {
		filename := "avatars/" + userID.String() + ".jpg"
		userWithAvatar := baseUser
		userWithAvatar.Avatar = &filename

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userWithAvatar, nil)
		mockAvatar.EXPECT().AvatarLink(ctx, userID, filename).Return("https://cdn.example.com/"+filename, nil)

		resp, err := service.Profile(ctx, userID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if *resp.Avatar != "https://cdn.example.com/"+filename {
			t.Errorf("expected avatar url, got %v", resp.Avatar)
		}
	})

	t.Run("Success_NoAvatar", func(t *testing.T) {
		userNoAvatar := baseUser
		userNoAvatar.Avatar = nil

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userNoAvatar, nil)

		resp, err := service.Profile(ctx, userID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.Avatar != nil {
			t.Errorf("expected nil avatar, got %v", resp.Avatar)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{}, pgx.ErrNoRows)

		_, err := service.Profile(ctx, userID)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAuthService_UploadAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockToken := mocks.NewMockManager(ctrl)
	mockAvatar := mocks.NewMockIAthRepo(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)

	service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)

	userID := uuid.New()
	ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

	createTestImage := func(width, height int) []byte {
		img := image.NewRGBA(image.Rect(0, 0, width, height))
		var buf bytes.Buffer
		_ = jpeg.Encode(&buf, img, nil)
		return buf.Bytes()
	}

	t.Run("Success", func(t *testing.T) {
		validImg := createTestImage(100, 100)
		filename := "avatars/" + userID.String() + ".jpg"
		url := "https://cdn.example.com/" + filename

		mockAvatar.EXPECT().UploadAvatar(ctx, filename, gomock.Any()).Return(url, nil)

		updatedUser := user_repo.User{
			ID:        userID,
			Username:  "updatedUser",
			Email:     "updated@example.com",
			Role:      user_repo.UserRoleCashier,
			Avatar:    &filename,
			CreatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		}

		mockRepo.EXPECT().UpdateUser(ctx, gomock.Any()).Return(updatedUser, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATEAVATAR, activitylog_repo.LogEntityTypeUSER, userID.String(), gomock.Any())

		resp, err := service.UploadAvatar(ctx, userID, validImg)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.Avatar == nil || *resp.Avatar != url {
			t.Errorf("expected avatar url %v, got %v", url, resp.Avatar)
		}
	})

	t.Run("FileTooLarge", func(t *testing.T) {
		largeData := make([]byte, 4*1024*1024)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.UploadAvatar(ctx, userID, largeData)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("InvalidImage", func(t *testing.T) {
		invalidData := []byte("not an image")
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.UploadAvatar(ctx, userID, invalidData)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("NotSquare", func(t *testing.T) {
		rectImg := createTestImage(100, 50)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.UploadAvatar(ctx, userID, rectImg)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestAuthService_UpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockToken := mocks.NewMockManager(ctrl)
	mockAvatar := mocks.NewMockIAthRepo(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)

	service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
	ctx := context.Background()
	userID := uuid.New()

	oldPassword := "oldPassword123"
	hashedOldPassword, _ := utils.HashPassword(oldPassword)

	existingUser := user_repo.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: hashedOldPassword,
		Role:         user_repo.UserRoleCashier,
	}

	req := user.UpdatePasswordRequest{
		OldPassword: oldPassword,
		NewPassword: "newPassword123",
	}

	t.Run("Success", func(t *testing.T) {

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)

		mockRepo.EXPECT().UpdateUserPassword(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, arg user_repo.UpdateUserPasswordParams) error {
			if arg.ID != userID {
				return common.ErrInternal
			}
			if !utils.CheckPassword(arg.PasswordHash, req.NewPassword) {
				return common.ErrInternal // Hash check failed
			}
			return nil
		})

		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATEPASSWORD, activitylog_repo.LogEntityTypeUSER, userID.String(), gomock.Any())
		mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()

		err := service.UpdatePassword(ctx, userID, req)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{}, pgx.ErrNoRows)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		err := service.UpdatePassword(ctx, userID, req)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrNotFound {
			t.Errorf("expected common.ErrNotFound, got %v", err)
		}
	})

	t.Run("IncorrectOldPassword", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		reqWrong := req
		reqWrong.OldPassword = "wrongPassword"

		err := service.UpdatePassword(ctx, userID, reqWrong)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrInvalidCredentials {
			t.Errorf("expected common.ErrInvalidCredentials, got %v", err)
		}
	})

	t.Run("UpdateFailure", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)
		mockRepo.EXPECT().UpdateUserPassword(ctx, gomock.Any()).Return(pgx.ErrTxClosed) // Some DB error

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		err := service.UpdatePassword(ctx, userID, req)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrInternal {
			t.Errorf("expected common.ErrInternal, got %v", err)
		}
	})
}

func TestAuthService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockToken := mocks.NewMockManager(ctrl)
	mockAvatar := mocks.NewMockIAthRepo(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)

	service := user.NewAuthService(mockRepo, mockLogger, mockToken, mockAvatar, mockActivity)
	ctx := context.Background()

	validRefreshToken := "valid_refresh_token"
	userID := uuid.New()
	u := user_repo.User{
		ID:           userID,
		Username:     "testuser",
		Email:        "test@example.com",
		Role:         user_repo.UserRoleCashier,
		RefreshToken: &validRefreshToken,
		CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
	}

	claims := utils.JWTClaims{
		UserID: userID,
		Type:   "refresh",
	}

	t.Run("Success", func(t *testing.T) {
		mockToken.EXPECT().VerifyToken(validRefreshToken).Return(claims, nil)
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(u, nil)

		newAccessToken := "new_access_token"
		newRefreshToken := "new_refresh_token"
		mockToken.EXPECT().GenerateToken(u.Username, u.Email, u.ID, user_repo.UserRole(u.Role)).
			Return(newAccessToken, time.Now().Add(time.Hour), nil)
		mockToken.EXPECT().GenerateRefreshToken(u.Username, u.Email, u.ID, user_repo.UserRole(u.Role)).
			Return(newRefreshToken, time.Now().Add(24*time.Hour), nil)

		mockRepo.EXPECT().UpdateRefreshToken(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, arg user_repo.UpdateRefreshTokenParams) error {
			if arg.ID != userID {
				return common.ErrInternal
			}
			if *arg.RefreshToken != newRefreshToken {
				return common.ErrInternal
			}
			return nil
		})

		resp, err := service.RefreshToken(ctx, validRefreshToken)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if resp.Token != newAccessToken {
			t.Errorf("expected access token %v, got %v", newAccessToken, resp.Token)
		}
		if resp.RefreshToken != newRefreshToken {
			t.Errorf("expected refresh token %v, got %v", newRefreshToken, resp.RefreshToken)
		}
	})

	t.Run("InvalidToken", func(t *testing.T) {
		invalidToken := "invalid_token"
		mockToken.EXPECT().VerifyToken(invalidToken).Return(utils.JWTClaims{}, errors.New("invalid token"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.RefreshToken(ctx, invalidToken)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrUnauthorized {
			t.Errorf("expected common.ErrUnauthorized, got %v", err)
		}
	})

	t.Run("InvalidTokenType", func(t *testing.T) {
		wrongTypeClaims := claims
		wrongTypeClaims.Type = "access"
		mockToken.EXPECT().VerifyToken(validRefreshToken).Return(wrongTypeClaims, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.RefreshToken(ctx, validRefreshToken)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrUnauthorized {
			t.Errorf("expected common.ErrUnauthorized, got %v", err)
		}
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockToken.EXPECT().VerifyToken(validRefreshToken).Return(claims, nil)
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.RefreshToken(ctx, validRefreshToken)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrUnauthorized {
			t.Errorf("expected common.ErrUnauthorized, got %v", err)
		}
	})

	t.Run("TokenMismatch", func(t *testing.T) {
		mockToken.EXPECT().VerifyToken(validRefreshToken).Return(claims, nil)

		userWithDifferentToken := u
		diffToken := "different_token"
		userWithDifferentToken.RefreshToken = &diffToken

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(userWithDifferentToken, nil)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.RefreshToken(ctx, validRefreshToken)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrUnauthorized {
			t.Errorf("expected common.ErrUnauthorized, got %v", err)
		}
	})

	t.Run("UpdateFailure", func(t *testing.T) {
		mockToken.EXPECT().VerifyToken(validRefreshToken).Return(claims, nil)
		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(u, nil)

		mockToken.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("new_access", time.Now(), nil)
		mockToken.EXPECT().GenerateRefreshToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return("new_refresh", time.Now(), nil)

		mockRepo.EXPECT().UpdateRefreshToken(ctx, gomock.Any()).Return(pgx.ErrTxClosed)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		_, err := service.RefreshToken(ctx, validRefreshToken)
		if err == nil {
			t.Error("expected error, got nil")
		}
		if err != common.ErrInternal {
			t.Errorf("expected common.ErrInternal, got %v", err)
		}
	})
}
