package orders_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/orders"
	orders_repo "POS-kasir/internal/orders/repository"
	products_repo "POS-kasir/internal/products/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	activitylog_repo "POS-kasir/internal/activitylog/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// setupTest creates basic mocks for tests that don't need pgxmock.
func setupTest(t *testing.T) (*mocks.MockStore, *mocks.MockOrderQuerier, *mocks.MockProductQuerier, *mocks.MockIMidtrans, *mocks.MockIActivityService, *mocks.MockFieldLogger, orders.IOrderService) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStore(ctrl)
	mockOrderRepo := mocks.NewMockOrderQuerier(ctrl)
	mockProductRepo := mocks.NewMockProductQuerier(ctrl)
	mockMidtrans := mocks.NewMockIMidtrans(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	service := orders.NewOrderService(mockStore, mockOrderRepo, mockProductRepo, mockMidtrans, mockActivity, mockLogger)
	return mockStore, mockOrderRepo, mockProductRepo, mockMidtrans, mockActivity, mockLogger, service
}

// allowAllLoggerCalls sets up AnyTimes expectations for all logger methods
// to prevent strict mock failures from variadic argument count mismatches.
func allowAllLoggerCalls(mockLogger *mocks.MockFieldLogger) {
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

// setupTestWithPgxMock creates mocks including a pgxmock pool for transaction testing.
func setupTestWithPgxMock(t *testing.T) (pgxmock.PgxPoolIface, *mocks.MockStore, *mocks.MockOrderQuerier, *mocks.MockProductQuerier, *mocks.MockIMidtrans, *mocks.MockIActivityService, *mocks.MockFieldLogger, orders.IOrderService) {
	ctrl := gomock.NewController(t)
	mockStore := mocks.NewMockStore(ctrl)
	mockOrderRepo := mocks.NewMockOrderQuerier(ctrl)
	mockProductRepo := mocks.NewMockProductQuerier(ctrl)
	mockMidtrans := mocks.NewMockIMidtrans(ctrl)
	mockActivity := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)

	mockPgx, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}

	service := orders.NewOrderService(mockStore, mockOrderRepo, mockProductRepo, mockMidtrans, mockActivity, mockLogger)
	return mockPgx, mockStore, mockOrderRepo, mockProductRepo, mockMidtrans, mockActivity, mockLogger, service
}

func TestOrderService_GetOrder(t *testing.T) {
	_, mockRepo, _, _, _, mockLogger, service := setupTest(t)
	ctx := context.Background()
	orderID := uuid.New()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		repoOrder := orders_repo.GetOrderWithDetailsRow{
			ID:         orderID,
			UserID:     pgtype.UUID{Bytes: uuid.New(), Valid: true},
			Status:     orders_repo.OrderStatusOpen,
			GrossTotal: 10000,
			NetTotal:   10000,
			CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
			Items:      nil,
		}

		mockRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(repoOrder, nil)

		resp, err := service.GetOrder(ctx, orderID)

		assert.NoError(t, err)
		assert.Equal(t, orderID, resp.ID)
		assert.Equal(t, orders_repo.OrderStatusOpen, resp.Status)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{}, pgx.ErrNoRows)
		mockLogger.EXPECT().Warn(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.GetOrder(ctx, orderID)

		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{}, errors.New("db error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

		resp, err := service.GetOrder(ctx, orderID)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOrderService_CreateOrder(t *testing.T) {
	mockPgx, mockStore, _, _, _, mockActivity, mockLogger, service := setupTestWithPgxMock(t)
	defer mockPgx.Close()

	productID := uuid.New()
	newOrderID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	// 18-column Order row template (reused for CreateOrder and UpdateOrderTotals)
	orderColumns := []string{
		"id", "user_id", "type", "status", "created_at", "updated_at",
		"gross_total", "discount_amount", "net_total", "applied_promotion_id",
		"payment_method_id", "payment_gateway_reference", "cash_received", "change_due",
		"cancellation_reason_id", "cancellation_notes", "payment_url", "payment_token",
	}

	// 19-column GetOrderWithDetails row (18 + items)
	orderWithDetailsColumns := append(append([]string{}, orderColumns...), "items")

	makeOrderRow := func(grossTotal, netTotal int64) []interface{} {
		return []interface{}{
			newOrderID, pgtype.UUID{Bytes: userID, Valid: true},
			orders_repo.OrderTypeDineIn, orders_repo.OrderStatusOpen,
			pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true},
			grossTotal, int64(0), netTotal, pgtype.UUID{},
			nil, nil, nil, nil, nil, nil, nil, nil,
		}
	}

	t.Run("Success", func(t *testing.T) {
		allowAllLoggerCalls(mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, common.UserIDKey, userID)

		req := orders.CreateOrderRequest{
			Type: orders_repo.OrderTypeDineIn,
			Items: []orders.CreateOrderItemRequest{
				{
					ProductID: productID,
					Quantity:  1,
					Options:   nil,
				},
			},
		}

		// Mock ExecTx: execute the callback with pgxmock, capture error
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(pgx.Tx) error) error {
				return fn(mockPgx)
			},
		)

		// 1. CreateOrder (INSERT INTO orders)
		mockPgx.ExpectQuery("INSERT INTO orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makeOrderRow(0, 0)...))

		// 2. GetProductsForUpdate (SELECT ... FOR UPDATE)
		mockPgx.ExpectQuery("SELECT .* FROM products WHERE id = ANY").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "name", "category_id", "image_url", "price", "stock",
				"created_at", "updated_at", "deleted_at", "cost_price",
			}).AddRow(
				productID, "Test Product", nil, nil, int64(10000), int32(10),
				now, now, nil, pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 3. BatchCreateOrderItems (INSERT INTO order_items)
		mockPgx.ExpectQuery("INSERT INTO order_items").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "order_id", "product_id", "quantity", "price_at_sale",
				"subtotal", "discount_amount", "net_subtotal", "cost_price_at_sale",
			}).AddRow(
				uuid.New(), newOrderID, productID, int32(1), int64(10000),
				int64(10000), int64(0), int64(10000),
				pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 4. BatchDecreaseProductStock (UPDATE products)
		mockPgx.ExpectExec("UPDATE products AS p SET stock").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// 5. CreateStockHistory (INSERT INTO stock_history)
		mockPgx.ExpectQuery("INSERT INTO stock_history").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "product_id", "change_amount", "previous_stock", "current_stock",
				"change_type", "reference_id", "note", "created_by", "created_at",
			}).AddRow(
				uuid.New(), productID, int32(-1), int32(10), int32(9),
				orders_repo.StockChangeTypeSale,
				pgtype.UUID{Bytes: newOrderID, Valid: true},
				utils.StringPtr("Order Created"),
				pgtype.UUID{Bytes: userID, Valid: true},
				now,
			))

		// 6. UpdateOrderTotals (UPDATE orders) - returns full 18-column Order
		mockPgx.ExpectQuery("UPDATE orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makeOrderRow(10000, 10000)...))

		// 7. GetOrderWithDetails (SELECT ... FROM orders o WHERE o.id) - Items is nil to bypass JSON unmarshal
		mockPgx.ExpectQuery("SELECT .* FROM orders o").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderWithDetailsColumns).AddRow(
				append(makeOrderRow(10000, 10000), nil)...,
			))

		// Activity Log (returns nothing)
		mockActivity.EXPECT().Log(gomock.Any(), userID, activitylog_repo.LogActionTypeCREATE, activitylog_repo.LogEntityTypeORDER, newOrderID.String(), gomock.Any())

		resp, err := service.CreateOrder(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, newOrderID, resp.ID)
		assert.Equal(t, orders_repo.OrderStatusOpen, resp.Status)
		assert.Equal(t, int64(10000), resp.GrossTotal)
		assert.NoError(t, mockPgx.ExpectationsWereMet())
	})

	t.Run("TransactionError", func(t *testing.T) {
		allowAllLoggerCalls(mockLogger)

		ctx := context.Background()
		ctx = context.WithValue(ctx, common.UserIDKey, userID)

		req := orders.CreateOrderRequest{
			Type: orders_repo.OrderTypeDineIn,
			Items: []orders.CreateOrderItemRequest{
				{ProductID: productID, Quantity: 1},
			},
		}

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(errors.New("transaction failed"))

		resp, err := service.CreateOrder(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "transaction failed")
	})
}

func TestOrderService_InitiateMidtransPayment(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	txnID := "midtrans-txn-123"

	baseOrder := orders_repo.GetOrderWithDetailsRow{
		ID:         orderID,
		UserID:     pgtype.UUID{Bytes: userID, Valid: true},
		Type:       orders_repo.OrderTypeDineIn,
		Status:     orders_repo.OrderStatusOpen,
		GrossTotal: 25000,
		NetTotal:   25000,
		CreatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Time: now, Valid: true},
		Items:      nil,
	}

	t.Run("Success", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, mockActivity, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		// Order has no existing payment reference
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)

		// Midtrans charge succeeds
		mockMidtrans.EXPECT().CreateQRISCharge(orderID.String(), int64(25000)).Return(&coreapi.ChargeResponse{
			TransactionID: txnID,
			OrderID:       orderID.String(),
			GrossAmount:   "25000.00",
			QRString:      "qris-string-data",
			ExpiryTime:    "2026-02-18 12:00:00",
			Actions: []coreapi.Action{
				{Name: "generate-qr-code", Method: "GET", URL: "https://api.midtrans.com/qr/123"},
			},
		}, nil)

		// Update payment info and URL
		mockOrderRepo.EXPECT().UpdateOrderPaymentInfo(ctx, gomock.Any()).Return(nil)
		mockOrderRepo.EXPECT().UpdateOrderPaymentUrl(ctx, gomock.Any()).Return(nil)

		// Activity log
		mockActivity.EXPECT().Log(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		resp, err := service.InitiateMidtransPayment(ctx, orderID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, orderID.String(), resp.OrderID)
		assert.Equal(t, txnID, resp.TransactionID)
		assert.Equal(t, "25000.00", resp.GrossAmount)
		assert.Equal(t, "qris-string-data", resp.QRString)
		assert.Len(t, resp.Actions, 1)
		assert.Equal(t, "generate-qr-code", resp.Actions[0].Name)
	})

	t.Run("AlreadyHasPaymentReference", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.Background()

		payRef := "existing-txn-456"
		actionsJSON := `[{"name":"generate-qr-code","method":"GET","url":"https://api.midtrans.com/qr/456"}]`

		orderWithPayment := baseOrder
		orderWithPayment.PaymentGatewayReference = &payRef
		orderWithPayment.PaymentUrl = &actionsJSON

		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orderWithPayment, nil)

		resp, err := service.InitiateMidtransPayment(ctx, orderID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, orderID.String(), resp.OrderID)
		assert.Equal(t, "existing-txn-456", resp.TransactionID)
		assert.Len(t, resp.Actions, 1)
		assert.Equal(t, "https://api.midtrans.com/qr/456", resp.Actions[0].URL)
	})

	t.Run("GetOrderError", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, _, service := setupTest(t)
		ctx := context.Background()

		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{}, errors.New("db error"))

		resp, err := service.InitiateMidtransPayment(ctx, orderID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "db error")
	})

	t.Run("CreateQRISChargeError", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.Background()

		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)
		mockMidtrans.EXPECT().CreateQRISCharge(orderID.String(), int64(25000)).Return(nil, errors.New("midtrans unavailable"))

		resp, err := service.InitiateMidtransPayment(ctx, orderID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "midtrans unavailable")
	})

	t.Run("UpdatePaymentInfoError", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.Background()

		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)
		mockMidtrans.EXPECT().CreateQRISCharge(orderID.String(), int64(25000)).Return(&coreapi.ChargeResponse{
			TransactionID: txnID,
			OrderID:       orderID.String(),
			GrossAmount:   "25000.00",
			Actions:       []coreapi.Action{},
		}, nil)
		mockOrderRepo.EXPECT().UpdateOrderPaymentInfo(ctx, gomock.Any()).Return(errors.New("update failed"))

		resp, err := service.InitiateMidtransPayment(ctx, orderID)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "update failed")
	})
}

func TestOrderService_HandleMidtransNotification(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	now := time.Now()
	txnID := "midtrans-txn-789"

	basePayload := payment.MidtransNotificationPayload{
		OrderID:           orderID.String(),
		TransactionID:     txnID,
		TransactionStatus: "settlement",
		StatusCode:        "200",
		SignatureKey:      "valid-signature",
	}

	baseOrder := orders_repo.GetOrderWithDetailsRow{
		ID:        orderID,
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		Status:    orders_repo.OrderStatusOpen,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		Items:     nil,
	}

	updatedOrder := orders_repo.Order{
		ID:        orderID,
		UserID:    pgtype.UUID{Bytes: userID, Valid: true},
		Status:    orders_repo.OrderStatusPaid,
		CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
	}

	t.Run("SettlementSuccess", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, mockActivity, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		mockMidtrans.EXPECT().VerifyNotificationSignature(basePayload).Return(nil)
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)
		mockOrderRepo.EXPECT().UpdateOrderStatusByGatewayRef(ctx, orders_repo.UpdateOrderStatusByGatewayRefParams{
			PaymentGatewayReference: &txnID,
			Status:                  orders_repo.OrderStatusPaid,
		}).Return(updatedOrder, nil)
		mockActivity.EXPECT().Log(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		err := service.HandleMidtransNotification(ctx, basePayload)

		assert.NoError(t, err)
	})

	t.Run("CaptureSuccess", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, mockActivity, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		capturePayload := basePayload
		capturePayload.TransactionStatus = "capture"

		mockMidtrans.EXPECT().VerifyNotificationSignature(capturePayload).Return(nil)
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)
		mockOrderRepo.EXPECT().UpdateOrderStatusByGatewayRef(ctx, gomock.Any()).Return(updatedOrder, nil)
		mockActivity.EXPECT().Log(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		err := service.HandleMidtransNotification(ctx, capturePayload)

		assert.NoError(t, err)
	})

	t.Run("CancelExpire", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, mockActivity, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		cancelPayload := basePayload
		cancelPayload.TransactionStatus = "expire"

		cancelledOrder := updatedOrder
		cancelledOrder.Status = orders_repo.OrderStatusCancelled

		mockMidtrans.EXPECT().VerifyNotificationSignature(cancelPayload).Return(nil)
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)
		mockOrderRepo.EXPECT().UpdateOrderStatusByGatewayRef(ctx, gomock.Any()).Return(cancelledOrder, nil)
		mockActivity.EXPECT().Log(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		err := service.HandleMidtransNotification(ctx, cancelPayload)

		assert.NoError(t, err)
	})

	t.Run("SignatureVerificationFailed", func(t *testing.T) {
		_, _, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		mockMidtrans.EXPECT().VerifyNotificationSignature(basePayload).Return(errors.New("invalid signature"))

		err := service.HandleMidtransNotification(ctx, basePayload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "signature verification failed")
	})

	t.Run("InvalidOrderID", func(t *testing.T) {
		_, _, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		invalidPayload := basePayload
		invalidPayload.OrderID = "not-a-uuid"

		mockMidtrans.EXPECT().VerifyNotificationSignature(invalidPayload).Return(nil)

		err := service.HandleMidtransNotification(ctx, invalidPayload)

		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		mockMidtrans.EXPECT().VerifyNotificationSignature(basePayload).Return(nil)
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{}, pgx.ErrNoRows)

		err := service.HandleMidtransNotification(ctx, basePayload)

		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("AlreadyFinalizedOrder", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		paidOrder := baseOrder
		paidOrder.Status = orders_repo.OrderStatusPaid

		mockMidtrans.EXPECT().VerifyNotificationSignature(basePayload).Return(nil)
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(paidOrder, nil)

		err := service.HandleMidtransNotification(ctx, basePayload)

		assert.NoError(t, err) // Should return nil (idempotent)
	})

	t.Run("UnknownStatusIgnored", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		pendingPayload := basePayload
		pendingPayload.TransactionStatus = "pending"

		mockMidtrans.EXPECT().VerifyNotificationSignature(pendingPayload).Return(nil)
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)

		err := service.HandleMidtransNotification(ctx, pendingPayload)

		assert.NoError(t, err) // Should return nil (ignored)
	})

	t.Run("UpdateStatusError", func(t *testing.T) {
		_, mockOrderRepo, _, mockMidtrans, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		mockMidtrans.EXPECT().VerifyNotificationSignature(basePayload).Return(nil)
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(baseOrder, nil)
		mockOrderRepo.EXPECT().UpdateOrderStatusByGatewayRef(ctx, gomock.Any()).Return(orders_repo.Order{}, errors.New("db error"))

		err := service.HandleMidtransNotification(ctx, basePayload)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestOrderService_ListOrders(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		_, mockOrderRepo, mockProductRepo, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		page := 1
		limit := 10
		req := orders.ListOrdersRequest{
			Page:  &page,
			Limit: &limit,
		}

		// ListOrders and CountOrders are called concurrently via goroutines,
		// so we use gomock.Any() for context matching
		payMethodID := int32(1)
		mockOrderRepo.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return([]orders_repo.ListOrdersRow{
			{
				ID:              orderID,
				UserID:          pgtype.UUID{Bytes: userID, Valid: true},
				Type:            orders_repo.OrderTypeDineIn,
				Status:          orders_repo.OrderStatusOpen,
				GrossTotal:      15000,
				NetTotal:        15000,
				CreatedAt:       pgtype.Timestamptz{Time: now, Valid: true},
				PaymentMethodID: &payMethodID,
			},
		}, nil)
		mockOrderRepo.EXPECT().CountOrders(gomock.Any(), gomock.Any()).Return(int64(1), nil)

		// GetOrderItemsByOrderID for the order
		mockOrderRepo.EXPECT().GetOrderItemsByOrderID(ctx, orderID).Return([]orders_repo.OrderItem{
			{
				ID:          uuid.New(),
				OrderID:     orderID,
				ProductID:   productID,
				Quantity:    2,
				PriceAtSale: 7500,
				Subtotal:    15000,
			},
		}, nil)

		// GetProductsByIDs for product names
		mockProductRepo.EXPECT().GetProductsByIDs(ctx, []uuid.UUID{productID}).Return([]products_repo.Product{
			{
				ID:   productID,
				Name: "Nasi Goreng",
			},
		}, nil)

		resp, err := service.ListOrders(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Len(t, resp.Orders, 1)
		assert.Equal(t, orderID, resp.Orders[0].ID)
		assert.Equal(t, orders_repo.OrderStatusOpen, resp.Orders[0].Status)
		assert.Equal(t, int64(15000), resp.Orders[0].NetTotal)
		assert.True(t, resp.Orders[0].IsPaid)
		assert.Len(t, resp.Orders[0].Items, 1)
		assert.Equal(t, "Nasi Goreng", resp.Orders[0].Items[0].ProductName)
		assert.Equal(t, 1, resp.Pagination.CurrentPage)
	})

	t.Run("EmptyList", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		page := 1
		limit := 10
		req := orders.ListOrdersRequest{
			Page:  &page,
			Limit: &limit,
		}

		mockOrderRepo.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return([]orders_repo.ListOrdersRow{}, nil)
		mockOrderRepo.EXPECT().CountOrders(gomock.Any(), gomock.Any()).Return(int64(0), nil)

		resp, err := service.ListOrders(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Empty(t, resp.Orders)
		assert.Equal(t, 0, resp.Pagination.TotalData)
	})

	t.Run("ListError", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		req := orders.ListOrdersRequest{}

		mockOrderRepo.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return(nil, errors.New("list failed"))
		mockOrderRepo.EXPECT().CountOrders(gomock.Any(), gomock.Any()).Return(int64(0), nil)

		resp, err := service.ListOrders(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("CountError", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		req := orders.ListOrdersRequest{}

		mockOrderRepo.EXPECT().ListOrders(gomock.Any(), gomock.Any()).Return([]orders_repo.ListOrdersRow{}, nil)
		mockOrderRepo.EXPECT().CountOrders(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("count failed"))

		resp, err := service.ListOrders(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})
}

func TestOrderService_CancelOrder(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	req := orders.CancelOrderRequest{
		CancellationReasonID: 1,
		CancellationNotes:    "Customer changed mind",
	}

	t.Run("Success", func(t *testing.T) {
		mockPgx, mockStore, _, _, _, mockActivity, mockLogger, service := setupTestWithPgxMock(t)
		defer mockPgx.Close()
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		orderColumns := []string{
			"id", "user_id", "type", "status", "created_at", "updated_at",
			"gross_total", "discount_amount", "net_total", "applied_promotion_id",
			"payment_method_id", "payment_gateway_reference", "cash_received", "change_due",
			"cancellation_reason_id", "cancellation_notes", "payment_url", "payment_token",
		}
		orderWithDetailsColumns := append(append([]string{}, orderColumns...), "items")
		now := time.Now()

		productID := uuid.New()
		itemID := uuid.New()

		// Build items JSON for GetOrderWithDetails to return
		items := []orders_repo.OrderItem{
			{ID: itemID, OrderID: orderID, ProductID: productID, Quantity: 2, PriceAtSale: 10000, Subtotal: 20000},
		}
		itemsJSON, _ := json.Marshal(items)

		// ExecTx: execute callback with pgxmock
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(pgx.Tx) error) error {
				return fn(mockPgx)
			},
		)

		// 1. GetOrderWithDetails (status=open, with items)
		mockPgx.ExpectQuery("SELECT .* FROM orders o").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderWithDetailsColumns).AddRow(
				orderID, pgtype.UUID{Bytes: userID, Valid: true},
				orders_repo.OrderTypeDineIn, orders_repo.OrderStatusOpen,
				pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true},
				int64(20000), int64(0), int64(20000), pgtype.UUID{},
				nil, nil, nil, nil, nil, nil, nil, nil,
				itemsJSON,
			))

		// 2. CancelOrder (UPDATE orders SET status='cancelled')
		mockPgx.ExpectQuery("UPDATE orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(
				orderID, pgtype.UUID{Bytes: userID, Valid: true},
				orders_repo.OrderTypeDineIn, orders_repo.OrderStatusCancelled,
				pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true},
				int64(20000), int64(0), int64(20000), pgtype.UUID{},
				nil, nil, nil, nil, nil, nil, nil, nil,
			))

		// 3. For each item: GetProductByID (from products_repo.New(tx) - has join with options, 11 cols)
		mockPgx.ExpectQuery("SELECT .* FROM products p WHERE").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "name", "category_id", "image_url", "price", "stock",
				"created_at", "updated_at", "deleted_at", "cost_price", "options",
			}).AddRow(
				productID, "Test Product", nil, nil, int64(10000), int32(8),
				now, now, nil, pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
				nil,
			))

		// 4. AddProductStock (from products_repo.New(tx) - returns full Product, 10 cols)
		mockPgx.ExpectQuery("UPDATE products").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "name", "category_id", "image_url", "price", "stock",
				"created_at", "updated_at", "deleted_at", "cost_price",
			}).AddRow(
				productID, "Test Product", nil, nil, int64(10000), int32(10),
				now, now, nil, pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 5. CreateStockHistory
		mockPgx.ExpectQuery("INSERT INTO stock_history").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "product_id", "change_amount", "previous_stock", "current_stock",
				"change_type", "reference_id", "note", "created_by", "created_at",
			}).AddRow(
				uuid.New(), productID, int32(2), int32(8), int32(10),
				orders_repo.StockChangeTypeReturn,
				pgtype.UUID{Bytes: orderID, Valid: true},
				utils.StringPtr("Order Cancelled"),
				pgtype.UUID{Bytes: userID, Valid: true},
				now,
			))

		// Activity log after successful cancellation
		mockActivity.EXPECT().Log(
			gomock.Any(),
			userID,
			activitylog_repo.LogActionTypeCANCEL,
			activitylog_repo.LogEntityTypeORDER,
			orderID.String(),
			gomock.Any(),
		)

		err := service.CancelOrder(ctx, orderID, req)

		assert.NoError(t, err)
		assert.NoError(t, mockPgx.ExpectationsWereMet())
	})

	t.Run("TransactionError", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(errors.New("transaction failed"))

		err := service.CancelOrder(ctx, orderID, req)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transaction failed")
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		// ExecTx callback returns ErrNotFound
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(common.ErrNotFound)

		err := service.CancelOrder(ctx, orderID, req)

		assert.ErrorIs(t, err, common.ErrNotFound)
	})

	t.Run("OrderNotCancellable", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		// ExecTx callback returns ErrOrderNotCancellable (order not in 'open' status)
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(common.ErrOrderNotCancellable)

		err := service.CancelOrder(ctx, orderID, req)

		assert.ErrorIs(t, err, common.ErrOrderNotCancellable)
	})
}

func TestOrderService_UpdateOrderItems(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	productID := uuid.New()

	reqs := []orders.UpdateOrderItemRequest{
		{
			ProductID: productID,
			Quantity:  3,
		},
	}

	t.Run("Success", func(t *testing.T) {
		mockPgx, mockStore, _, mockProductRepo, _, mockActivity, mockLogger, service := setupTestWithPgxMock(t)
		defer mockPgx.Close()
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		now := time.Now()
		orderColumns := []string{
			"id", "user_id", "type", "status", "created_at", "updated_at",
			"gross_total", "discount_amount", "net_total", "applied_promotion_id",
			"payment_method_id", "payment_gateway_reference", "cash_received", "change_due",
			"cancellation_reason_id", "cancellation_notes", "payment_url", "payment_token",
		}
		orderWithDetailsColumns := append(append([]string{}, orderColumns...), "items")

		existingItemID := uuid.New()

		// Build items: marshal OrderItem to JSON, then unmarshal to []interface{}
		// This matches how pgx scans JSON columns into interface{}
		rawItems := []orders_repo.OrderItem{
			{ID: existingItemID, OrderID: orderID, ProductID: productID, Quantity: 3, PriceAtSale: 10000, Subtotal: 30000, DiscountAmount: 0, NetSubtotal: 30000, CostPriceAtSale: pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true}},
		}
		rawJSON, _ := json.Marshal(rawItems)
		var itemsForPgxMock interface{}
		json.Unmarshal(rawJSON, &itemsForPgxMock)

		makeOrderRow := func(grossTotal, netTotal int64) []interface{} {
			return []interface{}{
				orderID, pgtype.UUID{Bytes: userID, Valid: true},
				orders_repo.OrderTypeDineIn, orders_repo.OrderStatusOpen,
				pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true},
				grossTotal, int64(0), netTotal, pgtype.UUID{},
				nil, nil, nil, nil, nil, nil, nil, nil,
			}
		}

		// ExecTx: execute callback with pgxmock
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(pgx.Tx) error) error {
				return fn(mockPgx)
			},
		)

		// 1. GetOrderForUpdate
		mockPgx.ExpectQuery("SELECT .* FROM orders").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makeOrderRow(10000, 10000)...))

		// 2. GetOrderItemsByOrderID - returns existing item (same product, qty=1)
		mockPgx.ExpectQuery("SELECT .* FROM order_items WHERE order_id").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "order_id", "product_id", "quantity", "price_at_sale",
				"subtotal", "discount_amount", "net_subtotal", "cost_price_at_sale",
			}).AddRow(
				existingItemID, orderID, productID, int32(1), int64(10000),
				int64(10000), int64(0), int64(10000),
				pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 3. GetProductByID (for the item update)
		mockPgx.ExpectQuery("SELECT .* FROM products WHERE id").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "name", "category_id", "image_url", "price", "stock",
				"created_at", "updated_at", "deleted_at", "cost_price",
			}).AddRow(
				productID, "Test Product", nil, nil, int64(10000), int32(10),
				now, now, nil, pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 4. qtyDiff=2 (3-1) > 0 → DecreaseProductStock (products_repo uses QueryRow, returns 10-col Product)
		mockPgx.ExpectQuery("UPDATE products").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "name", "category_id", "image_url", "price", "stock",
				"created_at", "updated_at", "deleted_at", "cost_price",
			}).AddRow(
				productID, "Test Product", nil, nil, int64(10000), int32(8),
				now, now, nil, pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 5. CreateStockHistory
		mockPgx.ExpectQuery("INSERT INTO stock_history").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "product_id", "change_amount", "previous_stock", "current_stock",
				"change_type", "reference_id", "note", "created_by", "created_at",
			}).AddRow(
				uuid.New(), productID, int32(-2), int32(10), int32(8),
				orders_repo.StockChangeTypeSale,
				pgtype.UUID{Bytes: orderID, Valid: true},
				utils.StringPtr("Order Item Qty Increase"),
				pgtype.UUID{Bytes: userID, Valid: true},
				now,
			))

		// 6. UpdateOrderItemQuantity
		mockPgx.ExpectQuery("UPDATE order_items").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "order_id", "product_id", "quantity", "price_at_sale",
				"subtotal", "discount_amount", "net_subtotal", "cost_price_at_sale",
			}).AddRow(
				existingItemID, orderID, productID, int32(3), int64(10000),
				int64(30000), int64(0), int64(30000),
				pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 7. UpdateOrderTotals
		mockPgx.ExpectQuery("UPDATE orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makeOrderRow(30000, 30000)...))

		// 8. GetOrderWithDetails (final, with items for buildOrderDetailResponseFromQueryResult)
		mockPgx.ExpectQuery("SELECT .* FROM orders o").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderWithDetailsColumns).AddRow(
				append(makeOrderRow(30000, 30000), itemsForPgxMock)...,
			))

		// Activity log
		mockActivity.EXPECT().Log(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		// buildOrderDetailResponseFromQueryResult calls GetProductsByIDs for item names
		mockProductRepo.EXPECT().GetProductsByIDs(gomock.Any(), gomock.Any()).Return([]products_repo.Product{
			{ID: productID, Name: "Test Product"},
		}, nil)

		resp, err := service.UpdateOrderItems(ctx, orderID, reqs)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, orderID, resp.ID)
		assert.NoError(t, mockPgx.ExpectationsWereMet())
	})

	t.Run("TransactionError", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(errors.New("transaction failed"))

		resp, err := service.UpdateOrderItems(ctx, orderID, reqs)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "transaction failed")
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(common.ErrNotFound)

		resp, err := service.UpdateOrderItems(ctx, orderID, reqs)

		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("OrderNotModifiable", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)

		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(common.ErrOrderNotModifiable)

		resp, err := service.UpdateOrderItems(ctx, orderID, reqs)

		assert.ErrorIs(t, err, common.ErrOrderNotModifiable)
		assert.Nil(t, resp)
	})
}

func TestOrderService_ConfirmManualPayment(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()

	req := orders.ConfirmManualPaymentRequest{
		PaymentMethodID: 1,
		CashReceived:    50000,
	}

	t.Run("Success", func(t *testing.T) {
		mockPgx, mockStore, _, mockProductRepo, _, mockActivity, mockLogger, service := setupTestWithPgxMock(t)
		defer mockPgx.Close()
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		now := time.Now()
		orderColumns := []string{
			"id", "user_id", "type", "status", "created_at", "updated_at",
			"gross_total", "discount_amount", "net_total", "applied_promotion_id",
			"payment_method_id", "payment_gateway_reference", "cash_received", "change_due",
			"cancellation_reason_id", "cancellation_notes", "payment_url", "payment_token",
		}
		orderWithDetailsColumns := append(append([]string{}, orderColumns...), "items")

		productID := uuid.New()
		paymentMethodID := int32(1)
		cashReceived := int64(50000)
		changeDue := int64(10000)

		// Build items: marshal OrderItem to JSON, then unmarshal to []interface{}
		rawItems := []orders_repo.OrderItem{
			{ID: uuid.New(), OrderID: orderID, ProductID: productID, Quantity: 2, PriceAtSale: 20000, Subtotal: 40000, DiscountAmount: 0, NetSubtotal: 40000, CostPriceAtSale: pgtype.Numeric{Int: big.NewInt(10000), Exp: 0, Valid: true}},
		}
		rawJSON, _ := json.Marshal(rawItems)
		var itemsForPgxMock interface{}
		json.Unmarshal(rawJSON, &itemsForPgxMock)

		makeOrderRow := func() []interface{} {
			return []interface{}{
				orderID, pgtype.UUID{Bytes: userID, Valid: true},
				orders_repo.OrderTypeDineIn, orders_repo.OrderStatusOpen,
				pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true},
				int64(40000), int64(0), int64(40000), pgtype.UUID{},
				nil, nil, nil, nil, nil, nil, nil, nil,
			}
		}

		makePaidOrderRow := func() []interface{} {
			return []interface{}{
				orderID, pgtype.UUID{Bytes: userID, Valid: true},
				orders_repo.OrderTypeDineIn, orders_repo.OrderStatusPaid,
				pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true},
				int64(40000), int64(0), int64(40000), pgtype.UUID{},
				&paymentMethodID, nil, &cashReceived, &changeDue, nil, nil, nil, nil,
			}
		}

		// ExecTx: execute callback with pgxmock
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(pgx.Tx) error) error {
				return fn(mockPgx)
			},
		)

		// 1. GetOrderForUpdate
		mockPgx.ExpectQuery("SELECT .* FROM orders").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makeOrderRow()...))

		// 2. UpdateOrderManualPayment
		mockPgx.ExpectQuery("UPDATE orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makePaidOrderRow()...))

		// 3. GetOrderWithDetails (final, with items)
		mockPgx.ExpectQuery("SELECT .* FROM orders o").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderWithDetailsColumns).AddRow(
				append(makePaidOrderRow(), itemsForPgxMock)...,
			))

		// Activity log after successful payment
		mockActivity.EXPECT().Log(
			gomock.Any(),
			userID,
			activitylog_repo.LogActionTypePROCESSPAYMENT,
			activitylog_repo.LogEntityTypeORDER,
			orderID.String(),
			gomock.Any(),
		)

		// buildOrderDetailResponseFromQueryResult calls GetProductsByIDs
		mockProductRepo.EXPECT().GetProductsByIDs(gomock.Any(), gomock.Any()).Return([]products_repo.Product{
			{ID: productID, Name: "Test Product"},
		}, nil)

		resp, err := service.ConfirmManualPayment(ctx, orderID, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, orderID, resp.ID)
		assert.NoError(t, mockPgx.ExpectationsWereMet())
	})

	t.Run("TransactionError", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(errors.New("transaction failed"))

		resp, err := service.ConfirmManualPayment(ctx, orderID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "transaction failed")
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(common.ErrNotFound)

		resp, err := service.ConfirmManualPayment(ctx, orderID, req)

		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("OrderNotModifiable", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		// Order is cancelled
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(common.ErrOrderNotModifiable)

		resp, err := service.ConfirmManualPayment(ctx, orderID, req)

		assert.ErrorIs(t, err, common.ErrOrderNotModifiable)
		assert.Nil(t, resp)
	})
}

func TestOrderService_UpdateOperationalStatus(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	t.Run("Success", func(t *testing.T) {
		_, mockOrderRepo, _, _, mockActivity, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		req := orders.UpdateOrderStatusRequest{
			Status: orders_repo.OrderStatusInProgress,
		}

		// GetOrderWithDetails for status check
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{
			ID:        orderID,
			UserID:    pgtype.UUID{Bytes: userID, Valid: true},
			Status:    orders_repo.OrderStatusOpen,
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			Items:     nil,
		}, nil)

		// UpdateOrderStatus
		mockOrderRepo.EXPECT().UpdateOrderStatus(ctx, orders_repo.UpdateOrderStatusParams{
			ID:     orderID,
			Status: orders_repo.OrderStatusInProgress,
		}).Return(orders_repo.Order{}, nil)

		// Activity log
		mockActivity.EXPECT().Log(gomock.Any(), userID, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())

		// GetOrder (called by s.GetOrder at end) — returns the updated order
		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{
			ID:        orderID,
			UserID:    pgtype.UUID{Bytes: userID, Valid: true},
			Status:    orders_repo.OrderStatusInProgress,
			CreatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			UpdatedAt: pgtype.Timestamptz{Time: now, Valid: true},
			Items:     nil,
		}, nil)

		resp, err := service.UpdateOperationalStatus(ctx, orderID, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, orders_repo.OrderStatusInProgress, resp.Status)
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		req := orders.UpdateOrderStatusRequest{
			Status: orders_repo.OrderStatusInProgress,
		}

		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(
			orders_repo.GetOrderWithDetailsRow{}, pgx.ErrNoRows,
		)

		resp, err := service.UpdateOperationalStatus(ctx, orderID, req)

		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("InvalidTransition", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		// Cancelled → InProgress is NOT in the allowed transitions
		req := orders.UpdateOrderStatusRequest{
			Status: orders_repo.OrderStatusInProgress,
		}

		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{
			ID:     orderID,
			Status: orders_repo.OrderStatusCancelled,
			Items:  nil,
		}, nil)

		resp, err := service.UpdateOperationalStatus(ctx, orderID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, common.ErrInvalidStatusTransition)
	})

	t.Run("UpdateError", func(t *testing.T) {
		_, mockOrderRepo, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.Background()

		req := orders.UpdateOrderStatusRequest{
			Status: orders_repo.OrderStatusInProgress,
		}

		mockOrderRepo.EXPECT().GetOrderWithDetails(ctx, orderID).Return(orders_repo.GetOrderWithDetailsRow{
			ID:     orderID,
			Status: orders_repo.OrderStatusOpen,
			Items:  nil,
		}, nil)
		mockOrderRepo.EXPECT().UpdateOrderStatus(ctx, gomock.Any()).Return(orders_repo.Order{}, errors.New("db error"))

		resp, err := service.UpdateOperationalStatus(ctx, orderID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestOrderService_ApplyPromotion(t *testing.T) {
	orderID := uuid.New()
	userID := uuid.New()
	promoID := uuid.New()

	req := orders.ApplyPromotionRequest{
		PromotionID: promoID,
	}

	t.Run("Success", func(t *testing.T) {
		mockPgx, mockStore, _, mockProductRepo, _, mockActivity, mockLogger, service := setupTestWithPgxMock(t)
		defer mockPgx.Close()
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		now := time.Now()
		orderColumns := []string{
			"id", "user_id", "type", "status", "created_at", "updated_at",
			"gross_total", "discount_amount", "net_total", "applied_promotion_id",
			"payment_method_id", "payment_gateway_reference", "cash_received", "change_due",
			"cancellation_reason_id", "cancellation_notes", "payment_url", "payment_token",
		}
		orderWithDetailsColumns := append(append([]string{}, orderColumns...), "items")

		productID := uuid.New()
		itemID := uuid.New()

		makeOrderRow := func(grossTotal, discountAmount, netTotal int64) []interface{} {
			return []interface{}{
				orderID, pgtype.UUID{Bytes: userID, Valid: true},
				orders_repo.OrderTypeDineIn, orders_repo.OrderStatusOpen,
				pgtype.Timestamptz{Time: now, Valid: true}, pgtype.Timestamptz{Time: now, Valid: true},
				grossTotal, discountAmount, netTotal, pgtype.UUID{},
				nil, nil, nil, nil, nil, nil, nil, nil,
			}
		}

		// ExecTx: execute callback with pgxmock
		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).DoAndReturn(
			func(ctx context.Context, fn func(pgx.Tx) error) error {
				return fn(mockPgx)
			},
		)

		// 1. GetOrderForUpdate (order with grossTotal=50000)
		mockPgx.ExpectQuery("SELECT .* FROM orders").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makeOrderRow(50000, 0, 50000)...))

		// 2. GetOrderItemsByOrderID
		mockPgx.ExpectQuery("SELECT .* FROM order_items WHERE order_id").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "order_id", "product_id", "quantity", "price_at_sale",
				"subtotal", "discount_amount", "net_subtotal", "cost_price_at_sale",
			}).AddRow(
				itemID, orderID, productID, int32(5), int64(10000),
				int64(50000), int64(0), int64(50000),
				pgtype.Numeric{Int: big.NewInt(5000), Exp: 0, Valid: true},
			))

		// 3. GetPromotionByID (ORDER scope, percentage 10%, active)
		mockPgx.ExpectQuery("SELECT .* FROM promotions WHERE id").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "name", "description", "scope", "discount_type",
				"discount_value", "max_discount_amount", "start_date", "end_date",
				"is_active", "created_at", "updated_at", "deleted_at",
			}).AddRow(
				promoID, "10% Off", nil, orders_repo.PromotionScopeORDER, orders_repo.DiscountTypePercentage,
				pgtype.Numeric{Int: big.NewInt(10), Exp: 0, Valid: true},
				pgtype.Numeric{Int: big.NewInt(0), Exp: 0, Valid: true},
				pgtype.Timestamptz{Time: now.Add(-24 * time.Hour), Valid: true},
				pgtype.Timestamptz{Time: now.Add(24 * time.Hour), Valid: true},
				true, pgtype.Timestamptz{Time: now, Valid: true},
				pgtype.Timestamptz{Time: now, Valid: true},
				pgtype.Timestamptz{},
			))

		// 4. GetPromotionRules (empty — no rules to check)
		mockPgx.ExpectQuery("SELECT .* FROM promotion_rules WHERE promotion_id").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows([]string{
				"id", "promotion_id", "rule_type", "rule_value", "description", "created_at", "updated_at",
			}))

		// 5. UpdateOrderTotals (10% of 50000 = 5000 discount, net=45000)
		mockPgx.ExpectQuery("UPDATE orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderColumns).AddRow(makeOrderRow(50000, 5000, 45000)...))

		// 6. UpdateOrderAppliedPromotion (exec, not query)
		mockPgx.ExpectExec("UPDATE orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg()).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// 7. GetOrderWithDetails (final)
		mockPgx.ExpectQuery("SELECT .* FROM orders o").
			WithArgs(pgxmock.AnyArg()).
			WillReturnRows(pgxmock.NewRows(orderWithDetailsColumns).AddRow(
				append(makeOrderRow(50000, 5000, 45000), nil)...,
			))

		// Activity log after successful promotion application
		mockActivity.EXPECT().Log(
			gomock.Any(),
			userID,
			activitylog_repo.LogActionTypeAPPLYPROMOTION,
			activitylog_repo.LogEntityTypeORDER,
			orderID.String(),
			gomock.Any(),
		)

		// Items=nil, so buildOrderDetailResponseFromQueryResult skips item fetch
		_ = mockProductRepo

		resp, err := service.ApplyPromotion(ctx, orderID, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, orderID, resp.ID)
		assert.Equal(t, int64(45000), resp.NetTotal)
		assert.Equal(t, int64(5000), resp.DiscountAmount)
		assert.NoError(t, mockPgx.ExpectationsWereMet())
	})

	t.Run("TransactionError", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(errors.New("transaction failed"))

		resp, err := service.ApplyPromotion(ctx, orderID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("OrderNotFound", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(common.ErrNotFound)

		resp, err := service.ApplyPromotion(ctx, orderID, req)

		assert.ErrorIs(t, err, common.ErrNotFound)
		assert.Nil(t, resp)
	})

	t.Run("PromotionNotApplicable", func(t *testing.T) {
		mockStore, _, _, _, _, mockLogger, service := setupTest(t)
		allowAllLoggerCalls(mockLogger)
		ctx := context.WithValue(context.Background(), common.UserIDKey, userID)

		mockStore.EXPECT().ExecTx(gomock.Any(), gomock.Any()).Return(
			fmt.Errorf("%w: promotion is not active", common.ErrPromotionNotApplicable),
		)

		resp, err := service.ApplyPromotion(ctx, orderID, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.ErrorIs(t, err, common.ErrPromotionNotApplicable)
	})
}
