package promotions

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
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
	service  IPromotionService
	log      logger.ILogger
	validate validator.Validator
}

func NewPromotionHandler(service IPromotionService, log logger.ILogger, validate validator.Validator) IPromotionHandler {
	return &PromotionHandler{
		service:  service,
		log:      log,
		validate: validate,
	}
}

// CreatePromotionHandler
// @Summary Create a new promotion
// @Tags Promotions
// @Accept json
// @Produce json
// @Param request body dto.CreatePromotionRequest true "Promotion details"
// @Success 201 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @x-roles ["admin", "manager"]
// @Router /promotions [post]
func (h *PromotionHandler) CreatePromotionHandler(c fiber.Ctx) error {
	var req dto.CreatePromotionRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Cannot parse create promotion request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("Create promotion request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
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

// UpdatePromotionHandler
// @Summary Update a promotion
// @Tags Promotions
// @Accept json
// @Produce json
// @Param id path string true "Promotion ID"
// @Param request body dto.UpdatePromotionRequest true "Promotion details"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @x-roles ["admin", "manager"]
// @Router /promotions/{id} [put]
func (h *PromotionHandler) UpdatePromotionHandler(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid ID format"})
	}

	var req dto.UpdatePromotionRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Cannot parse update promotion request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("Update promotion request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
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

// GetPromotionHandler
// @Summary Get a promotion by ID
// @Tags Promotions
// @Accept json
// @Produce json
// @Param id path string true "Promotion ID"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @x-roles ["cashier", "admin", "manager"]
// @Router /promotions/{id} [get]
func (h *PromotionHandler) GetPromotionHandler(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid ID format"})
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

// DeletePromotionHandler
// @Summary Delete (deactivate) a promotion
// @Tags Promotions
// @Accept json
// @Produce json
// @Param id path string true "Promotion ID"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @x-roles ["admin", "manager"]
// @Router /promotions/{id} [delete]
func (h *PromotionHandler) DeletePromotionHandler(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid ID format"})
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

// ListPromotionsHandler
// @Summary List all promotions
// @Tags Promotions
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Param trash query boolean false "Show trash items"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @x-roles ["cashier", "admin", "manager"]
// @Router /promotions [get]
func (h *PromotionHandler) ListPromotionsHandler(c fiber.Ctx) error {
	var req dto.ListPromotionsRequest
	if err := c.Bind().Query(&req); err != nil {
		h.log.Warnf("Cannot parse list promotions query", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid query parameters"})
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

// RestorePromotionHandler
// @Summary Restore a deleted promotion
// @Tags Promotions
// @Accept json
// @Produce json
// @Param id path string true "Promotion ID"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @x-roles ["admin", "manager"]
// @Router /promotions/{id}/restore [post]
func (h *PromotionHandler) RestorePromotionHandler(c fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid ID format"})
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

// fiber:context-methods migrated
