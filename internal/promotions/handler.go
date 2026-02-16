package promotions

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type IPromotionHandler interface {
	CreatePromotionHandler(c fiber.Ctx) error
	GetPromotionHandler(c fiber.Ctx) error
	UpdatePromotionHandler(c fiber.Ctx) error
	DeletePromotionHandler(c fiber.Ctx) error
	ListPromotionsHandler(c fiber.Ctx) error
	RestorePromotionHandler(c fiber.Ctx) error
}

type PromotionHandler struct {
	service IPromotionService
	log     logger.ILogger
}

func NewPromotionHandler(service IPromotionService, log logger.ILogger) IPromotionHandler {
	return &PromotionHandler{
		service: service,
		log:     log,
	}
}

// CreatePromotionHandler creates a new promotion
// @Summary      Create a new promotion
// @Description  Create a new promotion with rules and targets (Roles: admin, manager)
// @Tags         Promotions
// @Accept       json
// @Produce      json
// @Param        request body CreatePromotionRequest true "Promotion details"
// @Success      201 {object} common.SuccessResponse{data=PromotionResponse} "Promotion created successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body or validation failed"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager"]
// @Router       /promotions [post]
func (h *PromotionHandler) CreatePromotionHandler(c fiber.Ctx) error {
	var req CreatePromotionRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Create promotion validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	promo, err := h.service.CreatePromotion(c.RequestCtx(), req)
	if err != nil {
		h.log.Errorf("Failed to create promotion", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to create promotion"})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Promotion created successfully",
		Data:    promo,
	})
}

// UpdatePromotionHandler updates a promotion
// @Summary      Update a promotion
// @Description  Update details of an existing promotion by its ID (Roles: admin, manager)
// @Tags         Promotions
// @Accept       json
// @Produce      json
// @Param        id      path string true "Promotion ID" Format(uuid)
// @Param        request body UpdatePromotionRequest true "Promotion details"
// @Success      200 {object} common.SuccessResponse{data=PromotionResponse} "Promotion updated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid project ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Promotion not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager"]
// @Router       /promotions/{id} [put]
func (h *PromotionHandler) UpdatePromotionHandler(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid promotion ID format", "error", err, "id", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid promotion ID format"})
	}

	var req UpdatePromotionRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Update promotion validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	promo, err := h.service.UpdatePromotion(c.RequestCtx(), id, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Promotion not found"})
		}
		h.log.Errorf("Failed to update promotion", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update promotion"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Promotion updated successfully",
		Data:    promo,
	})
}

// GetPromotionHandler gets a promotion by ID
// @Summary      Get a promotion by ID
// @Description  Retrieve details of a specific promotion by its ID (Roles: admin, manager, cashier)
// @Tags         Promotions
// @Accept       json
// @Produce      json
// @Param        id path string true "Promotion ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse{data=PromotionResponse} "Promotion retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid promotion ID format"
// @Failure      404 {object} common.ErrorResponse "Promotion not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /promotions/{id} [get]
func (h *PromotionHandler) GetPromotionHandler(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid promotion ID format", "error", err, "id", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid promotion ID format"})
	}

	promo, err := h.service.GetPromotion(c.RequestCtx(), id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Promotion not found"})
		}
		h.log.Errorf("Failed to get promotion", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to get promotion"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Promotion retrieved successfully",
		Data:    promo,
	})
}

// DeletePromotionHandler deletes a promotion
// @Summary      Delete (deactivate) a promotion
// @Description  Soft delete a promotion by its ID (Roles: admin, manager)
// @Tags         Promotions
// @Accept       json
// @Produce      json
// @Param        id path string true "Promotion ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "Promotion deleted successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid promotion ID format"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager"]
// @Router       /promotions/{id} [delete]
func (h *PromotionHandler) DeletePromotionHandler(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid promotion ID format", "error", err, "id", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid promotion ID format"})
	}

	err = h.service.DeletePromotion(c.RequestCtx(), id)
	if err != nil {
		h.log.Errorf("Failed to delete promotion", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to delete promotion"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Promotion deleted successfully",
	})
}

// ListPromotionsHandler lists all promotions
// @Summary      List all promotions
// @Description  Get a list of promotions with pagination and optional trash filter (Roles: admin, manager, cashier)
// @Tags         Promotions
// @Accept       json
// @Produce      json
// @Param        page  query int     false "Page number"
// @Param        limit query int     false "Items per page"
// @Param        trash query boolean false "Show trash items"
// @Success      200 {object} common.SuccessResponse{data=PagedPromotionResponse} "Promotions retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /promotions [get]
func (h *PromotionHandler) ListPromotionsHandler(c fiber.Ctx) error {
	var req ListPromotionsRequest
	if err := c.Bind().Query(&req); err != nil {
		h.log.Warnf("List promotions validation failed", "error", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	resp, err := h.service.ListPromotions(c.RequestCtx(), req)
	if err != nil {
		h.log.Errorf("Failed to list promotions", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to list promotions"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Promotions retrieved successfully",
		Data:    resp,
	})
}

// RestorePromotionHandler restores a deleted promotion
// @Summary      Restore a deleted promotion
// @Description  Restore a soft-deleted promotion by its ID (Roles: admin, manager)
// @Tags         Promotions
// @Accept       json
// @Produce      json
// @Param        id path string true "Promotion ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "Promotion restored successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid promotion ID format"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager"]
// @Router       /promotions/{id}/restore [post]
func (h *PromotionHandler) RestorePromotionHandler(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid promotion ID format", "error", err, "id", c.Params("id"))
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid promotion ID format"})
	}

	err = h.service.RestorePromotion(c.RequestCtx(), id)
	if err != nil {
		h.log.Errorf("Failed to restore promotion", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to restore promotion"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Promotion restored successfully",
	})
}
