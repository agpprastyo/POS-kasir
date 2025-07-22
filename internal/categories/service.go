package categories

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"strconv"
)

type CtgService struct {
	repo            repository.Querier
	log             logger.ILogger
	activityService activitylog.IActivityService
}

type ICtgService interface {
	GetAllCategories(ctx context.Context, req ListCategoryRequest) ([]CategoryResponse, error)
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (CategoryResponse, error)
	GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error)
	UpdateCategory(ctx context.Context, id string, req CreateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id string) error
}

func NewCtgService(repo repository.Querier, log logger.ILogger, activityService activitylog.IActivityService) ICtgService {
	return &CtgService{
		repo:            repo,
		log:             log,
		activityService: activityService,
	}
}

func (s *CtgService) DeleteCategory(ctx context.Context, id string) error {
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		s.log.Errorf("Invalid category ID format", "error", err, "id", id)
		return common.ErrInvalidID
	}

	catID := int32(categoryID)
	// check existing category
	// ExistsCategory
	exists, err := s.repo.ExistsCategory(ctx, catID)
	if err != nil {
		s.log.Errorf("Failed to check if category exists", "error", err, "categoryID", catID)
	}

	if !exists {
		s.log.Warnf("Category not found", "categoryID", categoryID)
		return common.ErrCategoryNotFound
	}

	productCount, err := s.repo.CountProductsInCategory(ctx, &catID)
	if err != nil {
		s.log.Errorf("Failed to count products in category", "error", err, "categoryID", categoryID)
		return err
	}

	if productCount > 0 {
		s.log.Warnf(
			"Attempted to delete a category that is still in use",
			"categoryID", categoryID,
			"productCount", productCount,
		)
		return common.ErrCategoryInUse
	}

	err = s.repo.DeleteCategory(ctx, int32(categoryID))
	if err != nil {

		s.log.Errorf("Failed to delete category", "error", err, "categoryID", categoryID)
		return err
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"deleted_category_id": categoryID,
	}

	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeDELETE,
		repository.LogEntityTypeCATEGORY,
		id,
		logDetails,
	)

	s.log.Infof("Category deleted successfully", "categoryID", categoryID)
	return nil
}

func (s *CtgService) UpdateCategory(ctx context.Context, id string, req CreateCategoryRequest) (*CategoryResponse, error) {
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		s.log.Errorf("Invalid category ID", "error", err)
		return nil, err
	}

	params := repository.UpdateCategoryParams{
		ID:   int32(categoryID),
		Name: req.Name,
	}

	category, err := s.repo.UpdateCategory(ctx, params)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("Category not found", "id", id)
			return nil, nil
		default:
			s.log.Errorf("Failed to update category", "error", err)
			return nil, err
		}
	}

	response := &CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
	}

	return response, nil
}

func (s *CtgService) GetCategoryByID(ctx context.Context, id string) (*CategoryResponse, error) {
	categoryID, err := strconv.Atoi(id)
	if err != nil {
		s.log.Errorf("Invalid category ID", "error", err)
		return nil, err
	}

	category, err := s.repo.GetCategory(ctx, int32(categoryID))
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			s.log.Warnf("Category not found", "id", id)
			return nil, nil
		default:
			s.log.Errorf("Failed to get category by ID", "error", err)
			return nil, err
		}
	}

	response := &CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
	}

	return response, nil
}

func (s *CtgService) CreateCategory(ctx context.Context, req CreateCategoryRequest) (CategoryResponse, error) {

	category, err := s.repo.CreateCategory(ctx, req.Name)
	if err != nil {
		s.log.Errorf("Failed to create category", "error", err)
		return CategoryResponse{}, err
	}

	response := CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		CreatedAt: category.CreatedAt.Time,
		UpdatedAt: category.UpdatedAt.Time,
	}

	// Log the activity of creating a category
	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("Actor user ID not found in context for activity logging")
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

func (s *CtgService) GetAllCategories(ctx context.Context, req ListCategoryRequest) ([]CategoryResponse, error) {
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
		s.log.Errorf("Failed to get all categories", "error", err)
		return nil, err
	}

	var response []CategoryResponse
	for _, category := range categories {
		response = append(response, CategoryResponse{
			ID:        category.ID,
			Name:      category.Name,
			CreatedAt: category.CreatedAt.Time,
			UpdatedAt: category.UpdatedAt.Time,
		})
	}

	return response, nil
}
