package user_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/user"
	user_repo "POS-kasir/internal/user/repository"
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/mocks"
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

func TestDeleteUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()

		mockRepo.EXPECT().DeleteUser(ctx, userID).Return(nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeDELETE, activitylog_repo.LogEntityTypeUSER, userID.String(), gomock.Any()).Times(1)

		err := service.DeleteUser(ctx, userID)
		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()

		mockRepo.EXPECT().DeleteUser(ctx, userID).Return(pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteUser(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, common.ErrNotFound, err)
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()

		expectedErr := errors.New("db error")
		mockRepo.EXPECT().DeleteUser(ctx, userID).Return(expectedErr)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteUser(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestToggleUserStatus(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()

		mockRepo.EXPECT().ToggleUserActiveStatus(ctx, userID).Return(uuid.Nil, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATE, activitylog_repo.LogEntityTypeUSER, userID.String(), gomock.Any()).Times(1)

		err := service.ToggleUserStatus(ctx, userID)
		assert.NoError(t, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()

		mockRepo.EXPECT().ToggleUserActiveStatus(ctx, userID).Return(uuid.Nil, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		err := service.ToggleUserStatus(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, common.ErrNotFound, err)
	})

	t.Run("InternalError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()

		expectedErr := errors.New("db error")
		mockRepo.EXPECT().ToggleUserActiveStatus(ctx, userID).Return(uuid.Nil, expectedErr)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		err := service.ToggleUserStatus(ctx, userID)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("Success_WithAvatar", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()
		userID := uuid.New()
		avatarFile := "avatar.jpg"
		dbUser := user_repo.User{
			ID:           userID,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hashed",
			Role:         user_repo.UserRoleCashier,
			IsActive:     true,
			CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			Avatar:       &avatarFile,
		}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(dbUser, nil)
		mockAvatar.EXPECT().AvatarLink(ctx, userID, avatarFile).Return("http://link.to/avatar.jpg", nil)

		resp, err := service.GetUserByID(ctx, userID)
		assert.NoError(t, err)
		assert.Equal(t, userID, resp.ID)
		assert.Equal(t, "http://link.to/avatar.jpg", *resp.Avatar)
	})

	t.Run("Success_NoAvatar", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()
		userID := uuid.New()
		dbUser := user_repo.User{
			ID:           userID,
			Username:     "testuser",
			Email:        "test@example.com",
			PasswordHash: "hashed",
			Role:         user_repo.UserRoleCashier,
			IsActive:     true,
			CreatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			UpdatedAt:    pgtype.Timestamptz{Time: time.Now(), Valid: true},
			Avatar:       nil,
		}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(dbUser, nil)

		resp, err := service.GetUserByID(ctx, userID)
		assert.NoError(t, err)
		assert.Nil(t, resp.Avatar)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()
		userID := uuid.New()

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetUserByID(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrNotFound, err)
	})

	t.Run("AvatarLinkError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()
		userID := uuid.New()
		avatarFile := "avatar.jpg"
		dbUser := user_repo.User{
			ID:       userID,
			Username: "testuser",
			Avatar:   &avatarFile,
		}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(dbUser, nil)
		mockAvatar.EXPECT().AvatarLink(ctx, userID, avatarFile).Return("", errors.New("link error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetUserByID(ctx, userID)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()

		existingUser := user_repo.User{
			ID:       userID,
			Username: "olduser",
			Email:    "old@example.com",
			Role:     user_repo.UserRoleCashier,
		}
		newUsername := "newuser"
		newEmail := "new@example.com"
		newRole := user_repo.UserRoleManager
		isActive := false
		req := user.UpdateUserRequest{
			Username: &newUsername,
			Email:    &newEmail,
			Role:     &newRole,
			IsActive: &isActive,
		}
		updatedUser := user_repo.User{
			ID:       userID,
			Username: newUsername,
			Email:    newEmail,
			Role:     user_repo.UserRoleManager,
			IsActive: false,
		}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)
		mockRepo.EXPECT().CheckUserExistence(ctx, gomock.Any()).Return(user_repo.CheckUserExistenceRow{UsernameExists: false, EmailExists: false}, nil).Times(2)
		mockRepo.EXPECT().UpdateUser(ctx, gomock.Any()).Return(updatedUser, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeUPDATE, activitylog_repo.LogEntityTypeUSER, userID.String(), gomock.Any()).Times(1)

		resp, err := service.UpdateUser(ctx, userID, req)
		assert.NoError(t, err)
		assert.Equal(t, newUsername, resp.Username)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()
		req := user.UpdateUserRequest{}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(user_repo.User{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.UpdateUser(ctx, userID, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrNotFound, err)
	})

	t.Run("UsernameExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()
		existingUser := user_repo.User{ID: userID}
		username := "exists"
		req := user.UpdateUserRequest{Username: &username}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)
		mockRepo.EXPECT().CheckUserExistence(ctx, gomock.Any()).Return(user_repo.CheckUserExistenceRow{UsernameExists: true}, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.UpdateUser(ctx, userID, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrUsernameExists, err)
	})

	t.Run("EmailExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		userID := uuid.New()
		existingUser := user_repo.User{ID: userID}
		email := "exists@example.com"
		req := user.UpdateUserRequest{Email: &email}

		mockRepo.EXPECT().GetUserByID(ctx, userID).Return(existingUser, nil)
		mockRepo.EXPECT().CheckUserExistence(ctx, gomock.Any()).Return(user_repo.CheckUserExistenceRow{EmailExists: true}, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.UpdateUser(ctx, userID, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrEmailExists, err)
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()

		users := []user_repo.ListUsersRow{
			{ID: uuid.New(), Username: "u1", Email: "e1", Role: user_repo.UserRoleCashier, IsActive: true},
		}
		mockRepo.EXPECT().ListUsers(ctx, gomock.Any()).Return(users, nil)
		mockRepo.EXPECT().CountUsers(ctx, gomock.Any()).Return(int64(1), nil)

		req := user.UsersRequest{}
		resp, err := service.GetAllUsers(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Users))
	})

	t.Run("ComplexParams", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()

		page := 2
		limit := 5
		sortBy := user_repo.UserOrderColumnUsername
		sortOrder := user_repo.SortOrderDesc
		search := "test"
		role := user_repo.UserRoleManager
		isActive := true
		req := user.UsersRequest{
			Page:      &page,
			Limit:     &limit,
			SortBy:    &sortBy,
			SortOrder: &sortOrder,
			Search:    &search,
			Role:      &role,
			IsActive:  &isActive,
		}
		users := []user_repo.ListUsersRow{
			{ID: uuid.New(), Username: "u1", Email: "e1", Role: user_repo.UserRoleManager, IsActive: true},
		}

		mockRepo.EXPECT().ListUsers(ctx, gomock.Any()).Return(users, nil)
		mockRepo.EXPECT().CountUsers(ctx, gomock.Any()).Return(int64(10), nil)

		resp, err := service.GetAllUsers(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(resp.Users))
		assert.Equal(t, 2, resp.Pagination.TotalPage)
	})

	t.Run("EmptyList", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()

		mockRepo.EXPECT().ListUsers(ctx, gomock.Any()).Return([]user_repo.ListUsersRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := user.UsersRequest{}
		resp, err := service.GetAllUsers(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(resp.Users))
	})

	t.Run("ListError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.Background()

		mockRepo.EXPECT().ListUsers(ctx, gomock.Any()).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := user.UsersRequest{}
		resp, err := service.GetAllUsers(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		req := user.CreateUserRequest{
			Username: "newuser",
			Email:    "new@example.com",
			Password: "password123",
			Role:     "cashier",
		}

		mockRepo.EXPECT().CheckUserExistence(ctx, gomock.Any()).Return(user_repo.CheckUserExistenceRow{UsernameExists: false, EmailExists: false}, nil)
		newUser := user_repo.User{
			ID:       uuid.New(),
			Username: req.Username,
			Email:    req.Email,
			Role:     user_repo.UserRoleCashier,
			IsActive: true,
		}
		mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(newUser, nil)
		mockActivity.EXPECT().Log(ctx, gomock.Any(), activitylog_repo.LogActionTypeCREATE, activitylog_repo.LogEntityTypeUSER, gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.CreateUser(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Username, resp.Username)
	})

	t.Run("UsernameExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		req := user.CreateUserRequest{Username: "newuser"}

		mockRepo.EXPECT().CheckUserExistence(ctx, gomock.Any()).Return(user_repo.CheckUserExistenceRow{UsernameExists: true}, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.CreateUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrUsernameExists, err)
	})

	t.Run("EmailExists", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		req := user.CreateUserRequest{Email: "new@example.com"}

		mockRepo.EXPECT().CheckUserExistence(ctx, gomock.Any()).Return(user_repo.CheckUserExistenceRow{EmailExists: true}, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.CreateUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, common.ErrEmailExists, err)
	})

	t.Run("CreateError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockRepo := mocks.NewMockUserRepo(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockActivity := mocks.NewMockIActivityService(ctrl)
		mockAvatar := mocks.NewMockIAthRepo(ctrl)
		service := user.NewUsrService(mockRepo, mockLogger, mockActivity, mockAvatar)
		ctx := context.WithValue(context.Background(), common.UserIDKey, uuid.New())
		req := user.CreateUserRequest{Password: "pass"}

		mockRepo.EXPECT().CheckUserExistence(ctx, gomock.Any()).Return(user_repo.CheckUserExistenceRow{UsernameExists: false}, nil)
		mockRepo.EXPECT().CreateUser(ctx, gomock.Any()).Return(user_repo.User{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.CreateUser(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}
