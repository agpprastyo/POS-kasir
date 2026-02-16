package categories_test

import (
	"POS-kasir/internal/categories"
	categories_repo "POS-kasir/internal/categories/repository"
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/internal/common"
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

func setupTest(t *testing.T) (*mocks.MockCategoriesRepository, *mocks.MockFieldLogger, *mocks.MockIActivityService, categories.ICtgService) {
	ctrl := gomock.NewController(t)
	mockRepo := mocks.NewMockCategoriesRepository(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)
	service := categories.NewCtgService(mockRepo, mockLogger, mockActivity)
	return mockRepo, mockLogger, mockActivity, service
}

func TestCtgService_GetAllCategories(t *testing.T) {
	mockRepo, mockLogger, _, service := setupTest(t)
	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		repoCategories := []categories_repo.Category{
			{ID: 1, Name: "Food", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
			{ID: 2, Name: "Drink", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
		}

		mockRepo.EXPECT().ListCategories(ctx, categories_repo.ListCategoriesParams{Limit: 10, Offset: 0}).Return(repoCategories, nil)

		resp, err := service.GetAllCategories(ctx, categories.ListCategoryRequest{})

		assert.NoError(t, err)
		assert.Len(t, resp, 2)
		assert.Equal(t, "Food", resp[0].Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().ListCategories(ctx, gomock.Any()).Return([]categories_repo.Category{}, nil)
		mockLogger.EXPECT().Warnf(gomock.Any()).Times(1)

		resp, err := service.GetAllCategories(ctx, categories.ListCategoryRequest{})

		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
		assert.Nil(t, resp)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo.EXPECT().ListCategories(ctx, gomock.Any()).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetAllCategories(ctx, categories.ListCategoryRequest{})

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestCtgService_CreateCategory(t *testing.T) {
	mockRepo, mockLogger, mockActivity, service := setupTest(t)
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), common.UserIDKey, userID)
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		req := categories.CreateCategoryRequest{Name: "New Category"}
		repoCategory := categories_repo.Category{
			ID:        1,
			Name:      req.Name,
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		}

		mockRepo.EXPECT().CreateCategory(ctx, req.Name).Return(repoCategory, nil)
		mockActivity.EXPECT().Log(ctx, userID, activitylog_repo.LogActionTypeCREATE, activitylog_repo.LogEntityTypeCATEGORY, "1", gomock.Any())

		resp, err := service.CreateCategory(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})

	t.Run("CreateError", func(t *testing.T) {
		mockRepo.EXPECT().CreateCategory(ctx, gomock.Any()).Return(categories_repo.Category{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.CreateCategory(ctx, categories.CreateCategoryRequest{Name: "Fail"})

		assert.ErrorIs(t, err, common.ErrInternal)
		assert.Nil(t, resp)
	})

	t.Run("CreateMissingActor", func(t *testing.T) {
		ctxNoActor := context.Background()
		req := categories.CreateCategoryRequest{Name: "No Actor"}
		repoCategory := categories_repo.Category{ID: 2, Name: req.Name}

		mockRepo.EXPECT().CreateCategory(ctxNoActor, req.Name).Return(repoCategory, nil)
		mockLogger.EXPECT().Warnf(gomock.Any()).Times(1)
		mockActivity.EXPECT().Log(ctxNoActor, uuid.Nil, activitylog_repo.LogActionTypeCREATE, activitylog_repo.LogEntityTypeCATEGORY, "2", gomock.Any())

		resp, err := service.CreateCategory(ctxNoActor, req)

		assert.NoError(t, err)
		assert.Equal(t, req.Name, resp.Name)
	})
}

func TestCtgService_GetCategoryByID(t *testing.T) {
	mockRepo, mockLogger, _, service := setupTest(t)
	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		repoCategory := categories_repo.Category{ID: 1, Name: "Found", CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true}}
		mockRepo.EXPECT().GetCategory(ctx, int32(1)).Return(repoCategory, nil)

		resp, err := service.GetCategoryByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, "Found", resp.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetCategory(ctx, int32(1)).Return(categories_repo.Category{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetCategoryByID(ctx, 1)

		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
		assert.Nil(t, resp)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockRepo.EXPECT().GetCategory(ctx, int32(1)).Return(categories_repo.Category{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetCategoryByID(ctx, 1)

		assert.ErrorIs(t, err, common.ErrInternal)
		assert.Nil(t, resp)
	})
}

func TestCtgService_UpdateCategory(t *testing.T) {
	mockRepo, mockLogger, _, service := setupTest(t)
	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		req := categories.CreateCategoryRequest{Name: "Updated"}
		repoCategory := categories_repo.Category{ID: 1, Name: req.Name, CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true}}
		mockRepo.EXPECT().UpdateCategory(ctx, categories_repo.UpdateCategoryParams{ID: 1, Name: "Updated"}).Return(repoCategory, nil)

		resp, err := service.UpdateCategory(ctx, 1, req)

		assert.NoError(t, err)
		assert.Equal(t, "Updated", resp.Name)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().UpdateCategory(ctx, gomock.Any()).Return(categories_repo.Category{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.UpdateCategory(ctx, 1, categories.CreateCategoryRequest{Name: "X"})

		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
		assert.Nil(t, resp)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockRepo.EXPECT().UpdateCategory(ctx, gomock.Any()).Return(categories_repo.Category{}, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.UpdateCategory(ctx, 1, categories.CreateCategoryRequest{Name: "X"})

		assert.ErrorIs(t, err, common.ErrInternal)
		assert.Nil(t, resp)
	})
}

func TestCtgService_DeleteCategory(t *testing.T) {
	mockRepo, mockLogger, mockActivity, service := setupTest(t)
	userID := uuid.New()
	ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

	t.Run("Success", func(t *testing.T) {
		mockRepo.EXPECT().ExistsCategory(ctx, int32(1)).Return(true, nil)
		mockRepo.EXPECT().CountProductsInCategory(ctx, gomock.Any()).Return(int64(0), nil)
		mockRepo.EXPECT().DeleteCategory(ctx, int32(1)).Return(nil)
		mockActivity.EXPECT().Log(ctx, userID, activitylog_repo.LogActionTypeDELETE, activitylog_repo.LogEntityTypeCATEGORY, "1", gomock.Any())

		err := service.DeleteCategory(ctx, 1)

		assert.NoError(t, err)
	})

	t.Run("ExistsCheckError", func(t *testing.T) {
		mockRepo.EXPECT().ExistsCategory(ctx, int32(1)).Return(false, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.ErrorIs(t, err, common.ErrInternal)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().ExistsCategory(ctx, int32(1)).Return(false, nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
	})

	t.Run("CountProductsError", func(t *testing.T) {
		mockRepo.EXPECT().ExistsCategory(ctx, int32(1)).Return(true, nil)
		mockRepo.EXPECT().CountProductsInCategory(ctx, gomock.Any()).Return(int64(0), errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.ErrorIs(t, err, common.ErrInternal)
	})

	t.Run("InUse", func(t *testing.T) {
		mockRepo.EXPECT().ExistsCategory(ctx, int32(1)).Return(true, nil)
		mockRepo.EXPECT().CountProductsInCategory(ctx, gomock.Any()).Return(int64(5), nil)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.ErrorIs(t, err, common.ErrCategoryInUse)
	})

	t.Run("DeleteErrorNotFound", func(t *testing.T) {
		mockRepo.EXPECT().ExistsCategory(ctx, int32(1)).Return(true, nil)
		mockRepo.EXPECT().CountProductsInCategory(ctx, gomock.Any()).Return(int64(0), nil)
		mockRepo.EXPECT().DeleteCategory(ctx, int32(1)).Return(pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
	})

	t.Run("DeleteErrorInternal", func(t *testing.T) {
		mockRepo.EXPECT().ExistsCategory(ctx, int32(1)).Return(true, nil)
		mockRepo.EXPECT().CountProductsInCategory(ctx, gomock.Any()).Return(int64(0), nil)
		mockRepo.EXPECT().DeleteCategory(ctx, int32(1)).Return(errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		err := service.DeleteCategory(ctx, 1)

		assert.ErrorIs(t, err, common.ErrInternal)
	})

	t.Run("DeleteMissingActor", func(t *testing.T) {
		ctxNoActor := context.Background()
		mockRepo.EXPECT().ExistsCategory(ctxNoActor, int32(1)).Return(true, nil)
		mockRepo.EXPECT().CountProductsInCategory(ctxNoActor, gomock.Any()).Return(int64(0), nil)
		mockRepo.EXPECT().DeleteCategory(ctxNoActor, int32(1)).Return(nil)
		mockLogger.EXPECT().Warnf(gomock.Any()).Times(1)
		mockActivity.EXPECT().Log(ctxNoActor, uuid.Nil, activitylog_repo.LogActionTypeDELETE, activitylog_repo.LogEntityTypeCATEGORY, "1", gomock.Any())

		err := service.DeleteCategory(ctxNoActor, 1)

		assert.NoError(t, err)
	})
}

func TestCtgService_GetCategoryWithProductCount(t *testing.T) {
	mockRepo, mockLogger, _, service := setupTest(t)
	ctx := context.Background()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		repoData := []categories_repo.ListCategoriesWithProductsRow{
			{ID: 1, Name: "C1", ProductCount: 2, CreatedAt: pgtype.Timestamptz{Time: now, Valid: true}, UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true}},
		}
		mockRepo.EXPECT().ListCategoriesWithProducts(ctx, gomock.Any()).Return(repoData, nil)

		resp, err := service.GetCategoryWithProductCount(ctx)

		assert.NoError(t, err)
		assert.Len(t, *resp, 1)
		assert.Equal(t, int32(2), (*resp)[0].ProductCount)
	})

	t.Run("NoRowsError", func(t *testing.T) {
		mockRepo.EXPECT().ListCategoriesWithProducts(ctx, gomock.Any()).Return(nil, pgx.ErrNoRows)
		mockLogger.EXPECT().Warnf(gomock.Any()).Times(1)

		resp, err := service.GetCategoryWithProductCount(ctx)

		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
		assert.Nil(t, resp)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockRepo.EXPECT().ListCategoriesWithProducts(ctx, gomock.Any()).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, err := service.GetCategoryWithProductCount(ctx)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("Empty", func(t *testing.T) {
		mockRepo.EXPECT().ListCategoriesWithProducts(ctx, gomock.Any()).Return([]categories_repo.ListCategoriesWithProductsRow{}, nil)
		mockLogger.EXPECT().Warnf(gomock.Any()).Times(1)

		resp, err := service.GetCategoryWithProductCount(ctx)

		assert.ErrorIs(t, err, common.ErrCategoryNotFound)
		assert.Nil(t, resp)
	})
}
