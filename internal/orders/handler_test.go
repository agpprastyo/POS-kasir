package orders_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/middleware"
	"POS-kasir/internal/orders"
	orders_repo "POS-kasir/internal/orders/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/validator"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(t *testing.T) (*mocks.MockIOrderService, *mocks.MockFieldLogger, *orders.OrderHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIOrderService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	handler := orders.NewOrderHandler(mockService, mockLogger).(*orders.OrderHandler)
	app := fiber.New(fiber.Config{
		StructValidator: validator.NewValidator(),
	})
	return mockService, mockLogger, handler, app
}

func allowAllHandlerLoggerCalls(mockLogger *mocks.MockFieldLogger) {
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

// ====================== CreateOrderHandler ======================

func TestOrderHandler_CreateOrderHandler(t *testing.T) {
	productID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders", handler.CreateOrderHandler)

		reqBody := orders.CreateOrderRequest{
			Type: orders_repo.OrderTypeDineIn,
			Items: []orders.CreateOrderItemRequest{
				{ProductID: productID, Quantity: 2},
			},
		}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(&orders.OrderDetailResponse{
			ID:     uuid.New(),
			Type:   orders_repo.OrderTypeDineIn,
			Status: orders_repo.OrderStatusOpen,
		}, nil)

		req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders", handler.CreateOrderHandler)

		req := httptest.NewRequest("POST", "/orders", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationError", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders", handler.CreateOrderHandler)

		// Missing required fields
		reqBody := orders.CreateOrderRequest{}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders", handler.CreateOrderHandler)

		reqBody := orders.CreateOrderRequest{
			Type: orders_repo.OrderTypeDineIn,
			Items: []orders.CreateOrderItemRequest{
				{ProductID: productID, Quantity: 2},
			},
		}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		req := httptest.NewRequest("POST", "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== GetOrderHandler ======================

func TestOrderHandler_GetOrderHandler(t *testing.T) {
	orderID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders/:id", handler.GetOrderHandler)

		mockService.EXPECT().GetOrder(gomock.Any(), orderID).Return(&orders.OrderDetailResponse{
			ID:     orderID,
			Status: orders_repo.OrderStatusOpen,
		}, nil)

		req := httptest.NewRequest("GET", "/orders/"+orderID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders/:id", handler.GetOrderHandler)

		req := httptest.NewRequest("GET", "/orders/not-a-uuid", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders/:id", handler.GetOrderHandler)

		mockService.EXPECT().GetOrder(gomock.Any(), orderID).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("GET", "/orders/"+orderID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders/:id", handler.GetOrderHandler)

		mockService.EXPECT().GetOrder(gomock.Any(), orderID).Return(nil, errors.New("db error"))

		req := httptest.NewRequest("GET", "/orders/"+orderID.String(), nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== InitiateMidtransPaymentHandler ======================

func TestOrderHandler_InitiateMidtransPaymentHandler(t *testing.T) {
	orderID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/midtrans", handler.InitiateMidtransPaymentHandler)

		mockService.EXPECT().InitiateMidtransPayment(gomock.Any(), orderID).Return(&orders.MidtransPaymentResponse{
			OrderID:       orderID.String(),
			TransactionID: "txn-123",
		}, nil)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/midtrans", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/midtrans", handler.InitiateMidtransPaymentHandler)

		req := httptest.NewRequest("POST", "/orders/bad-id/pay/midtrans", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/midtrans", handler.InitiateMidtransPaymentHandler)

		mockService.EXPECT().InitiateMidtransPayment(gomock.Any(), orderID).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/midtrans", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/midtrans", handler.InitiateMidtransPaymentHandler)

		mockService.EXPECT().InitiateMidtransPayment(gomock.Any(), orderID).Return(nil, errors.New("midtrans error"))

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/midtrans", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== MidtransNotificationHandler ======================

func TestOrderHandler_MidtransNotificationHandler(t *testing.T) {

	t.Run("Success", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/webhook/midtrans", handler.MidtransNotificationHandler)

		payload := payment.MidtransNotificationPayload{
			OrderID:           uuid.New().String(),
			TransactionID:     "txn-456",
			TransactionStatus: "settlement",
			StatusCode:        "200",
			SignatureKey:      "sig",
		}
		body, _ := json.Marshal(payload)

		mockService.EXPECT().HandleMidtransNotification(gomock.Any(), gomock.Any()).Return(nil)

		req := httptest.NewRequest("POST", "/webhook/midtrans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/webhook/midtrans", handler.MidtransNotificationHandler)

		req := httptest.NewRequest("POST", "/webhook/midtrans", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/webhook/midtrans", handler.MidtransNotificationHandler)

		payload := payment.MidtransNotificationPayload{
			OrderID:           uuid.New().String(),
			TransactionStatus: "settlement",
		}
		body, _ := json.Marshal(payload)

		mockService.EXPECT().HandleMidtransNotification(gomock.Any(), gomock.Any()).Return(errors.New("notification error"))

		req := httptest.NewRequest("POST", "/webhook/midtrans", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== ListOrdersHandler ======================

func TestOrderHandler_ListOrdersHandler(t *testing.T) {

	t.Run("SuccessAdmin", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders", func(c fiber.Ctx) error {
			c.Locals("role", middleware.UserRoleAdmin)
			c.Locals("user_id", uuid.New())
			return c.Next()
		}, handler.ListOrdersHandler)

		mockService.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return(&orders.PagedOrderResponse{
			Orders: []orders.OrderListResponse{},
		}, nil)

		req := httptest.NewRequest("GET", "/orders", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("SuccessCashier", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		cashierID := uuid.New()
		app.Get("/orders", func(c fiber.Ctx) error {
			c.Locals("role", middleware.UserRoleCashier)
			c.Locals("user_id", cashierID)
			return c.Next()
		}, handler.ListOrdersHandler)

		// For cashier, the user_id should be set to their own
		mockService.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return(&orders.PagedOrderResponse{
			Orders: []orders.OrderListResponse{},
		}, nil)

		req := httptest.NewRequest("GET", "/orders", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("RoleAsString", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders", func(c fiber.Ctx) error {
			c.Locals("role", "admin") // string, not UserRole type
			c.Locals("user_id", uuid.New())
			return c.Next()
		}, handler.ListOrdersHandler)

		mockService.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return(&orders.PagedOrderResponse{
			Orders: []orders.OrderListResponse{},
		}, nil)

		req := httptest.NewRequest("GET", "/orders", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NoRole", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders", func(c fiber.Ctx) error {
			c.Locals("role", 123) // invalid type
			c.Locals("user_id", uuid.New())
			return c.Next()
		}, handler.ListOrdersHandler)

		req := httptest.NewRequest("GET", "/orders", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("NoUserID", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders", func(c fiber.Ctx) error {
			c.Locals("role", middleware.UserRoleAdmin)
			c.Locals("user_id", "not-a-uuid") // invalid type
			return c.Next()
		}, handler.ListOrdersHandler)

		req := httptest.NewRequest("GET", "/orders", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Get("/orders", func(c fiber.Ctx) error {
			c.Locals("role", middleware.UserRoleAdmin)
			c.Locals("user_id", uuid.New())
			return c.Next()
		}, handler.ListOrdersHandler)

		mockService.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))

		req := httptest.NewRequest("GET", "/orders", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== CancelOrderHandler ======================

func TestOrderHandler_CancelOrderHandler(t *testing.T) {
	orderID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/cancel", handler.CancelOrderHandler)

		reqBody := orders.CancelOrderRequest{
			CancellationReasonID: 1,
			CancellationNotes:    "Customer changed mind",
		}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CancelOrder(gomock.Any(), orderID, gomock.Any()).Return(nil)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/cancel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/cancel", handler.CancelOrderHandler)

		req := httptest.NewRequest("POST", "/orders/bad-id/cancel", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/cancel", handler.CancelOrderHandler)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/cancel", bytes.NewReader([]byte("bad")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationError", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/cancel", handler.CancelOrderHandler)

		// Missing required cancellation_reason_id
		reqBody := orders.CancelOrderRequest{}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/cancel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/cancel", handler.CancelOrderHandler)

		reqBody := orders.CancelOrderRequest{CancellationReasonID: 1}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CancelOrder(gomock.Any(), orderID, gomock.Any()).Return(common.ErrNotFound)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/cancel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("NotCancellable", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/cancel", handler.CancelOrderHandler)

		reqBody := orders.CancelOrderRequest{CancellationReasonID: 1}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CancelOrder(gomock.Any(), orderID, gomock.Any()).Return(common.ErrOrderNotCancellable)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/cancel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/cancel", handler.CancelOrderHandler)

		reqBody := orders.CancelOrderRequest{CancellationReasonID: 1}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().CancelOrder(gomock.Any(), orderID, gomock.Any()).Return(errors.New("db error"))

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/cancel", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== UpdateOrderItemsHandler ======================

func TestOrderHandler_UpdateOrderItemsHandler(t *testing.T) {
	orderID := uuid.New()
	productID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockIOrderService(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		handler := orders.NewOrderHandler(mockService, mockLogger).(*orders.OrderHandler)
		// Use Fiber app WITHOUT struct validator since Validate.Struct can't handle slices
		app := fiber.New()
		app.Put("/orders/:id/items", handler.UpdateOrderItemsHandler)

		jsonBody := fmt.Sprintf(`[{"product_id":"%s","quantity":3}]`, productID.String())

		mockService.EXPECT().UpdateOrderItems(gomock.Any(), orderID, gomock.Any()).Return(&orders.OrderDetailResponse{
			ID:     orderID,
			Status: orders_repo.OrderStatusOpen,
		}, nil)

		req := httptest.NewRequest("PUT", "/orders/"+orderID.String()+"/items", bytes.NewReader([]byte(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, _, handler, app := setupHandlerTest(t)
		app.Put("/orders/:id/items", handler.UpdateOrderItemsHandler)

		req := httptest.NewRequest("PUT", "/orders/bad-id/items", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		_, _, handler, app := setupHandlerTest(t)
		app.Put("/orders/:id/items", handler.UpdateOrderItemsHandler)

		req := httptest.NewRequest("PUT", "/orders/"+orderID.String()+"/items", bytes.NewReader([]byte("bad json")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockIOrderService(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		handler := orders.NewOrderHandler(mockService, mockLogger).(*orders.OrderHandler)
		app := fiber.New()
		app.Put("/orders/:id/items", handler.UpdateOrderItemsHandler)

		jsonBody := fmt.Sprintf(`[{"product_id":"%s","quantity":3}]`, productID.String())

		mockService.EXPECT().UpdateOrderItems(gomock.Any(), orderID, gomock.Any()).Return(nil, errors.New("error"))

		req := httptest.NewRequest("PUT", "/orders/"+orderID.String()+"/items", bytes.NewReader([]byte(jsonBody)))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== ConfirmManualPaymentHandler ======================

func TestOrderHandler_ConfirmManualPaymentHandler(t *testing.T) {
	orderID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/manual", handler.ConfirmManualPaymentHandler)

		reqBody := orders.ConfirmManualPaymentRequest{
			PaymentMethodID: 1,
			CashReceived:    50000,
		}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ConfirmManualPayment(gomock.Any(), orderID, gomock.Any()).Return(&orders.OrderDetailResponse{
			ID:     orderID,
			Status: orders_repo.OrderStatusPaid,
		}, nil)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/manual", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/manual", handler.ConfirmManualPaymentHandler)

		req := httptest.NewRequest("POST", "/orders/bad-id/pay/manual", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/manual", handler.ConfirmManualPaymentHandler)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/manual", bytes.NewReader([]byte("bad")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationError", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/manual", handler.ConfirmManualPaymentHandler)

		// Missing required payment_method_id
		reqBody := orders.ConfirmManualPaymentRequest{}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/manual", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/manual", handler.ConfirmManualPaymentHandler)

		reqBody := orders.ConfirmManualPaymentRequest{PaymentMethodID: 1, CashReceived: 50000}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ConfirmManualPayment(gomock.Any(), orderID, gomock.Any()).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/manual", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("OrderNotModifiable", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/manual", handler.ConfirmManualPaymentHandler)

		reqBody := orders.ConfirmManualPaymentRequest{PaymentMethodID: 1, CashReceived: 50000}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ConfirmManualPayment(gomock.Any(), orderID, gomock.Any()).Return(nil, common.ErrOrderNotModifiable)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/manual", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/pay/manual", handler.ConfirmManualPaymentHandler)

		reqBody := orders.ConfirmManualPaymentRequest{PaymentMethodID: 1, CashReceived: 50000}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ConfirmManualPayment(gomock.Any(), orderID, gomock.Any()).Return(nil, errors.New("db error"))

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/pay/manual", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== UpdateOperationalStatusHandler ======================

func TestOrderHandler_UpdateOperationalStatusHandler(t *testing.T) {
	orderID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/update-status", handler.UpdateOperationalStatusHandler)

		reqBody := orders.UpdateOrderStatusRequest{
			Status: orders_repo.OrderStatusInProgress,
		}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateOperationalStatus(gomock.Any(), orderID, gomock.Any()).Return(&orders.OrderDetailResponse{
			ID:     orderID,
			Status: orders_repo.OrderStatusInProgress,
		}, nil)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/update-status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/update-status", handler.UpdateOperationalStatusHandler)

		req := httptest.NewRequest("POST", "/orders/bad-id/update-status", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/update-status", handler.UpdateOperationalStatusHandler)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/update-status", bytes.NewReader([]byte("bad")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationError", func(t *testing.T) {
		_, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/update-status", handler.UpdateOperationalStatusHandler)

		// Invalid status value
		body, _ := json.Marshal(map[string]string{"status": "invalid_status"})

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/update-status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/update-status", handler.UpdateOperationalStatusHandler)

		reqBody := orders.UpdateOrderStatusRequest{Status: orders_repo.OrderStatusInProgress}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateOperationalStatus(gomock.Any(), orderID, gomock.Any()).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/update-status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InvalidTransition", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/update-status", handler.UpdateOperationalStatusHandler)

		reqBody := orders.UpdateOrderStatusRequest{Status: orders_repo.OrderStatusInProgress}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateOperationalStatus(gomock.Any(), orderID, gomock.Any()).Return(nil,
			fmt.Errorf("%w: cancelled to in_progress", common.ErrInvalidStatusTransition),
		)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/update-status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/update-status", handler.UpdateOperationalStatusHandler)

		reqBody := orders.UpdateOrderStatusRequest{Status: orders_repo.OrderStatusInProgress}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().UpdateOperationalStatus(gomock.Any(), orderID, gomock.Any()).Return(nil, errors.New("db error"))

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/update-status", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== ApplyPromotionHandler ======================

func TestOrderHandler_ApplyPromotionHandler(t *testing.T) {
	orderID := uuid.New()
	promoID := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockService, _, handler, app := setupHandlerTest(t)
		app.Post("/orders/:id/apply-promotion", handler.ApplyPromotionHandler)

		reqBody := orders.ApplyPromotionRequest{
			PromotionID: promoID,
		}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ApplyPromotion(gomock.Any(), orderID, gomock.Any()).Return(&orders.OrderDetailResponse{
			ID:     orderID,
			Status: orders_repo.OrderStatusOpen,
		}, nil)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/apply-promotion", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, _, handler, app := setupHandlerTest(t)
		app.Post("/orders/:id/apply-promotion", handler.ApplyPromotionHandler)

		req := httptest.NewRequest("POST", "/orders/bad-id/apply-promotion", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		_, _, handler, app := setupHandlerTest(t)
		app.Post("/orders/:id/apply-promotion", handler.ApplyPromotionHandler)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/apply-promotion", bytes.NewReader([]byte("bad")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationError", func(t *testing.T) {
		_, _, handler, app := setupHandlerTest(t)
		app.Post("/orders/:id/apply-promotion", handler.ApplyPromotionHandler)

		// Missing required promotion_id
		reqBody := orders.ApplyPromotionRequest{}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/apply-promotion", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/apply-promotion", handler.ApplyPromotionHandler)

		reqBody := orders.ApplyPromotionRequest{PromotionID: promoID}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ApplyPromotion(gomock.Any(), orderID, gomock.Any()).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/apply-promotion", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("PromotionNotApplicable", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/apply-promotion", handler.ApplyPromotionHandler)

		reqBody := orders.ApplyPromotionRequest{PromotionID: promoID}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ApplyPromotion(gomock.Any(), orderID, gomock.Any()).Return(nil,
			fmt.Errorf("%w: promotion expired", common.ErrPromotionNotApplicable),
		)

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/apply-promotion", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService, mockLogger, handler, app := setupHandlerTest(t)
		allowAllHandlerLoggerCalls(mockLogger)
		app.Post("/orders/:id/apply-promotion", handler.ApplyPromotionHandler)

		reqBody := orders.ApplyPromotionRequest{PromotionID: promoID}
		body, _ := json.Marshal(reqBody)

		mockService.EXPECT().ApplyPromotion(gomock.Any(), orderID, gomock.Any()).Return(nil, errors.New("db error"))

		req := httptest.NewRequest("POST", "/orders/"+orderID.String()+"/apply-promotion", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

// ====================== Helpers ======================

// Suppress unused imports
var _ = time.Now
