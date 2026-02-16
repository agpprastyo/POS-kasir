package shift

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type Handler interface {
	StartShiftHandler(c fiber.Ctx) error
	EndShiftHandler(c fiber.Ctx) error
	GetOpenShiftHandler(c fiber.Ctx) error
	CreateCashTransactionHandler(c fiber.Ctx) error
}

type handler struct {
	service Service
	log     logger.ILogger
}

func NewHandler(service Service, log logger.ILogger) Handler {
	return &handler{
		service: service,
		log:     log,
	}
}

// StartShiftHandler handles the request to start a new shift
// @Summary      Start a new shift
// @Description  Create a new shift session for the authenticated user (Roles: admin, manager, cashier)
// @Tags         Shifts
// @Accept       json
// @Produce      json
// @Param        request body StartShiftRequest true "Start Shift Request"
// @Success      201 {object} common.SuccessResponse{data=ShiftResponse} "Shift started successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body or validation failure"
// @Failure      409 {object} common.ErrorResponse "User already has an open shift"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /shifts/start [post]
func (h *handler) StartShiftHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Locals("user_id").(uuid.UUID)

	var req StartShiftRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Start shift validation failed", "error", err)
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
// @Summary      End current shift
// @Description  Close the active shift session for the authenticated user (Roles: admin, manager, cashier)
// @Tags         Shifts
// @Accept       json
// @Produce      json
// @Param        request body EndShiftRequest true "End Shift Request"
// @Success      200 {object} common.SuccessResponse{data=ShiftResponse} "Shift ended successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body or validation failure"
// @Failure      404 {object} common.ErrorResponse "No open shift found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /shifts/end [post]
func (h *handler) EndShiftHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Locals("user_id").(uuid.UUID)

	var req EndShiftRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("End shift validation failed", "error", err)
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
// @Summary      Get current open shift
// @Description  Check for and retrieve the details of an active shift session for the authenticated user (Roles: admin, manager, cashier)
// @Tags         Shifts
// @Accept       json
// @Produce      json
// @Success      200 {object} common.SuccessResponse{data=ShiftResponse} "Open shift retrieved successfully"
// @Failure      404 {object} common.ErrorResponse "No open shift found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /shifts/current [get]
func (h *handler) GetOpenShiftHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
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
// @Summary      Create a cash transaction (Drop/Expense/In)
// @Description  Record a manual cash entry or exit within the active shift (Roles: admin, manager, cashier)
// @Tags         Shifts
// @Accept       json
// @Produce      json
// @Param        request body CashTransactionRequest true "Cash Transaction Request"
// @Success      201 {object} common.SuccessResponse{data=CashTransactionResponse} "Cash transaction created successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body or validation failure"
// @Failure      404 {object} common.ErrorResponse "No open shift found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /shifts/cash-transaction [post]
func (h *handler) CreateCashTransactionHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Locals("user_id").(uuid.UUID)

	var req CashTransactionRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Create cash transaction validation failed", "error", err)
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
