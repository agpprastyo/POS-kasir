package shift

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler interface {
	StartShiftHandler(c *fiber.Ctx) error
	EndShiftHandler(c *fiber.Ctx) error
	GetOpenShiftHandler(c *fiber.Ctx) error
	CreateCashTransactionHandler(c *fiber.Ctx) error
}

type handler struct {
	service   Service
	log       logger.ILogger
	validator validator.Validator
}

func NewHandler(service Service, log logger.ILogger, validator validator.Validator) Handler {
	return &handler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

// StartShiftHandler handles the request to start a new shift
// @Summary Start a new shift
// @Tags Shifts
// @Accept json
// @Produce json
// @Param request body dto.StartShiftRequest true "Start Shift Request"
// @Success 201 {object} common.SuccessResponse{data=dto.ShiftResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 409 {object} common.ErrorResponse "User already has an open shift"
// @Failure 500 {object} common.ErrorResponse
// @Router /shifts/start [post]
func (h *handler) StartShiftHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.StartShiftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	shift, err := h.service.StartShift(ctx, userID, req)
	if err != nil {
		h.log.Errorf("StartShiftHandler | Failed to start shift: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Shift started successfully",
		Data:    shift,
	})
}

// EndShiftHandler handles the request to end the current shift
// @Summary End current shift
// @Tags Shifts
// @Accept json
// @Produce json
// @Param request body dto.EndShiftRequest true "End Shift Request"
// @Success 200 {object} common.SuccessResponse{data=dto.ShiftResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse "No open shift found"
// @Failure 500 {object} common.ErrorResponse
// @Router /shifts/end [post]
func (h *handler) EndShiftHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.EndShiftRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	shift, err := h.service.EndShift(ctx, userID, req)
	if err != nil {
		h.log.Errorf("EndShiftHandler | Failed to end shift: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Shift ended successfully",
		Data:    shift,
	})
}

// GetOpenShiftHandler handles the request to get the current open shift
// @Summary Get current open shift
// @Tags Shifts
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse{data=dto.ShiftResponse}
// @Failure 404 {object} common.ErrorResponse "No open shift found"
// @Failure 500 {object} common.ErrorResponse
// @Router /shifts/current [get]
func (h *handler) GetOpenShiftHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	shift, err := h.service.GetOpenShift(ctx, userID)
	if err != nil {
		if err == common.ErrNotFound {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "No open shift found",
			})
		}
		h.log.Errorf("GetOpenShiftHandler | Failed to get open shift: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to get open shift",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Open shift retrieved successfully",
		Data:    shift,
	})
}

// CreateCashTransactionHandler handles the request to create a cash transaction
// @Summary Create a cash transaction (Drop/Expense/In)
// @Tags Shifts
// @Accept json
// @Produce json
// @Param request body dto.CashTransactionRequest true "Cash Transaction Request"
// @Success 201 {object} common.SuccessResponse{data=dto.CashTransactionResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse "No open shift found"
// @Failure 500 {object} common.ErrorResponse
// @Router /shifts/cash-transaction [post]
func (h *handler) CreateCashTransactionHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userID := c.Locals("user_id").(uuid.UUID)

	var req dto.CashTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	tx, err := h.service.CreateCashTransaction(ctx, userID, req)
	if err != nil {
		h.log.Errorf("CreateCashTransactionHandler | Failed to create cash transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Cash transaction created successfully",
		Data:    tx,
	})
}
