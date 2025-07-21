package categories

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"errors"
	"github.com/gofiber/fiber/v2"
)

type CtgHandler struct {
	service ICtgService
	log     logger.ILogger
}

func (h *CtgHandler) DeleteCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category ID is required",
		})
	}

	err := h.service.DeleteCategory(ctx, id)
	if err != nil {
		h.log.Errorf("Failed to delete category", "error", err, "categoryID", id)
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

func (h *CtgHandler) UpdateCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category ID is required",
		})
	}

	req := new(CreateCategoryRequest)
	if err := c.BodyParser(req); err != nil {
		h.log.Errorf("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category name is required",
		})
	}

	category, err := h.service.UpdateCategory(ctx, id, *req)
	if err != nil {
		h.log.Errorf("Failed to update category", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to update category",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Category updated successfully",
		Data:    category,
	})
}

func (h *CtgHandler) GetCategoryByIDHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Category ID is required",
		})
	}

	category, err := h.service.GetCategoryByID(ctx, id)
	if err != nil {
		h.log.Errorf("Failed to get category by ID", "error", err)
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

func (h *CtgHandler) CreateCategoryHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(CreateCategoryRequest)
	if err := c.BodyParser(req); err != nil {
		h.log.Errorf("Failed to parse request body", "error", err)
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
		h.log.Errorf("Failed to create category", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to create category",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Category created successfully",
		Data:    category,
	})
}

func (h *CtgHandler) GetAllCategoriesHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(ListCategoryRequest)
	if err := c.QueryParser(req); err != nil {
		h.log.Errorf("Failed to parse query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
		})
	}

	categories, err := h.service.GetAllCategories(ctx, *req)
	if err != nil {
		h.log.Errorf("Failed to get all categories", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to retrieve categories",
		})
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

type ICtgHandler interface {
	GetAllCategoriesHandler(c *fiber.Ctx) error
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
