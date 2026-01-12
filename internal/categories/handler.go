package categories

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type ICtgHandler interface {
	GetAllCategoriesHandler(c *fiber.Ctx) error
	GetCategoryCountHandler(c *fiber.Ctx) error
	CreateCategoryHandler(c *fiber.Ctx) error
	GetCategoryByIDHandler(c *fiber.Ctx) error
	UpdateCategoryHandler(c *fiber.Ctx) error
	DeleteCategoryHandler(c *fiber.Ctx) error
}

func NewCtgHandler(service ICtgService, log logger.ILogger) ICtgHandler {
	return &CtgHandler{
		service: service,
		log:     log,
	}
}

type CtgHandler struct {
	service ICtgService
	log     logger.ILogger
}

// DeleteCategoryHandler
// @Summary Delete category by ID
// @Description Delete category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} common.SuccessResponse "Category deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid category ID format"
// @Failure 404 {object} common.ErrorResponse "Category not found"
// @Failure 409 {object} common.ErrorResponse "Category cannot be deleted because it is in use"
// @Failure 500 {object} common.ErrorResponse "Failed to delete category"
// @Router /categories/{id} [delete]
func (h *CtgHandler) DeleteCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	idStr := c.Params("id")

	categoryID, err := strconv.Atoi(idStr)
	if err != nil {
		h.log.Warnf("DeleteCategoryHandler | Invalid category ID format provided: %s", idStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid category ID format. ID must be a number.",
		})
	}

	err = h.service.DeleteCategory(ctx, int32(categoryID))
	if err != nil {
		h.log.Warnf("DeleteCategoryHandler | Failed to delete category: %v, categoryID: %d", err, categoryID)
		switch {
		case errors.Is(err, common.ErrCategoryNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "Category not found",
			})
		case errors.Is(err, common.ErrCategoryInUse):
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Category cannot be deleted because it is in use",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to delete category",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Category deleted successfully",
	})
}

// UpdateCategoryHandler
// @Summary Update category by ID
// @Description Update category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body dto.CreateCategoryRequest true "Category details"
// @Success 200 {object} common.SuccessResponse "Category deleted successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request body"
// @Failure 404 {object} common.ErrorResponse "Category not found"
// @Failure 409 {object} common.ErrorResponse "Category with this name already exists"
// @Failure 500 {object} common.ErrorResponse "Failed to update category"
// @Router /categories/{id} [put]
func (h *CtgHandler) UpdateCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category ID is required",
		})
	}

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		h.log.Warnf("UpdateCategoryHandler | Invalid category ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid category ID format. ID must be a number.",
		})
	}

	req := new(dto.CreateCategoryRequest)
	if err := c.BodyParser(req); err != nil {
		h.log.Warnf("UpdateCategoryHandler | Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category name is required",
		})
	}

	category, err := h.service.UpdateCategory(ctx, int32(categoryID), *req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrCategoryNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "Category not found",
			})
		case errors.Is(err, common.ErrCategoryExists):
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Category with this name already exists",
			})
		default:

			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to update category",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Category updated successfully",
		Data:    category,
	})
}

// GetCategoryByIDHandler
// @Summary Get category by ID
// @Description Get category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} common.SuccessResponse{data=dto.CategoryResponse} "Category retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid category ID format"
// @Failure 404 {object} common.ErrorResponse "Category not found"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve category"
// @Router /categories/{id} [get]
func (h *CtgHandler) GetCategoryByIDHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category ID is required",
		})
	}

	categoryID, err := strconv.Atoi(id)
	if err != nil {
		h.log.Warnf("GetCategoryByIDHandler | Invalid category ID format: %s", id)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid category ID format. ID must be a number.",
		})
	}

	category, err := h.service.GetCategoryByID(ctx, int32(categoryID))
	if err != nil {
		h.log.Warnf("GetCategoryByIDHandler | Failed to get category by ID: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to retrieve category",
		})
	}

	if category == nil {
		return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
			Message: "Category not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Category retrieved successfully",
		Data:    category,
	})
}

// CreateCategoryHandler
// @Summary Create a new category
// @Description Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body dto.CreateCategoryRequest true "Category details"
// @Success 201 {object} common.SuccessResponse{data=dto.CategoryResponse} "Category created successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid request body"
// @Failure 409 {object} common.ErrorResponse "Category with this name already exists"
// @Failure 500 {object} common.ErrorResponse "Failed to create category"
// @Router /categories [post]
func (h *CtgHandler) CreateCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(dto.CreateCategoryRequest)
	if err := c.BodyParser(req); err != nil {
		h.log.Warnf("CreateCategoryHandler | Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category name is required",
		})
	}

	category, err := h.service.CreateCategory(ctx, *req)
	if err != nil {
		h.log.Warnf("CreateCategoryHandler | Failed to create category: %v", err)
		switch {
		case errors.Is(err, common.ErrCategoryExists):
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Category with this name already exists",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to create category",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Category created successfully",
		Data:    category,
	})
}

// GetAllCategoriesHandler
// @Summary Get all categories
// @Description Get all categories
// @Tags Categories
// @Accept json
// @Produce json
// @Param limit query int false "Number of categories to return"
// @Param offset query int false "Offset for pagination"
// @Success 200 {object} common.SuccessResponse{data=[]dto.CategoryWithCountResponse} "Categories retrieved successfully"
// @Failure 400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve categories"
// @Router /categories [get]
func (h *CtgHandler) GetAllCategoriesHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(dto.ListCategoryRequest)
	if err := c.QueryParser(req); err != nil {
		h.log.Warnf("GetAllCategoriesHandler | Failed to parse query parameters: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
		})
	}

	categories, err := h.service.GetAllCategories(ctx, *req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrCategoryNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "No categories found",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to retrieve categories",
			})
		}
	}

	if len(categories) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
			Message: "No categories found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

// GetCategoryCountHandler
// @Summary Get total number of categories
// @Description Get total number of categories
// @Tags Categories
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse{data=int} "Category count retrieved successfully"
// @Failure 500 {object} common.ErrorResponse "Failed to retrieve category count"
// @Router /categories/count [get]
func (h *CtgHandler) GetCategoryCountHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	count, err := h.service.GetCategoryWithProductCount(ctx)
	if err != nil {
		h.log.Warnf("GetCategoryCountHandler | Failed to get category count: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to retrieve category count",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Category count retrieved successfully",
		Data:    count,
	})

}
