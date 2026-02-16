package orders

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/middleware"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/validator"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type IOrderHandler interface {
	CreateOrderHandler(c fiber.Ctx) error
	GetOrderHandler(c fiber.Ctx) error
	InitiateMidtransPaymentHandler(c fiber.Ctx) error
	MidtransNotificationHandler(c fiber.Ctx) error
	ListOrdersHandler(c fiber.Ctx) error
	CancelOrderHandler(c fiber.Ctx) error
	UpdateOrderItemsHandler(c fiber.Ctx) error
	ConfirmManualPaymentHandler(c fiber.Ctx) error
	UpdateOperationalStatusHandler(c fiber.Ctx) error
	ApplyPromotionHandler(c fiber.Ctx) error
}

type OrderHandler struct {
	orderService IOrderService
	log          logger.ILogger
}

func NewOrderHandler(orderService IOrderService, log logger.ILogger) IOrderHandler {
	return &OrderHandler{
		orderService: orderService,
		log:          log,
	}
}

// ApplyPromotionHandler applies a promotion to an order
// @Summary      Apply promotion to an order
// @Description  Apply a specific promotion to an existing order by its ID (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Param        request body ApplyPromotionRequest true "Promotion details"
// @Success      200 {object} common.SuccessResponse{data=OrderDetailResponse} "Promotion applied successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Order or Promotion not found"
// @Failure      500 {object} common.ErrorResponse "Failed to apply promotion"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/apply-promotion [post]
func (h *OrderHandler) ApplyPromotionHandler(c fiber.Ctx) error {

	orderID, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req ApplyPromotionRequest
	if err := c.Bind().Body(&req); err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	orderResponse, err := h.orderService.ApplyPromotion(c.RequestCtx(), orderID, req)
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

// UpdateOperationalStatusHandler updates the operational status of an order
// @Summary      Update order operational status
// @Description  Update the status of an existing order (e.g., to in_progress, served) (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Param        request body UpdateOrderStatusRequest true "Order status details"
// @Success      200 {object} common.SuccessResponse{data=OrderDetailResponse} "Order status updated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Order not found"
// @Failure      409 {object} common.ErrorResponse "Invalid status transition"
// @Failure      500 {object} common.ErrorResponse "Failed to update order status"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/update-status [post]
func (h *OrderHandler) UpdateOperationalStatusHandler(c fiber.Ctx) error {

	orderID, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid order ID format for status update", "error", err, "id", orderID)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req UpdateOrderStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Cannot parse update status request body", "error", err)
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	orderResponse, err := h.orderService.UpdateOperationalStatus(c.RequestCtx(), orderID, req)
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
// @Summary      Confirm manual payment for an order
// @Description  Process a manual payment (Cash) and finalize an order (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Param        request body ConfirmManualPaymentRequest true "Manual payment details"
// @Success      200 {object} common.SuccessResponse{data=OrderDetailResponse} "Payment completed successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Order not found"
// @Failure      409 {object} common.ErrorResponse "Order might have been paid or cancelled"
// @Failure      500 {object} common.ErrorResponse "Failed to complete payment"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/pay/manual [post]
func (h *OrderHandler) ConfirmManualPaymentHandler(c fiber.Ctx) error {

	orderID, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid order ID format for manual payment", "error", err, "id", orderID)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req ConfirmManualPaymentRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Cannot parse complete manual payment request body", "error", err)
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	orderResponse, err := h.orderService.ConfirmManualPayment(c.RequestCtx(), orderID, req)
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
// @Summary      Update items in an order
// @Description  Update, add, or remove items in an existing open order (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Param        request body []UpdateOrderItemRequest true "Update order items"
// @Success      200 {object} common.SuccessResponse{data=OrderDetailResponse} "Order items updated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Order not found"
// @Failure      500 {object} common.ErrorResponse "Failed to update order items"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/items [put]
func (h *OrderHandler) UpdateOrderItemsHandler(c fiber.Ctx) error {
	orderID, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req []UpdateOrderItemRequest
	if err := c.Bind().Body(&req); err != nil {
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body, expected an array of actions"})
	}

	updatedOrder, err := h.orderService.UpdateOrderItems(c.RequestCtx(), orderID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update order items"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Order items updated successfully",
		Data:    updatedOrder,
	})
}

// CancelOrderHandler cancels an order
// @Summary      Cancel an order
// @Description  Cancel an existing order with a reason (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Param        request body CancelOrderRequest true "Cancel order details"
// @Success      200 {object} common.SuccessResponse "Order cancelled successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID format or request body"
// @Failure      404 {object} common.ErrorResponse "Order not found"
// @Failure      409 {object} common.ErrorResponse "Order cannot be cancelled"
// @Failure      500 {object} common.ErrorResponse "Failed to cancel order"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/cancel [post]
func (h *OrderHandler) CancelOrderHandler(c fiber.Ctx) error {
	orderID, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("CancelOrderHandler | Invalid order ID format for cancellation", "error", err, "id", orderID)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	var req CancelOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("CancelOrderHandler | Cannot parse cancel order request body", "error", err)
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	err = h.orderService.CancelOrder(c.RequestCtx(), orderID, req)
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

// ListOrdersHandler lists all orders with filtering and pagination
// @Summary      List orders
// @Description  Get a list of orders with filtering by status and user (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Number of orders per page"
// @Param        statuses query []string false "Order statuses" collectionFormat(multi) Enums(open, in_progress, served, paid, cancelled)
// @Param        user_id query string false "Filter by User ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse{data=PagedOrderResponse} "Orders retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Failed to retrieve orders"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders [get]
func (h *OrderHandler) ListOrdersHandler(c fiber.Ctx) error {

	var req ListOrdersRequest
	if err := c.Bind().Query(&req); err != nil {
		h.log.Warnf("Cannot parse list orders query parameters", "error", err)
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid query parameters"})
	}

	h.log.Infof("List orders request", "request", req)

	currentRoleRaw := c.Locals("role")
	currentRole, ok := currentRoleRaw.(middleware.UserRole)
	if !ok {
		roleStr, okStr := currentRoleRaw.(string)
		if okStr {
			currentRole = middleware.UserRole(roleStr)
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

	if currentRole == middleware.UserRoleCashier {
		req.UserID = &currentUserID
	}

	pagedResponse, err := h.orderService.ListOrders(c.RequestCtx(), req)
	if err != nil {
		h.log.Errorf("Failed to list orders from service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve orders"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Orders retrieved successfully",
		Data:    pagedResponse,
	})
}

// CreateOrderHandler creates a new order
// @Summary      Create an order
// @Description  Create a new order with multiple items (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        request body CreateOrderRequest true "Create order details"
// @Success      201 {object} common.SuccessResponse{data=OrderDetailResponse} "Order created successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body"
// @Failure      500 {object} common.ErrorResponse "Failed to create order"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders [post]
func (h *OrderHandler) CreateOrderHandler(c fiber.Ctx) error {
	var req CreateOrderRequest
	if err := c.Bind().Body(&req); err != nil {
		h.log.Warnf("Cannot parse create order request body", "error", err)
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}
	h.log.Infof("CreateOrderHandler payload 1: %+v", req)

	orderResponse, err := h.orderService.CreateOrder(c.RequestCtx(), req)
	if err != nil {
		h.log.Errorf("Failed to create order in service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to create order"})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Order created successfully",
		Data:    orderResponse,
	})
}

// GetOrderHandler gets an order by ID
// @Summary      Get an order by ID
// @Description  Retrieve detailed information of a specific order by its ID (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse{data=OrderDetailResponse} "Order retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID format"
// @Failure      404 {object} common.ErrorResponse "Order not found"
// @Failure      500 {object} common.ErrorResponse "Failed to retrieve order"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetOrderHandler(c fiber.Ctx) error {

	orderID, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid order ID format", "error", err, "id", orderID)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	orderResponse, err := h.orderService.GetOrder(c.RequestCtx(), orderID)
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
// @Summary      Initiate Midtrans payment for an order
// @Description  Create a QRIS/Gopay payment session via Midtrans for an existing order (Roles: admin, manager, cashier)
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse{data=MidtransPaymentResponse} "QRIS payment initiated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID format"
// @Failure      404 {object} common.ErrorResponse "Order not found"
// @Failure      500 {object} common.ErrorResponse "Failed to process payment"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/pay/midtrans [post]
func (h *OrderHandler) InitiateMidtransPaymentHandler(c fiber.Ctx) error {

	orderID, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		h.log.Warnf("Invalid order ID format for payment", "error", err, "id", orderID)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid order ID format"})
	}

	qrisResponse, err := h.orderService.InitiateMidtransPayment(c.RequestCtx(), orderID)
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

// MidtransNotificationHandler handles midtrans payment notifications
// @Summary      Midtrans Payment Notification Callback
// @Description  Webhook for Midtrans to notify order payment status updates
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Param        payload body payment.MidtransNotificationPayload true "Midtrans Notification Payload"
// @Success      200 {object} common.SuccessResponse "Notification received successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid notification format"
// @Failure      500 {object} common.ErrorResponse "Failed to handle notification"
// @Router       /orders/webhook/midtrans [post]
func (h *OrderHandler) MidtransNotificationHandler(c fiber.Ctx) error {

	var payload payment.MidtransNotificationPayload
	if err := c.Bind().Body(&payload); err != nil {
		h.log.Warnf("Cannot parse Midtrans notification body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid notification format"})
	}

	err := h.orderService.HandleMidtransNotification(c.RequestCtx(), payload)
	if err != nil {
		h.log.Errorf("Error handling Midtrans notification", "error", err, "orderID", payload.OrderID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to handle notification"})
	}

	h.log.Infof("Successfully handled Midtrans notification", "orderID", payload.OrderID, "status", payload.TransactionStatus)
	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{Message: "Notification received successfully"})
}
