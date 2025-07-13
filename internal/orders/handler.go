package orders

import (
	"POS-kasir/internal/common"
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

func (h *OrderHandler) CreateOrderHandler(c *fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warn("Cannot parse create order request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warn("Create order request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
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
	var payload MidtransNotificationPayload
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
