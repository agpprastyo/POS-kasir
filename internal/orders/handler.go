package orders

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IOrderHandler interface {
	CreateOrderHandler(c *fiber.Ctx) error
	GetOrderHandler(c *fiber.Ctx) error
	InitiateMidtransPaymentHandler(c *fiber.Ctx) error
	MidtransNotificationHandler(c *fiber.Ctx) error
	ListOrdersHandler(c *fiber.Ctx) error
	CancelOrderHandler(c *fiber.Ctx) error
	UpdateOrderItemsHandler(c *fiber.Ctx) error
	ConfirmManualPaymentHandler(c *fiber.Ctx) error
	UpdateOperationalStatusHandler(c *fiber.Ctx) error
	ApplyPromotionHandler(c *fiber.Ctx) error
}

type OrderHandler struct {
	orderService IOrderService
	log          logger.ILogger
	validate     validator.Validator
}

func NewOrderHandler(orderService IOrderService, log logger.ILogger, validate validator.Validator) IOrderHandler {
	return &OrderHandler{
		orderService: orderService,
		log:          log,
		validate:     validate,
	}
}

// ListProductOptionsHandler is a placeholder for listing product options
// @Summary Apply promotion to an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.ApplyPromotionRequest true "Promotion details"
// @Success 200 {object} common.SuccessResponse{data=dto.OrderDetailResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/{id}/apply-promotion [post]
func (h *OrderHandler) ApplyPromotionHandler(c *fiber.Ctx) error {

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req dto.ApplyPromotionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	orderResponse, err := h.orderService.ApplyPromotion(c.Context(), orderID, req)
	if err != nil {

		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order or Promotion not found"})
		}
		if errors.Is(err, common.ErrPromotionNotApplicable) {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{Message: "Promotion cannot be applied", Error: err.Error()})
		}
		h.log.Errorf("Failed to apply promotion in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to apply promotion"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Promotion applied successfully",
		Data:    orderResponse,
	})
}

// UpdateOperationalStatusHandler is a placeholder for updating order operational status
// @Summary Update order operational status
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.UpdateOrderStatusRequest true "Order status details"
// @Success 200 {object} common.SuccessResponse{data=dto.OrderDetailResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/{id}/update-status [post]
func (h *OrderHandler) UpdateOperationalStatusHandler(c *fiber.Ctx) error {

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warnf("Invalid order ID format for status update", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req dto.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warnf("Cannot parse update status request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("Update status request validation failed", "error", err)
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
		h.log.Errorf("Failed to update operational status in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update order status"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order status updated successfully",
		Data:    orderResponse,
	})
}

// ConfirmManualPaymentHandler confirms manual payment for an order
// @Summary Confirm manual payment for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.ConfirmManualPaymentRequest true "Manual payment details"
// @Success 200 {object} common.SuccessResponse{data=dto.OrderDetailResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/{id}/pay/manual [post]
func (h *OrderHandler) ConfirmManualPaymentHandler(c *fiber.Ctx) error {

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warnf("Invalid order ID format for manual payment", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req dto.ConfirmManualPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warnf("Cannot parse complete manual payment request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("Complete manual payment request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	orderResponse, err := h.orderService.ConfirmManualPayment(c.Context(), orderID, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		if errors.Is(err, common.ErrOrderNotModifiable) {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{Message: "Order cannot be processed", Error: "Order might have been paid or cancelled."})
		}
		h.log.Errorf("Failed to complete manual payment in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to complete payment"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Payment completed successfully",
		Data:    orderResponse,
	})
}

// UpdateOrderItemsHandler updates items in an order
// @Summary Update items in an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body []dto.UpdateOrderItemRequest true "Update order items"
// @Success 200 {object} common.SuccessResponse{data=dto.OrderDetailResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/{id}/items [put]
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

// CancelOrderHandler is a placeholder for cancelling an order
// @Summary Cancel an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Param request body dto.CancelOrderRequest true "Cancel order details"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrderHandler(c *fiber.Ctx) error {
	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warnf("CancelOrderHandler | Invalid order ID format for cancellation", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req dto.CancelOrderRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warnf("CancelOrderHandler | Cannot parse cancel order request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("Cancel order request validation failed", "error", err)
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
		h.log.Errorf("Failed to cancel order in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to cancel order"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order cancelled successfully",
	})
}

// ListOrdersHandler is a placeholder for listing orders
// @Summary List orders
// @Tags Orders
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of orders per page"
// @Param statuses query []string false "Order statuses" collectionFormat(multi) Enums(open, in_progress, served, paid, cancelled)
// @Param user_id query string false "Filter by User ID"
// @Success 200 {object} common.SuccessResponse{data=dto.PagedOrderResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders [get]
func (h *OrderHandler) ListOrdersHandler(c *fiber.Ctx) error {

	var req dto.ListOrdersRequest
	if err := c.QueryParser(&req); err != nil {
		h.log.Warnf("Cannot parse list orders query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid query parameters"})
	}

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("List orders request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error()})
	}

	h.log.Infof("List orders request", "request", req)

	currentRoleRaw := c.Locals("role")
	currentRole, ok := currentRoleRaw.(repository.UserRole)
	if !ok {
		roleStr, okStr := currentRoleRaw.(string)
		if okStr {
			currentRole = repository.UserRole(roleStr)
		} else {
			h.log.Errorf("Failed to retrieve role from context. Type: %T", currentRoleRaw)
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{Message: "Unauthorized"})
		}
	}

	currentUserIDRaw := c.Locals("user_id")
	currentUserID, ok := currentUserIDRaw.(uuid.UUID)
	if !ok {
		h.log.Errorf("Failed to retrieve user_id from context. Type: %T", currentUserIDRaw)
		return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{Message: "Unauthorized"})
	}

	if currentRole == repository.UserRoleCashier {
		req.UserID = &currentUserID
	}

	pagedResponse, err := h.orderService.ListOrders(c.Context(), req)
	if err != nil {
		h.log.Errorf("Failed to list orders from service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve orders"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    pagedResponse,
	})
}

// CreateOrderHandler is a placeholder for creating an order
// @Summary Create an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body dto.CreateOrderRequest true "Create order details"
// @Success 201 {object} common.SuccessResponse{data=dto.OrderDetailResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders [post]
func (h *OrderHandler) CreateOrderHandler(c *fiber.Ctx) error {
	var req dto.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		h.log.Warnf("Cannot parse create order request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}
	h.log.Infof("CreateOrderHandler payload 1: %+v", req)

	if err := h.validate.Validate(req); err != nil {
		h.log.Warnf("Create order request validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Validation failed", Error: err.Error(), Data: req})
	}

	orderResponse, err := h.orderService.CreateOrder(c.Context(), req)
	if err != nil {
		h.log.Errorf("Failed to create order in service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to create order"})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Order created successfully",
		Data:    orderResponse,
	})
}

// GetOrderHandler is a placeholder for getting an order by ID
// @Summary Get an order by ID
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} common.SuccessResponse{data=dto.OrderDetailResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/{id} [get]
func (h *OrderHandler) GetOrderHandler(c *fiber.Ctx) error {
	// 1. Ambil ID dari parameter URL
	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warnf("Invalid order ID format", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	orderResponse, err := h.orderService.GetOrder(c.Context(), orderID)
	if err != nil {

		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		h.log.Errorf("Failed to get order from service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve order"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order retrieved successfully",
		Data:    orderResponse,
	})
}

// InitiateMidtransPaymentHandler initiates midtrans payment for an order
// @Summary Initiate Midtrans payment for an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} common.SuccessResponse{data=dto.MidtransPaymentResponse}
// @Failure 400 {object} common.ErrorResponse
// @Failure 404 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /orders/{id}/pay/midtrans [post]
func (h *OrderHandler) InitiateMidtransPaymentHandler(c *fiber.Ctx) error {

	orderIDStr := c.Params("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		h.log.Warnf("Invalid order ID format for payment", "error", err, "id", orderIDStr)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	qrisResponse, err := h.orderService.InitiateMidtransPayment(c.Context(), orderID)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Order not found"})
		}
		h.log.Errorf("Failed to process payment in service", "error", err, "orderID", orderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to process payment: " + err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "QRIS payment initiated successfully",
		Data:    qrisResponse,
	})
}

func (h *OrderHandler) MidtransNotificationHandler(c *fiber.Ctx) error {

	var payload dto.MidtransNotificationPayload
	if err := c.BodyParser(&payload); err != nil {
		h.log.Warnf("Cannot parse Midtrans notification body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid notification format"})
	}

	err := h.orderService.HandleMidtransNotification(c.Context(), payload)
	if err != nil {
		h.log.Errorf("Error handling Midtrans notification", "error", err, "orderID", payload.OrderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to handle notification"})
	}

	h.log.Infof("Successfully handled Midtrans notification", "orderID", payload.OrderID, "status", payload.TransactionStatus)
	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{Message: "Notification received successfully"})
}
