package categories

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CtgService struct {
	repo            repository.Querier
	log             logger.ILogger
	activityService activitylog.IActivityService
}

type ICtgService interface {
	GetAllCategories(ctx context.Context, req dto.ListCategoryRequest) ([]dto.CategoryResponse, error)
	CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	GetCategoryByID(ctx context.Context, categoryID int32) (*dto.CategoryResponse, error)
	UpdateCategory(ctx context.Context, categoryID int32, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error)
	DeleteCategory(ctx context.Context, categoryID int32) error
	GetCategoryWithProductCount(ctx context.Context) (*[]dto.CategoryWithCountResponse, error)
}

func NewCtgService(repo repository.Querier, log logger.ILogger, activityService activitylog.IActivityService) ICtgService {
	return &CtgService{
		repo:            repo,
		log:             log,
		activityService: activityService,
	}
}

func (s *CtgService) GetCategoryWithProductCount(ctx context.Context) (*[]dto.CategoryWithCountResponse, error) {
	params := repository.ListCategoriesWithProductsParams{
		Limit:  100,
		Offset: 0,
	}

	categories, err := s.repo.ListCategoriesWithProducts(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("GetCategoryWithProductCount | No categories found with product count")
			return nil, common.ErrCategoryNotFound
		default:
			s.log.Errorf("GetCategoryWithProductCount | Failed to get categories with product count: %v", err)
			return nil, err
		}
	}
	var response []dto.CategoryWithCountResponse
	for _, category := range categories {
		response = append(response, dto.CategoryWithCountResponse{
			ID:           category.ID,
			Name:         category.Name,
			ProductCount: int32(category.ProductCount),
			CreatedAt:    category.CreatedAt.Time,
			UpdatedAt:    category.UpdatedAt.Time,
		})
	}

	if len(response) == 0 {
		s.log.Warnf("GetCategoryWithProductCount | No categories found")
		return nil, common.ErrCategoryNotFound
	}

	return &response, nil
}

func (s *CtgService) DeleteCategory(ctx context.Context, categoryID int32) error {
	exists, err := s.repo.ExistsCategory(ctx, categoryID)
	if err != nil {
		s.log.Errorf("DeleteCategory | Failed to check if category exists: %v, categoryID=%d", err, categoryID)
		return common.ErrInternal
	}

	if !exists {
		s.log.Warnf("DeleteCategory | Category not found: categoryID=%d", categoryID)
		return common.ErrCategoryNotFound
	}

	productCount, err := s.repo.CountProductsInCategory(ctx, &categoryID)
	if err != nil {
		s.log.Errorf("DeleteCategory | Failed to count products in category: %v, categoryID=%d", err, categoryID)
		return common.ErrInternal
	}

	if productCount > 0 {
		s.log.Warnf(
			"DeleteCategory | Attempted to delete a category that is still in use: categoryID=%d, productCount=%d",
			categoryID,
			productCount,
		)
		return common.ErrCategoryInUse
	}

	err = s.repo.DeleteCategory(ctx, categoryID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("DeleteCategory | Category not found for deletion: categoryID=%d", categoryID)
			return common.ErrCategoryNotFound
		default:
			s.log.Errorf("DeleteCategory | Failed to delete category: %v, categoryID=%d", err, categoryID)
			return common.ErrInternal
		}
	}
	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("DeleteCategory | Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"deleted_category_id": categoryID,
	}

	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeDELETE,
		repository.LogEntityTypeCATEGORY,
		string(categoryID),
		logDetails,
	)
	return nil
}

func (s *CtgService) UpdateCategory(ctx context.Context, categoryID int32, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	params := repository.UpdateCategoryParams{
		ID:   categoryID,
		Name: req.Name,
	}

	category, err := s.repo.UpdateCategory(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("UpdateCategory | Category not found: id=%d", categoryID)
			return nil, common.ErrCategoryNotFound
		default:
			s.log.Errorf("UpdateCategory | Failed to update category: %v", err)
			return nil, common.ErrInternal
		}
	}

	response := &dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
	}

	return response, nil
}

func (s *CtgService) GetCategoryByID(ctx context.Context, categoryID int32) (*dto.CategoryResponse, error) {

	category, err := s.repo.GetCategory(ctx, categoryID)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("GetCategoryByID | Category not found: id=%d", categoryID)
			return nil, common.ErrCategoryNotFound

		default:
			s.log.Errorf("GetCategoryByID | Failed to get category by ID: %v", err)
			return nil, common.ErrInternal
		}
	}

	response := &dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
	}

	return response, nil
}

func (s *CtgService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {

	category, err := s.repo.CreateCategory(ctx, req.Name)
	if err != nil {
		s.log.Errorf("CreateCategory | Failed to create category: %v", err)
		return nil, common.ErrInternal
	}

	response := &dto.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
	}

	// Log the activity of creating a category
	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("CreateCategory | Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"category_id":   category.ID,
		"category_name": category.Name,
	}

	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeCREATE,
		repository.LogEntityTypeCATEGORY,
		strconv.FormatUint(uint64(category.ID), 10),
		logDetails,
	)

	return response, nil
}

func (s *CtgService) GetAllCategories(ctx context.Context, req dto.ListCategoryRequest) ([]dto.CategoryResponse, error) {
	limit := int32(10)
	if req.Limit > 0 {
		limit = req.Limit
	}

	offset := int32(0)
	if req.Offset > 0 {
		offset = req.Offset
	}

	params := repository.ListCategoriesParams{
		Limit:  limit,
		Offset: offset,
	}

	categories, err := s.repo.ListCategories(ctx, params)
	if err != nil {
		s.log.Errorf("GetAllCategories | Failed to get all categories: %v", err)
		return nil, err
	}

	if len(categories) == 0 {
		s.log.Warnf("GetAllCategories | No categories found")
		return nil, common.ErrCategoryNotFound
	}

	var response []dto.CategoryResponse
	for _, category := range categories {
		response = append(response, dto.CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Time,
			UpdatedAt: category.UpdatedAt.Time,
		})
	}

	return response, nil
}
