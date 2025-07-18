package orders

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IOrderHandler interface {
	CreateOrderHandler(c *fiber.Ctx) error
	GetOrderHandler(c *fiber.Ctx) error
	ProcessPaymentHandler(c *fiber.Ctx) error
	MidtransNotificationHandler(c *fiber.Ctx) error
	ListOrdersHandler(c *fiber.Ctx) error
	CancelOrderHandler(c *fiber.Ctx) error
	UpdateOrderItemsHandler(c *fiber.Ctx) error
	CompleteManualPaymentHandler(c *fiber.Ctx) error
	UpdateOperationalStatusHandler(c *fiber.Ctx) error
	ApplyPromotionHandler(c *fiber.Ctx) error
}

type OrderHandler struct {
	orderService IOrderService
	log          *logger.Logger
	validate     validator.Validator
}

func NewOrderHandler(orderService IOrderService, log *logger.Logger, validate validator.Validator) IOrderHandler {
	return &OrderHandler{
		orderService: orderService,
		log:          log,
		validate:     validate,
	}
}

func (h *OrderHandler) ApplyPromotionHandler(c *fiber.Ctx) error {
	// 1. Ambil ID pesanan dari parameter URL
	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	// 2. Parse request body
	var req dto.ApplyPromotionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	// 3. Validasi DTO
	if err := h.validate.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	// 4. Panggil service untuk menerapkan promosi
	orderResponse, err := h.orderService.ApplyPromotion(c.Context(), orderID, req)
	if err != nil {
		// Tangani error spesifik dari service, seperti promosi tidak valid atau pesanan tidak memenuhi syarat
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order or Promotion not found"})
		}
		if errors.Is(err, common.ErrPromotionNotApplicable) {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{Message: "Promotion cannot be applied", Error: err.Error()})
		}
		h.log.Error("Failed to apply promotion in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to apply promotion"})
	}

	// 5. Kirim respons sukses
	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Promotion applied successfully",
		Data:    orderResponse,
	})
}

func (h *OrderHandler) UpdateOperationalStatusHandler(c *fiber.Ctx) error {

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warn("Invalid order ID format for status update", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse update status request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Update status request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	orderResponse, err := h.orderService.UpdateOperationalStatus(c.Context(), orderID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		if errors.Is(err, common.ErrInvalidStatusTransition) {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{Message: "Invalid status transition", Error: err.Error()})
		}
		h.log.Error("Failed to update operational status in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update order status"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order status updated successfully",
		Data:    orderResponse,
	})
}
func (h *OrderHandler) CompleteManualPaymentHandler(c *fiber.Ctx) error {

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warn("Invalid order ID format for manual payment", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req dto.CompleteManualPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse complete manual payment request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Complete manual payment request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	orderResponse, err := h.orderService.CompleteManualPayment(c.Context(), orderID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		if errors.Is(err, common.ErrOrderNotModifiable) {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{Message: "Order cannot be processed", Error: "Order might have been paid or cancelled."})
		}
		h.log.Error("Failed to complete manual payment in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to complete payment"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Payment completed successfully",
		Data:    orderResponse,
	})
}

func (h *OrderHandler) UpdateOrderItemsHandler(c *fiber.Ctx) error {
	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req []dto.UpdateOrderItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body, expected an array of actions"})
	}

	if err := h.validate.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	updatedOrder, err := h.orderService.UpdateOrderItems(c.Context(), orderID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update order items"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order items updated successfully",
		Data:    updatedOrder,
	})
}

func (h *OrderHandler) CancelOrderHandler(c *fiber.Ctx) error {
	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warn("Invalid order ID format for cancellation", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req dto.CancelOrderRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse cancel order request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Cancel order request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	err = h.orderService.CancelOrder(c.Context(), orderID, req)
	if err != nil {

		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		if errors.Is(err, common.ErrOrderNotCancellable) {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{Message: "Order cannot be cancelled", Error: "Order might have been paid or already cancelled."})
		}
		h.log.Error("Failed to cancel order in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to cancel order"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order cancelled successfully",
	})
}

func (h *OrderHandler) ListOrdersHandler(c *fiber.Ctx) error {

	var req dto.ListOrdersRequest
	if err := c.QueryParser(&req); err != nil {
		h.log.Warn("Cannot parse list orders query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid query parameters"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("List orders request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	pagedResponse, err := h.orderService.ListOrders(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to list orders from service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve orders"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    pagedResponse,
	})
}

func (h *OrderHandler) CreateOrderHandler(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse create order request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Create order request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	orderResponse, err := h.orderService.CreateOrder(c.Context(), req)
	if err != nil {
		h.log.Error("Failed to create order in service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to create order"})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Order created successfully",
		Data:    orderResponse,
	})
}

func (h *OrderHandler) GetOrderHandler(c *fiber.Ctx) error {
	// 1. Ambil ID dari parameter URL
	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warn("Invalid order ID format", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	orderResponse, err := h.orderService.GetOrder(c.Context(), orderID)
	if err != nil {

		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		h.log.Error("Failed to get order from service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve order"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order retrieved successfully",
		Data:    orderResponse,
	})
}

func (h *OrderHandler) ProcessPaymentHandler(c *fiber.Ctx) error {

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warn("Invalid order ID format for payment", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	qrisResponse, err := h.orderService.ProcessPayment(c.Context(), orderID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		h.log.Error("Failed to process payment in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to process payment"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "QRIS payment initiated successfully",
		Data:    qrisResponse,
	})
}

func (h *OrderHandler) MidtransNotificationHandler(c *fiber.Ctx) error {

	var payload dto.MidtransNotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		h.log.Warn("Cannot parse Midtrans notification body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid notification format"})
	}

	err := h.orderService.HandleMidtransNotification(c.Context(), payload)
	if err != nil {
		h.log.Error("Error handling Midtrans notification", "error", err, "orderID", payload.OrderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to handle notification"})
	}

	h.log.Info("Successfully handled Midtrans notification", "orderID", payload.OrderID, "status", payload.TransactionStatus)
	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{Message: "Notification received successfully"})
}
