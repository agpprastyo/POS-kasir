package orders

import (
	"POS-kasir/internal/activitylog"
	activity_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/common/store"
	orders_repo "POS-kasir/internal/orders/repository"
	products_repo "POS-kasir/internal/products/repository"
	"POS-kasir/internal/settings"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/metrics"
	"strconv"
	"time"

	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	ws "POS-kasir/internal/websocket"
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderDetailResponse, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*OrderDetailResponse, error)
	InitiateMidtransPayment(ctx context.Context, orderID uuid.UUID) (*MidtransPaymentResponse, error)
	HandleMidtransNotification(ctx context.Context, payload payment.MidtransNotificationPayload) error
	ListOrders(ctx context.Context, req ListOrdersRequest) (*PagedOrderResponse, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID, req CancelOrderRequest) error
	UpdateOrderItems(ctx context.Context, orderID uuid.UUID, req UpdateOrderItemsRequest) (*OrderDetailResponse, error)
	ConfirmManualPayment(ctx context.Context, orderID uuid.UUID, req ConfirmManualPaymentRequest) (*OrderDetailResponse, error)
	UpdateOperationalStatus(ctx context.Context, orderID uuid.UUID, req UpdateOrderStatusRequest) (*OrderDetailResponse, error)
	ApplyPromotion(ctx context.Context, orderID uuid.UUID, req ApplyPromotionRequest) (*OrderDetailResponse, error)
	RefundOrder(ctx context.Context, orderID uuid.UUID, req RefundOrderRequest) (*OrderDetailResponse, error)
	CheckoutOrder(ctx context.Context, req CheckoutOrderRequest) (*OrderDetailResponse, error)
	CalculateOrder(ctx context.Context, req CalculateOrderRequest) (*CalculateOrderResponse, error)
}

// Business constants
const (
	// Payment method IDs matching the payment_methods table
	PaymentMethodCash       int32 = 1
	PaymentMethodQRIS       int32 = 2
	PaymentMethodStaticQRIS int32 = 3
)

type OrderService struct {
	store           store.Store
	ordersRepo      orders_repo.Querier
	productsRepo    products_repo.Querier
	settingsService settings.ISettingsService
	midtransService payment.IMidtrans
	activityService activitylog.IActivityService
	log             logger.ILogger
	wsHub           *ws.Hub
}

func NewOrderService(store store.Store, ordersRepo orders_repo.Querier, productsRepo products_repo.Querier, settingsService settings.ISettingsService, midtransService payment.IMidtrans, activityService activitylog.IActivityService, log logger.ILogger, wsHub *ws.Hub) IOrderService {
	return &OrderService{
		store:           store,
		ordersRepo:      ordersRepo,
		productsRepo:    productsRepo,
		settingsService: settingsService,
		midtransService: midtransService,
		activityService: activityService,
		log:             log,
		wsHub:           wsHub,
	}
}

var allowedStatusTransitions = map[orders_repo.OrderStatus]map[orders_repo.OrderStatus]bool{
	orders_repo.OrderStatusOpen: {
		orders_repo.OrderStatusInProgress: true,
		orders_repo.OrderStatusCancelled:  true,
		orders_repo.OrderStatusOpen:       true,
	},
	orders_repo.OrderStatusInProgress: {
		orders_repo.OrderStatusServed:     true,
		orders_repo.OrderStatusCancelled:  true,
		orders_repo.OrderStatusInProgress: true,
		orders_repo.OrderStatusOpen:       true, // Allow going back to open if items changed significantly?
	},
	orders_repo.OrderStatusServed: {
		orders_repo.OrderStatusPaid:       true,
		orders_repo.OrderStatusServed:     true,
		orders_repo.OrderStatusInProgress: true,
	},
	orders_repo.OrderStatusPaid: {
		orders_repo.OrderStatusPaid: true,
	},
	orders_repo.OrderStatusCancelled: {
		orders_repo.OrderStatusCancelled: true,
	},
}

func (s *OrderService) ApplyPromotion(ctx context.Context, orderID uuid.UUID, req ApplyPromotionRequest) (*OrderDetailResponse, error) {
	var finalOrder orders_repo.GetOrderWithDetailsRow

	txErr := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			return common.ErrNotFound
		}

		if order.Status != orders_repo.OrderStatusOpen {
			return common.ErrOrderNotModifiable
		}

		orderItems, err := qtx.GetOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return fmt.Errorf("failed to get order items: %w", err)
		}

		calc, err := s.calculateOrderTotals(ctx, tx, qtx, orderItems, &req.PromotionID, order.GrossTotal)
		if err != nil {
			return err
		}

		_, err = qtx.UpdateOrderTotals(ctx, orders_repo.UpdateOrderTotalsParams{
			ID:                  orderID,
			GrossTotal:          calc.GrossTotal,
			DiscountAmount:      calc.DiscountAmount,
			NetTotal:            calc.NetTotal,
			TaxAmount:           calc.TaxAmount,
			ServiceChargeAmount: calc.ServiceChargeAmount,
			Version:             order.Version,
		})
		if err != nil {
			return err
		}

		err = qtx.UpdateOrderAppliedPromotion(ctx, orders_repo.UpdateOrderAppliedPromotionParams{
			ID:                 orderID,
			AppliedPromotionID: pgtype.UUID{Bytes: req.PromotionID, Valid: true},
		})
		if err != nil {
			return fmt.Errorf("failed to update applied promotion: %w", err)
		}

		finalOrder, err = qtx.GetOrderWithDetails(ctx, orderID)
		return err
	})

	if txErr != nil {
		return nil, txErr
	}

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warnf("UpdateOrder | Actor user ID not found in context for activity logging")
	}

	logDetails := map[string]interface{}{
		"updated_order_id":     orderID,
		"updated_order_status": finalOrder.Status,
		"promotion_id":         req.PromotionID,
	}

	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypeAPPLYPROMOTION,
		activity_repo.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderUpdated, map[string]interface{}{"order_id": orderID})
	}

	return s.buildOrderDetailResponseFromQueryResult(ctx, finalOrder)
}

func (s *OrderService) UpdateOperationalStatus(ctx context.Context, orderID uuid.UUID, req UpdateOrderStatusRequest) (*OrderDetailResponse, error) {

	order, err := s.ordersRepo.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Order not found for status update", "orderID", orderID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get order for status update", "error", err)
		return nil, err
	}

	currentStatus := order.Status
	newStatus := req.Status

	if currentStatus == newStatus {
		return s.GetOrder(ctx, orderID)
	}

	allowed, ok := allowedStatusTransitions[currentStatus][newStatus]
	if !ok || !allowed {
		errMsg := fmt.Sprintf("invalid status transition from '%s' to '%s'", currentStatus, newStatus)
		s.log.Warn(errMsg, "orderID", orderID, "currentStatus", currentStatus, "newStatus", newStatus)
		return nil, fmt.Errorf("%w: %s", common.ErrInvalidStatusTransition, errMsg)
	}

	_, err = s.ordersRepo.UpdateOrderStatus(ctx, orders_repo.UpdateOrderStatusParams{
		ID:     orderID,
		Status: newStatus,
	})
	if err != nil {
		s.log.Error("Failed to update order status in repository", "error", err, "orderID", orderID)
		return nil, err
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"order_id":    orderID.String(),
		"status_from": currentStatus,
		"status_to":   newStatus,
	}
	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypeUPDATE,
		activity_repo.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderUpdated, map[string]interface{}{"order_id": orderID})
	}

	return s.GetOrder(ctx, orderID)
}

func (s *OrderService) ConfirmManualPayment(ctx context.Context, orderID uuid.UUID, req ConfirmManualPaymentRequest) (*OrderDetailResponse, error) {
	var finalOrder orders_repo.GetOrderWithDetailsRow

	txErr := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return common.ErrNotFound
			}
			return err
		}

		if order.Status == orders_repo.OrderStatusCancelled {
			return common.ErrOrderNotModifiable
		}

		if order.PaymentMethodID != nil {
			return fmt.Errorf("order already paid")
		}

		netTotal := order.NetTotal
		cashReceived := req.CashReceived

		if req.PaymentMethodID == PaymentMethodStaticQRIS {
			if cashReceived == 0 {
				cashReceived = netTotal
			}
		}

		if cashReceived < netTotal {
			return fmt.Errorf("uang kurang: tagihan %d, diterima %d", netTotal, cashReceived)
		}

		changeDue := cashReceived - netTotal

		_, err = qtx.UpdateOrderManualPayment(ctx, orders_repo.UpdateOrderManualPaymentParams{
			ID:              orderID,
			PaymentMethodID: utils.Int32Ptr(int(req.PaymentMethodID)),
			CashReceived:    &cashReceived,
			ChangeDue:       &changeDue,
			Version:         req.Version,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return common.ErrOrderConflict
			}
			return err
		}

		finalOrder, err = qtx.GetOrderWithDetails(ctx, orderID)
		return err
	})

	if txErr != nil {
		return nil, txErr
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"order_id":          orderID.String(),
		"payment_method_id": req.PaymentMethodID,
		"amount":            finalOrder.NetTotal,
	}
	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypePROCESSPAYMENT,
		activity_repo.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderUpdated, map[string]interface{}{"order_id": orderID})
	}

	// Record payment metrics
	metrics.PaymentProcessedTotal.WithLabelValues("manual", "success").Inc()
	metrics.OrderRevenueTotal.Add(float64(finalOrder.NetTotal))

	return s.buildOrderDetailResponseFromQueryResult(ctx, finalOrder)
}

func (s *OrderService) UpdateOrderItems(ctx context.Context, orderID uuid.UUID, req UpdateOrderItemsRequest) (*OrderDetailResponse, error) {
	var finalOrder orders_repo.GetOrderWithDetailsRow
	actorID, userIdOk := ctx.Value(common.UserIDKey).(uuid.UUID)

	txErr := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		qPrd := products_repo.New(tx)

		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			return common.ErrNotFound
		}
		if order.Status != orders_repo.OrderStatusOpen {
			return common.ErrOrderNotModifiable
		}

		if order.Version != req.Version {
			return common.ErrOrderConflict
		}

		existingItems, err := qtx.GetOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return err
		}

		currentMap := make(map[uuid.UUID]orders_repo.OrderItem)
		for _, item := range existingItems {
			currentMap[item.ProductID] = item
		}

		var grossTotal int64 = 0

		for _, reqItem := range req.Items {
			product, err := qtx.GetProductByID(ctx, reqItem.ProductID)
			if err != nil {
				return err
			}

			price := product.Price

			subtotal := price * int64(reqItem.Quantity)
			grossTotal += subtotal

			if existingItem, exists := currentMap[reqItem.ProductID]; exists {

				qtyDiff := reqItem.Quantity - existingItem.Quantity

				if qtyDiff > 0 {

					if product.Stock < qtyDiff {
						return fmt.Errorf("insufficient stock for update %s", product.Name)
					}
					if _, err := qPrd.DecreaseProductStock(ctx, products_repo.DecreaseProductStockParams{ID: reqItem.ProductID, Quantity: qtyDiff}); err != nil {
						return fmt.Errorf("failed to decrease stock for %s: %w", reqItem.ProductID, err)
					}

					// Log Stock Decrease
					if _, err := qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
						ProductID:     reqItem.ProductID,
						ChangeAmount:  -qtyDiff,
						PreviousStock: product.Stock,
						CurrentStock:  product.Stock - qtyDiff,
						ChangeType:    orders_repo.StockChangeTypeSale,
						ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
						Note:          utils.StringPtr("Order Item Qty Increase"),
						CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
					}); err != nil {
						return fmt.Errorf("failed to create stock history: %w", err)
					}
				} else if qtyDiff < 0 {

					restoreQty := -qtyDiff
					if _, err := qtx.AddProductStock(ctx, orders_repo.AddProductStockParams{ID: reqItem.ProductID, Stock: restoreQty}); err != nil {
						return fmt.Errorf("failed to restore stock for %s: %w", reqItem.ProductID, err)
					}

					// Log Stock Increase
					if _, err := qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
						ProductID:     reqItem.ProductID,
						ChangeAmount:  restoreQty,
						PreviousStock: product.Stock,
						CurrentStock:  product.Stock + restoreQty,
						ChangeType:    orders_repo.StockChangeTypeReturn,
						ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
						Note:          utils.StringPtr("Order Item Qty Decrease"),
						CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
					}); err != nil {
						return fmt.Errorf("failed to create stock history: %w", err)
					}
				}

				if _, err := qtx.UpdateOrderItemQuantity(ctx, orders_repo.UpdateOrderItemQuantityParams{
					ID:          existingItem.ID,
					OrderID:     orderID,
					Quantity:    reqItem.Quantity,
					Subtotal:    subtotal,
					NetSubtotal: subtotal,
				}); err != nil {
					return fmt.Errorf("failed to update order item quantity: %w", err)
				}

				delete(currentMap, reqItem.ProductID)

			} else {
				if product.Stock < reqItem.Quantity {
					return fmt.Errorf("insufficient stock for new item %s", product.Name)
				}

				if _, err := qPrd.DecreaseProductStock(ctx, products_repo.DecreaseProductStockParams{ID: reqItem.ProductID, Quantity: reqItem.Quantity}); err != nil {
					return fmt.Errorf("failed to decrease stock for new item %s: %w", reqItem.ProductID, err)
				}

				// Log Stock Decrease (New Item)
				if _, err := qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
					ProductID:     reqItem.ProductID,
					ChangeAmount:  -reqItem.Quantity,
					PreviousStock: product.Stock,
					CurrentStock:  product.Stock - reqItem.Quantity,
					ChangeType:    orders_repo.StockChangeTypeSale,
					ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
					Note:          utils.StringPtr("Order Item Added"),
					CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
				}); err != nil {
					return fmt.Errorf("failed to create stock history: %w", err)
				}

				costPrice := 0.0
				if product.CostPrice.Valid {
					f, _ := product.CostPrice.Float64Value()
					costPrice = f.Float64
				}
				numericCost := pgtype.Numeric{}
				numericCost.Scan(fmt.Sprintf("%f", costPrice))

				if _, err := qtx.CreateOrderItem(ctx, orders_repo.CreateOrderItemParams{
					OrderID:         orderID,
					ProductID:       reqItem.ProductID,
					Quantity:        reqItem.Quantity,
					PriceAtSale:     price,
					Subtotal:        subtotal,
					NetSubtotal:     subtotal,
					CostPriceAtSale: numericCost,
				}); err != nil {
					return fmt.Errorf("failed to create order item: %w", err)
				}
			}
		}

		for productID, item := range currentMap {

			params := orders_repo.AddProductStockParams{ID: productID, Stock: item.Quantity}
			if _, err := qtx.AddProductStock(ctx, params); err != nil {
				return fmt.Errorf("failed to restore stock for removed item %s: %w", productID, err)
			}

			prod, err := qtx.GetProductByID(ctx, productID)
			if err != nil {
				return fmt.Errorf("failed to fetch product for stock history: %w", err)
			}

			if _, err := qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
				ProductID:     productID,
				ChangeAmount:  item.Quantity,
				PreviousStock: prod.Stock,
				CurrentStock:  prod.Stock + item.Quantity,
				ChangeType:    orders_repo.StockChangeTypeReturn,
				ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
				Note:          utils.StringPtr("Order Item Removed"),
				CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
			}); err != nil {
				return fmt.Errorf("failed to create stock history: %w", err)
			}

			if err := qtx.DeleteOrderItem(ctx, orders_repo.DeleteOrderItemParams{ID: item.ID, OrderID: orderID}); err != nil {
				return fmt.Errorf("failed to delete order item: %w", err)
			}
		}

		// Recalculate everything including promotions
		promotionID := utils.NullableUUIDToPointer(order.AppliedPromotionID)

		updatedItems, err := qtx.GetOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return err
		}

		calc, err := s.calculateOrderTotals(ctx, tx, qtx, updatedItems, promotionID, grossTotal)
		if err != nil {
			return err
		}

		_, err = qtx.UpdateOrderTotals(ctx, orders_repo.UpdateOrderTotalsParams{
			ID:                  orderID,
			GrossTotal:          calc.GrossTotal,
			NetTotal:            calc.NetTotal,
			DiscountAmount:      calc.DiscountAmount,
			TaxAmount:           calc.TaxAmount,
			ServiceChargeAmount: calc.ServiceChargeAmount,
			Version:             req.Version,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return common.ErrOrderConflict
			}
			return err
		}

		finalOrder, err = qtx.GetOrderWithDetails(ctx, orderID)
		return err
	})

	if txErr != nil {
		return nil, txErr
	}

	logDetails := map[string]interface{}{
		"updated_order_id":     orderID,
		"updated_order_status": finalOrder.Status,
	}

	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypeUPDATE,
		activity_repo.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderUpdated, map[string]interface{}{"order_id": orderID})
	}

	return s.buildOrderDetailResponseFromQueryResult(ctx, finalOrder)
}

func (s *OrderService) buildOrderDetailResponseFromQueryResult(ctx context.Context, orderWithDetails orders_repo.GetOrderWithDetailsRow) (*OrderDetailResponse, error) {
	var itemResponses []OrderItemResponse

	if orderWithDetails.Items != nil {

		itemsJSON, err := json.Marshal(orderWithDetails.Items)
		if err != nil {
			s.log.Error("Failed to re-marshal order items interface", "error", err)
			return nil, fmt.Errorf("could not process order items")
		}

		var tempItems []struct {
			orders_repo.OrderItem
			Options []orders_repo.OrderItemOption `json:"options"`
		}

		if err := json.Unmarshal(itemsJSON, &tempItems); err != nil {
			s.log.Error("Failed to unmarshal order items JSON", "error", err)
			return nil, fmt.Errorf("could not parse order items")
		}

		// Collect IDs
		var productIDs []uuid.UUID
		var optionIDs []uuid.UUID
		for _, tempItem := range tempItems {
			productIDs = append(productIDs, tempItem.ProductID)
			for _, opt := range tempItem.Options {
				optionIDs = append(optionIDs, opt.ProductOptionID)
			}
		}

		// Fetch Names
		productNameMap := make(map[uuid.UUID]string)
		productPrintCategoryMap := make(map[uuid.UUID]string)
		if len(productIDs) > 0 {
			products, err := s.productsRepo.GetProductsByIDs(ctx, productIDs)
			if err == nil {
				for _, p := range products {
					productNameMap[p.ID] = p.Name
					productPrintCategoryMap[p.ID] = string(p.PrintCategory)
				}
			} else {
				s.log.Warn("Failed to fetch product names for order detail", "error", err)
			}
		}

		optionNameMap := make(map[uuid.UUID]string)
		if len(optionIDs) > 0 {
			options, err := s.productsRepo.GetProductOptionsByIDs(ctx, optionIDs)
			if err == nil {
				for _, o := range options {
					optionNameMap[o.ID] = o.Name
				}
			} else {
				s.log.Warn("Failed to fetch option names for order detail", "error", err)
			}
		}

		for _, tempItem := range tempItems {
			var optionResponses []OrderItemOptionResponse
			for _, opt := range tempItem.Options {
				name := optionNameMap[opt.ProductOptionID]
				optionResponses = append(optionResponses, OrderItemOptionResponse{
					ProductOptionID: opt.ProductOptionID,
					OptionName:      name,
					PriceAtSale:     opt.PriceAtSale,
				})
			}
			pName := productNameMap[tempItem.ProductID]
			itemResponses = append(itemResponses, OrderItemResponse{
				ID:            tempItem.ID,
				ProductID:     tempItem.ProductID,
				ProductName:   pName,
				Quantity:      tempItem.Quantity,
				PriceAtSale:   tempItem.PriceAtSale,
				Subtotal:      tempItem.Subtotal,
				PrintCategory: productPrintCategoryMap[tempItem.ProductID],
				Options:       optionResponses,
			})
		}
	}

	return &OrderDetailResponse{
		ID:                      orderWithDetails.ID,
		UserID:                  utils.NullableUUIDToPointer(orderWithDetails.UserID),
		CustomerID:              utils.NullableUUIDToPointer(orderWithDetails.CustomerID),
		Type:                    orderWithDetails.Type,
		Status:                  orderWithDetails.Status,
		GrossTotal:              orderWithDetails.GrossTotal,
		DiscountAmount:          orderWithDetails.DiscountAmount,
		NetTotal:                orderWithDetails.NetTotal,
		TaxAmount:               orderWithDetails.TaxAmount,
		ServiceChargeAmount:     orderWithDetails.ServiceChargeAmount,
		PaymentMethodID:         orderWithDetails.PaymentMethodID,
		PaymentGatewayReference: orderWithDetails.PaymentGatewayReference,
		CashReceived:            orderWithDetails.CashReceived,
		ChangeDue:               orderWithDetails.ChangeDue,
		AppliedPromotionID:      utils.NullableUUIDToPointer(orderWithDetails.AppliedPromotionID),
		CreatedAt:               orderWithDetails.CreatedAt.Time,
		UpdatedAt:               orderWithDetails.UpdatedAt.Time,
		Version:                 orderWithDetails.Version,
		IsPaid:                  orderWithDetails.PaymentMethodID != nil,
		Items:                   itemResponses,
	}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID, req CancelOrderRequest) error {
	actorID, userIdOk := ctx.Value(common.UserIDKey).(uuid.UUID)

	txErr := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		qPrd := products_repo.New(tx)
		orderWithDetails, err := qtx.GetOrderWithDetails(ctx, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				s.log.Warn("Order not found for cancellation", "orderID", orderID)
				return common.ErrNotFound
			}
			s.log.Error("Failed to get order details for cancellation", "error", err)
			return err
		}

		if orderWithDetails.Status != orders_repo.OrderStatusOpen {
			s.log.Warn("Attempted to cancel an order that is not in 'open' state", "orderID", orderID, "currentStatus", orderWithDetails.Status)
			return common.ErrOrderNotCancellable
		}

		// Cancel Midtrans Transaction if exists
		if orderWithDetails.PaymentGatewayReference != nil && *orderWithDetails.PaymentGatewayReference != "" {
			s.log.Infof("Cancelling Midtrans transaction for order %s", orderID)
			_, err := s.midtransService.CancelTransaction(orderID.String())
			if err != nil {
				s.log.Errorf("Failed to cancel Midtrans transaction for order %s: %v", orderID, err)
				return fmt.Errorf("failed to cancel payment gateway transaction: %w", err)
			}
		}

		_, err = qtx.CancelOrder(ctx, orders_repo.CancelOrderParams{
			ID:                   orderID,
			CancellationReasonID: &req.CancellationReasonID,
			CancellationNotes:    &req.CancellationNotes,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				s.log.Warn("Attempted to cancel an order that is not in 'open' state", "orderID", orderID, "currentStatus", orderWithDetails.Status)
				return common.ErrOrderNotCancellable
			}
			s.log.Error("Failed to execute cancel order query", "error", err)
			return err
		}

		if orderWithDetails.Items != nil {
			itemsJSON, _ := json.Marshal(orderWithDetails.Items)
			var items []orders_repo.OrderItem
			json.Unmarshal(itemsJSON, &items)
			for _, item := range items {
				prod, err := qPrd.GetProductByID(ctx, item.ProductID)
				if err != nil {
					return err
				}

				_, stockErr := qPrd.AddProductStock(ctx, products_repo.AddProductStockParams{
					ID:       item.ProductID,
					Quantity: item.Quantity,
				})
				if stockErr != nil {
					return stockErr
				}

				qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
					ProductID:     item.ProductID,
					ChangeAmount:  item.Quantity,
					PreviousStock: prod.Stock,
					CurrentStock:  prod.Stock + item.Quantity,
					ChangeType:    orders_repo.StockChangeTypeReturn,
					ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
					Note:          utils.StringPtr("Order Cancelled"),
					CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
				})
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	logDetails := map[string]interface{}{
		"cancelled_order_id": orderID.String(),
		"reason_id":          req.CancellationReasonID,
		"notes":              req.CancellationNotes,
	}
	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypeCANCEL,
		activity_repo.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderUpdated, map[string]interface{}{"order_id": orderID})
	}

	metrics.OrdersCancelledTotal.Inc()

	return nil
}

func (s *OrderService) RefundOrder(ctx context.Context, orderID uuid.UUID, req RefundOrderRequest) (*OrderDetailResponse, error) {
	actorID, userIdOk := ctx.Value(common.UserIDKey).(uuid.UUID)

	var finalOrder orders_repo.GetOrderWithDetailsRow

	txErr := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		qPrd := products_repo.New(tx)

		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return common.ErrNotFound
			}
			return err
		}

		if order.PaymentMethodID == nil {
			return errors.New("only paid orders can be refunded")
		}

		_, err = qtx.RefundOrder(ctx, orderID)
		if err != nil {
			return err
		}

		items, err := qtx.GetOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return err
		}

		for _, item := range items {
			prod, err := qPrd.GetProductByID(ctx, item.ProductID)
			if err != nil {
				return err
			}

			_, stockErr := qPrd.AddProductStock(ctx, products_repo.AddProductStockParams{
				ID:       item.ProductID,
				Quantity: item.Quantity,
			})
			if stockErr != nil {
				return stockErr
			}

			qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
				ProductID:     item.ProductID,
				ChangeAmount:  item.Quantity,
				PreviousStock: prod.Stock,
				CurrentStock:  prod.Stock + item.Quantity,
				ChangeType:    orders_repo.StockChangeTypeReturn,
				ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
				Note:          utils.StringPtr("Order Refunded: " + req.Reason),
				CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
			})
		}

		finalOrder, err = qtx.GetOrderWithDetails(ctx, orderID)
		return err
	})

	if txErr != nil {
		s.log.Error("RefundOrder transaction failed", "error", txErr)
		return nil, txErr
	}

	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypeUPDATE,
		activity_repo.LogEntityTypeORDER,
		orderID.String(),
		map[string]interface{}{
			"action": "refund",
			"reason": req.Reason,
		},
	)

	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderUpdated, map[string]interface{}{"order_id": orderID})
	}

	metrics.OrdersRefundedTotal.Inc()

	return s.buildOrderDetailResponseFromQueryResult(ctx, finalOrder)
}

func (s *OrderService) ListOrders(ctx context.Context, req ListOrdersRequest) (*PagedOrderResponse, error) {
	req.SetDefaults()

	page := req.Page
	limit := req.Limit
	offset := (page - 1) * limit

	var nullUserID pgtype.UUID
	if req.UserID != nil {
		nullUserID.Valid = true
		nullUserID.Bytes = *req.UserID
	}

	var statusStrings []string
	if req.Statuses != nil {
		statusStrings = make([]string, len(req.Statuses))
		for i, st := range req.Statuses {
			statusStrings[i] = string(st)
		}
	}

	listParams := orders_repo.ListOrdersParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Statuses: statusStrings,
		UserID:   nullUserID,
	}
	countParams := orders_repo.CountOrdersParams{
		Statuses: statusStrings,
		UserID:   nullUserID,
	}
	var wg sync.WaitGroup
	var orders []orders_repo.ListOrdersRow
	var totalData int64
	var listErr, countErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		orders, listErr = s.ordersRepo.ListOrders(ctx, listParams)
	}()

	go func() {
		defer wg.Done()
		totalData, countErr = s.ordersRepo.CountOrders(ctx, countParams)
	}()

	wg.Wait()

	if listErr != nil {
		s.log.Error("Failed to list orders from repository", "error", listErr)
		return nil, listErr
	}
	if countErr != nil {
		s.log.Error("Failed to count orders from repository", "error", countErr)
		return nil, countErr
	}

	var ordersResponse []OrderListResponse
	if len(orders) == 0 {
		return &PagedOrderResponse{
			Orders:     []OrderListResponse{},
			Pagination: pagination.BuildPagination(page, int(totalData), limit),
		}, nil
	}

	// 1. Collect all order IDs
	var orderIDs []uuid.UUID
	for _, order := range orders {
		orderIDs = append(orderIDs, order.ID)
	}

	// 2. Batch fetch all items
	allItems, err := s.ordersRepo.GetOrderItemsByOrderIDs(ctx, orderIDs)
	if err != nil {
		s.log.Error("Failed to fetch order items in batch", "error", err)
	}

	// Group items by order ID and collect product IDs
	itemsByOrderID := make(map[uuid.UUID][]orders_repo.OrderItem)
	var productIDs []uuid.UUID
	productIDMap := make(map[uuid.UUID]bool)

	for _, item := range allItems {
		itemsByOrderID[item.OrderID] = append(itemsByOrderID[item.OrderID], item)
		if !productIDMap[item.ProductID] {
			productIDMap[item.ProductID] = true
			productIDs = append(productIDs, item.ProductID)
		}
	}

	// 3. Batch fetch all products
	productMap := make(map[uuid.UUID]string)
	productPrintCategoryMap := make(map[uuid.UUID]string)
	if len(productIDs) > 0 {
		products, err := s.productsRepo.GetProductsByIDs(ctx, productIDs)
		if err != nil {
			s.log.Error("Failed to fetch products for order list items", "error", err)
		} else {
			for _, p := range products {
				productMap[p.ID] = p.Name
				productPrintCategoryMap[p.ID] = string(p.PrintCategory)
			}
		}
	}

	// 4. Build response
	for _, order := range orders {
		netTotal := order.NetTotal
		items := itemsByOrderID[order.ID]

		var itemResponses []OrderItemResponse
		for _, item := range items {
			name := productMap[item.ProductID]

			itemResponses = append(itemResponses, OrderItemResponse{
				ID:            item.ID,
				ProductID:     item.ProductID,
				ProductName:   name,
				Quantity:      item.Quantity,
				PriceAtSale:   item.PriceAtSale,
				Subtotal:      item.Subtotal,
				PrintCategory: productPrintCategoryMap[item.ProductID],
			})
		}

		queueNumber := order.ID.String()[len(order.ID.String())-4:]

		isPaid := false
		if order.PaymentMethodID != nil {
			isPaid = true
		}

		ordersResponse = append(ordersResponse, OrderListResponse{
			ID:          order.ID,
			UserID:      utils.NullableUUIDToPointer(order.UserID),
			Type:        order.Type,
			Status:      order.Status,
			NetTotal:    netTotal,
			CreatedAt:   order.CreatedAt.Time,
			Items:       itemResponses,
			QueueNumber: "#" + queueNumber,
			IsPaid:      isPaid,
		})
	}

	pagedResponse := &PagedOrderResponse{
		Orders: ordersResponse,
		Pagination: pagination.BuildPagination(
			page,
			int(totalData),
			limit,
		),
	}

	return pagedResponse, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderDetailResponse, error) {
	var newOrderID uuid.UUID
	var finalOrder orders_repo.GetOrderWithDetailsRow

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for order creation")
	}

	txErr := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		qPrd := products_repo.New(tx)

		// 1. Prepare and Validate Items
		prepared, err := s.prepareAndValidateOrderItems(ctx, qPrd, req.Items)
		if err != nil {
			return err
		}

		// 2. Create Order Header
		var nullCustomerID pgtype.UUID
		if req.CustomerID != nil {
			nullCustomerID.Valid = true
			nullCustomerID.Bytes = *req.CustomerID
		}

		orderHeader, err := qtx.CreateOrder(ctx, orders_repo.CreateOrderParams{
			UserID:     pgtype.UUID{Bytes: actorID, Valid: ok},
			Type:       req.Type,
			CustomerID: nullCustomerID,
		})
		if err != nil {
			return fmt.Errorf("failed to create order header: %w", err)
		}
		newOrderID = orderHeader.ID

		// 3. Persist Items and Options
		createdItems, err := s.persistOrderItems(ctx, qtx, newOrderID, prepared, actorID, ok)
		if err != nil {
			return err
		}

		// 4. Calculate Totals
		calc, err := s.calculateOrderTotals(ctx, tx, qtx, createdItems, nil, prepared.GrossTotal)
		if err != nil {
			return err
		}

		_, err = qtx.UpdateOrderTotals(ctx, orders_repo.UpdateOrderTotalsParams{
			ID:                  newOrderID,
			GrossTotal:          calc.GrossTotal,
			NetTotal:            calc.NetTotal,
			DiscountAmount:      calc.DiscountAmount,
			TaxAmount:           calc.TaxAmount,
			ServiceChargeAmount: calc.ServiceChargeAmount,
			Version:             1,
		})
		if err != nil {
			return fmt.Errorf("failed to update order totals: %w", err)
		}

		finalOrder, err = qtx.GetOrderWithDetails(ctx, newOrderID)
		return err
	})

	if txErr != nil {
		s.log.Error("CreateOrder transaction failed", "error", txErr)
		return nil, txErr
	}

	s.activityService.Log(ctx, actorID, activity_repo.LogActionTypeCREATE, activity_repo.LogEntityTypeORDER, newOrderID.String(), map[string]interface{}{"action": "create_order", "order_id": newOrderID})
	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderCreated, map[string]interface{}{"order_id": newOrderID})
	}

	return s.buildOrderDetailResponseFromQueryResult(ctx, finalOrder)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*OrderDetailResponse, error) {

	orderWithDetails, err := s.ordersRepo.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Order not found by ID", "orderID", orderID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get order details from repository", "error", err)
		return nil, err
	}

	return s.buildOrderDetailResponseFromQueryResult(ctx, orderWithDetails)
}

func (s *OrderService) InitiateMidtransPayment(ctx context.Context, orderID uuid.UUID) (*MidtransPaymentResponse, error) {
	order, err := s.ordersRepo.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order.PaymentGatewayReference != nil {
		s.log.Infof("Order %s already has payment reference: %s. Returning existing.", orderID, *order.PaymentGatewayReference)

		if order.PaymentUrl != nil && *order.PaymentUrl != "" {
			var actions []PaymentAction
			if err := json.Unmarshal([]byte(*order.PaymentUrl), &actions); err == nil {
				return &MidtransPaymentResponse{
					OrderID:       order.ID.String(),
					TransactionID: *order.PaymentGatewayReference,
					GrossAmount:   fmt.Sprintf("%d.00", order.NetTotal),
					Actions:       actions,
				}, nil
			}
		}

	}

	chargeResp, err := s.midtransService.CreateQRISCharge(order.ID.String(), order.NetTotal)
	if err != nil {
		return nil, err
	}

	s.log.Infof("QRIS charge created successfully for Order ID: %s. Transaction ID: %s", order.ID.String(), chargeResp.TransactionID)

	var paymentActions []PaymentAction
	for _, act := range chargeResp.Actions {
		paymentActions = append(paymentActions, PaymentAction{
			Name:   act.Name,
			Method: act.Method,
			URL:    act.URL,
		})
	}

	actionsJSON, _ := json.Marshal(paymentActions)

	err = s.ordersRepo.UpdateOrderPaymentInfo(ctx, orders_repo.UpdateOrderPaymentInfoParams{
		ID:                      order.ID,
		PaymentMethodID:         nil,
		PaymentGatewayReference: utils.StringPtr(chargeResp.TransactionID),
	})
	if err != nil {
		return nil, err
	}

	paymentUrlStr := string(actionsJSON)
	err = s.ordersRepo.UpdateOrderPaymentUrl(ctx, orders_repo.UpdateOrderPaymentUrlParams{
		ID:           order.ID,
		PaymentUrl:   &paymentUrlStr,
		PaymentToken: nil,
	})
	if err != nil {
		s.log.Warnf("Failed to update payment url for order %s: %v", order.ID, err)
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypePROCESSPAYMENT,
		activity_repo.LogEntityTypeORDER,
		order.ID.String(),
		map[string]interface{}{
			"payment_gateway": "midtrans",
			"transaction_id":  chargeResp.TransactionID,
			"amount":          chargeResp.GrossAmount,
		},
	)

	response := &MidtransPaymentResponse{
		OrderID:       chargeResp.OrderID,
		TransactionID: chargeResp.TransactionID,
		GrossAmount:   chargeResp.GrossAmount,
		QRString:      chargeResp.QRString,
		ExpiryTime:    chargeResp.ExpiryTime,
		Actions:       paymentActions,
	}

	return response, nil
}

func (s *OrderService) HandleMidtransNotification(ctx context.Context, payload payment.MidtransNotificationPayload) error {
	s.log.Infof("Handling Midtrans notification for Order ID: %s", payload.OrderID)

	if err := s.midtransService.VerifyNotificationSignature(payload); err != nil {
		s.log.Error("Midtrans notification signature verification failed", "error", err, "orderID", payload.OrderID)
		return fmt.Errorf("signature verification failed")
	}

	orderIDFromPayload, err := uuid.Parse(payload.OrderID)
	if err != nil {
		s.log.Error("Invalid order ID in notification", "orderID", payload.OrderID)
		return common.ErrNotFound
	}

	order, err := s.ordersRepo.GetOrderWithDetails(ctx, orderIDFromPayload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Order not found for Midtrans notification", "orderID", payload.OrderID)
			return common.ErrNotFound
		}
		s.log.Error("Failed to get order for notification", "error", err)
		return err
	}

	if order.Status == orders_repo.OrderStatusPaid || order.Status == orders_repo.OrderStatusCancelled {
		s.log.Warn("Received notification for an already finalized order", "orderID", order.ID, "status", order.Status)
		return nil
	}

	var newStatus orders_repo.OrderStatus
	var paymentMethodID *int32

	switch payload.TransactionStatus {
	case "settlement", "capture":
		if order.Status == orders_repo.OrderStatusOpen {
			newStatus = orders_repo.OrderStatusInProgress
		} else {
			newStatus = order.Status
		}
		paymentMethodID = utils.Int32Ptr(int(PaymentMethodQRIS))
	case "cancel", "deny", "expire":
		newStatus = orders_repo.OrderStatusCancelled
	default:
		s.log.Infof("Ignoring Midtrans notification with status: %s", payload.TransactionStatus)
		return nil
	}

	updatedOrder, err := s.ordersRepo.UpdateOrderStatusByGatewayRef(ctx, orders_repo.UpdateOrderStatusByGatewayRefParams{
		PaymentGatewayReference: &payload.TransactionID,
		Status:                  newStatus,
		PaymentMethodID:         paymentMethodID,
	})
	if err != nil {
		s.log.Error("Failed to update order status from notification", "error", err, "orderID", order.ID)
		return err
	}

	userUUID := utils.NullableUUIDToPointer(updatedOrder.UserID)
	if userUUID != nil {
		s.activityService.Log(
			ctx,
			*userUUID,
			activity_repo.LogActionTypeUPDATE,
			activity_repo.LogEntityTypeORDER,
			updatedOrder.ID.String(),
			map[string]interface{}{
				"status_from":     order.Status,
				"status_to":       newStatus,
				"payment_gateway": "midtrans",
				"gateway_status":  payload.TransactionStatus,
			},
		)
	}

	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderUpdated, map[string]interface{}{"order_id": updatedOrder.ID})
	}

	return nil
}

func (s *OrderService) CalculateOrder(ctx context.Context, req CalculateOrderRequest) (*CalculateOrderResponse, error) {
	var resp CalculateOrderResponse

	err := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		qPrd := products_repo.New(tx)

		prepared, err := s.prepareAndValidateOrderItems(ctx, qPrd, req.Items)
		if err != nil {
			return err
		}

		var items []orders_repo.OrderItem
		for _, pi := range prepared.Items {
			items = append(items, pi.OrderItem)
		}

		calc, err := s.calculateOrderTotals(ctx, tx, qtx, items, req.PromotionID, prepared.GrossTotal)
		if err != nil {
			return err
		}

		resp.GrossTotal = calc.GrossTotal
		resp.DiscountAmount = calc.DiscountAmount
		resp.TaxAmount = calc.TaxAmount
		resp.ServiceChargeAmount = calc.ServiceChargeAmount
		resp.NetTotal = calc.NetTotal

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (s *OrderService) CheckoutOrder(ctx context.Context, req CheckoutOrderRequest) (*OrderDetailResponse, error) {
	var finalOrder orders_repo.GetOrderWithDetailsRow
	var newOrderID uuid.UUID

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for checkout creation")
	}

	txErr := s.store.ExecTx(ctx, func(tx pgx.Tx) error {
		qtx := orders_repo.New(tx)
		qPrd := products_repo.New(tx)

		// 1. Prepare and Validate
		prepared, err := s.prepareAndValidateOrderItems(ctx, qPrd, req.Items)
		if err != nil {
			return err
		}

		// 2. Create Order Header
		var nullCustomerID pgtype.UUID
		if req.CustomerID != nil {
			nullCustomerID.Valid = true
			nullCustomerID.Bytes = *req.CustomerID
		}

		orderHeader, err := qtx.CreateOrder(ctx, orders_repo.CreateOrderParams{
			UserID:     pgtype.UUID{Bytes: actorID, Valid: ok},
			Type:       req.Type,
			CustomerID: nullCustomerID,
		})
		if err != nil {
			return fmt.Errorf("failed to create order header: %w", err)
		}
		newOrderID = orderHeader.ID

		// 3. Persist Items
		createdItems, err := s.persistOrderItems(ctx, qtx, newOrderID, prepared, actorID, ok)
		if err != nil {
			return err
		}

		// 4. Calculate Totals and Apply Promotion
		calc, err := s.calculateOrderTotals(ctx, tx, qtx, createdItems, req.PromotionID, prepared.GrossTotal)
		if err != nil {
			return err
		}

		if req.PromotionID != nil {
			err = qtx.UpdateOrderAppliedPromotion(ctx, orders_repo.UpdateOrderAppliedPromotionParams{
				ID:                 newOrderID,
				AppliedPromotionID: pgtype.UUID{Bytes: *req.PromotionID, Valid: true},
			})
			if err != nil {
				return fmt.Errorf("failed to update order applied promotion: %w", err)
			}
		}

		_, err = qtx.UpdateOrderTotals(ctx, orders_repo.UpdateOrderTotalsParams{
			ID:                  newOrderID,
			GrossTotal:          calc.GrossTotal,
			DiscountAmount:      calc.DiscountAmount,
			NetTotal:            calc.NetTotal,
			TaxAmount:           calc.TaxAmount,
			ServiceChargeAmount: calc.ServiceChargeAmount,
			Version:             1,
		})
		if err != nil {
			return fmt.Errorf("failed to update order totals: %w", err)
		}

		// 5. Confirm Manual Payment if provided
		if req.PaymentMethodID != nil {
			cashReceived := calc.NetTotal
			if req.CashReceived != nil {
				cashReceived = *req.CashReceived
			}
			if *req.PaymentMethodID == PaymentMethodStaticQRIS {
				cashReceived = calc.NetTotal
			}
			if cashReceived < calc.NetTotal {
				return fmt.Errorf("uang kurang: tagihan %d, diterima %d", calc.NetTotal, cashReceived)
			}
			changeDue := cashReceived - calc.NetTotal

			_, err = qtx.UpdateOrderManualPayment(ctx, orders_repo.UpdateOrderManualPaymentParams{
				ID:              newOrderID,
				PaymentMethodID: utils.Int32Ptr(int(*req.PaymentMethodID)),
				CashReceived:    &cashReceived,
				ChangeDue:       &changeDue,
				Version:         2, // Updated by UpdateOrderTotals
			})
			if err != nil {
				return fmt.Errorf("failed to update manual payment: %w", err)
			}
		}

		finalOrder, err = qtx.GetOrderWithDetails(ctx, newOrderID)
		return err
	})

	if txErr != nil {
		return nil, txErr
	}

	s.activityService.Log(ctx, actorID, activity_repo.LogActionTypeCREATE, activity_repo.LogEntityTypeORDER, newOrderID.String(), map[string]interface{}{"action": "checkout_order", "order_id": newOrderID})
	if s.wsHub != nil {
		s.wsHub.BroadcastEvent(ws.EventOrderCreated, map[string]interface{}{"order_id": newOrderID})
	}

	// Record business metrics
	hasPayment := "false"
	if req.PaymentMethodID != nil {
		hasPayment = "true"
		metrics.PaymentProcessedTotal.WithLabelValues("checkout", "success").Inc()
		metrics.OrderRevenueTotal.Add(float64(finalOrder.NetTotal))
	}
	metrics.OrdersCreatedTotal.WithLabelValues(string(req.Type), hasPayment).Inc()

	return s.buildOrderDetailResponseFromQueryResult(ctx, finalOrder)
}

// --- Helpers for Refactoring ---

type orderCalculationResult struct {
	GrossTotal          int64
	DiscountAmount      int64
	TaxAmount           int64
	ServiceChargeAmount int64
	NetTotal            int64
}

type preparedOrderItems struct {
	Items      []preparedItem
	GrossTotal int64
}

type preparedItem struct {
	OrderItem orders_repo.OrderItem
	Options   []orders_repo.OrderItemOption
	Product   products_repo.Product
}

func (s *OrderService) calculateOrderTotals(
	ctx context.Context,
	tx pgx.Tx,
	qtx orders_repo.Querier,
	items []orders_repo.OrderItem,
	promotionID *uuid.UUID,
	grossTotal int64,
) (orderCalculationResult, error) {
	res := orderCalculationResult{
		GrossTotal: grossTotal,
	}

	if promotionID != nil {
		promo, err := qtx.GetPromotionByID(ctx, *promotionID)
		if err == nil && promo.IsActive {
			now := time.Now()
			if !now.Before(promo.StartDate.Time) && !now.After(promo.EndDate.Time) {
				rules, _ := qtx.GetPromotionRules(ctx, promo.ID)

				var productIDs []uuid.UUID
				for _, item := range items {
					productIDs = append(productIDs, item.ProductID)
				}
				productCategoryCache, _ := s.fetchProductCategoriesBatch(ctx, tx, productIDs)

				rulesMet := true
				for _, rule := range rules {
					switch rule.RuleType {
					case orders_repo.PromotionRuleTypeMINIMUMORDERAMOUNT:
						minAmount, _ := strconv.ParseInt(rule.RuleValue, 10, 64)
						if grossTotal < minAmount {
							rulesMet = false
						}
					case orders_repo.PromotionRuleTypeREQUIREDPRODUCT:
						reqID, _ := uuid.Parse(rule.RuleValue)
						found := false
						for _, i := range items {
							if i.ProductID == reqID {
								found = true
								break
							}
						}
						if !found {
							rulesMet = false
						}
					case orders_repo.PromotionRuleTypeREQUIREDCATEGORY:
						reqCatID, _ := strconv.Atoi(rule.RuleValue)
						found := false
						for _, i := range items {
							for _, c := range productCategoryCache[i.ProductID] {
								if c == reqCatID {
									found = true
									break
								}
							}
						}
						if !found {
							rulesMet = false
						}
					}
				}

				if rulesMet {
					var discountAmount int64
					if promo.Scope == orders_repo.PromotionScopeITEM {
						targets, _ := qtx.GetPromotionTargets(ctx, promo.ID)
						var eligibleTotal int64
						for _, item := range items {
							isEligible := false
							for _, target := range targets {
								if target.TargetType == orders_repo.PromotionTargetTypePRODUCT && target.TargetID == item.ProductID.String() {
									isEligible = true
									break
								} else if target.TargetType == orders_repo.PromotionTargetTypeCATEGORY {
									tid, _ := strconv.Atoi(target.TargetID)
									for _, c := range productCategoryCache[item.ProductID] {
										if c == tid {
											isEligible = true
											break
										}
									}
								}
							}
							if isEligible {
								eligibleTotal += item.Subtotal
							}
						}
						if eligibleTotal > 0 {
							if promo.DiscountType == orders_repo.DiscountTypePercentage {
								discountAmount = (eligibleTotal * utils.NumericToInt64(promo.DiscountValue)) / 100
							} else {
								discountAmount = utils.NumericToInt64(promo.DiscountValue)
							}
						}
					} else {
						if promo.DiscountType == orders_repo.DiscountTypePercentage {
							discountAmount = (grossTotal * utils.NumericToInt64(promo.DiscountValue)) / 100
						} else {
							discountAmount = utils.NumericToInt64(promo.DiscountValue)
						}
					}
					maxDisc := utils.NumericToInt64(promo.MaxDiscountAmount)
					if maxDisc > 0 && discountAmount > maxDisc {
						discountAmount = maxDisc
					}
					if discountAmount > grossTotal {
						discountAmount = grossTotal
					}
					res.DiscountAmount = discountAmount
				}
			}
		}
	}

	taxRate := 0.11
	serviceChargeRate := 0.0
	taxSettings, err := s.settingsService.GetTaxSettings(ctx)
	if err == nil {
		taxRate = taxSettings.TaxRate
		serviceChargeRate = taxSettings.ServiceChargeRate
	} else {
		s.log.Warn("Failed to fetch tax settings, using defaults", "error", err)
	}

	res.TaxAmount = int64(float64(res.GrossTotal-res.DiscountAmount) * taxRate)
	res.ServiceChargeAmount = int64(float64(res.GrossTotal-res.DiscountAmount) * serviceChargeRate)
	res.NetTotal = res.GrossTotal - res.DiscountAmount + res.TaxAmount + res.ServiceChargeAmount

	return res, nil
}

func (s *OrderService) prepareAndValidateOrderItems(ctx context.Context, qPrd products_repo.Querier, reqItems []CreateOrderItemRequest) (preparedOrderItems, error) {
	res := preparedOrderItems{}
	if len(reqItems) == 0 {
		return res, errors.New("order must have at least one item")
	}

	productIDs := make([]uuid.UUID, len(reqItems))
	for i, item := range reqItems {
		productIDs[i] = item.ProductID
	}

	products, err := qPrd.GetProductsForUpdate(ctx, productIDs)
	if err != nil {
		return res, fmt.Errorf("failed to lock products: %w", err)
	}

	productMap := make(map[uuid.UUID]products_repo.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	var allOptionIDs []uuid.UUID
	for _, item := range reqItems {
		for _, opt := range item.Options {
			allOptionIDs = append(allOptionIDs, opt.ProductOptionID)
		}
	}

	optionMap := make(map[uuid.UUID]products_repo.ProductOption)
	if len(allOptionIDs) > 0 {
		options, err := qPrd.GetProductOptionsByIDs(ctx, allOptionIDs)
		if err != nil {
			return res, fmt.Errorf("failed to fetch options: %w", err)
		}
		for _, opt := range options {
			optionMap[opt.ID] = opt
		}
	}

	for _, reqItem := range reqItems {
		product, exists := productMap[reqItem.ProductID]
		if !exists {
			return res, fmt.Errorf("product %s not found", reqItem.ProductID)
		}
		if product.Stock < reqItem.Quantity {
			return res, fmt.Errorf("insufficient stock for %s: available %d, requested %d", product.Name, product.Stock, reqItem.Quantity)
		}

		priceAtSale := product.Price
		var itemOptions []orders_repo.OrderItemOption
		for _, optReq := range reqItem.Options {
			option, exists := optionMap[optReq.ProductOptionID]
			if !exists {
				return res, fmt.Errorf("option %s not found", optReq.ProductOptionID)
			}
			priceAtSale += option.AdditionalPrice
			itemOptions = append(itemOptions, orders_repo.OrderItemOption{
				ProductOptionID: optReq.ProductOptionID,
				PriceAtSale:     option.AdditionalPrice,
			})
		}

		subtotal := priceAtSale * int64(reqItem.Quantity)
		res.GrossTotal += subtotal

		res.Items = append(res.Items, preparedItem{
			OrderItem: orders_repo.OrderItem{
				ProductID:   reqItem.ProductID,
				Quantity:    reqItem.Quantity,
				PriceAtSale: priceAtSale,
				Subtotal:    subtotal,
				NetSubtotal: subtotal,
			},
			Options: itemOptions,
			Product: product,
		})
	}

	return res, nil
}

func (s *OrderService) persistOrderItems(ctx context.Context, qtx orders_repo.Querier, orderID uuid.UUID, prepared preparedOrderItems, actorID uuid.UUID, ok bool) ([]orders_repo.OrderItem, error) {
	var (
		productIDs []uuid.UUID
		quantities []int32
		prices     []pgtype.Numeric
		subtotals  []pgtype.Numeric
		netSubs    []pgtype.Numeric
		costPrices []pgtype.Numeric
	)

	for _, item := range prepared.Items {
		productIDs = append(productIDs, item.OrderItem.ProductID)
		quantities = append(quantities, item.OrderItem.Quantity)
		prices = append(prices, utils.Int64ToNumeric(item.OrderItem.PriceAtSale))
		subtotals = append(subtotals, utils.Int64ToNumeric(item.OrderItem.Subtotal))
		netSubs = append(netSubs, utils.Int64ToNumeric(item.OrderItem.NetSubtotal))

		costPrice := 0.0
		if item.Product.CostPrice.Valid {
			f, _ := item.Product.CostPrice.Float64Value()
			costPrice = f.Float64
		}
		numericCost := pgtype.Numeric{}
		numericCost.Scan(fmt.Sprintf("%f", costPrice))
		costPrices = append(costPrices, numericCost)
	}

	createdItems, err := qtx.BatchCreateOrderItems(ctx, orders_repo.BatchCreateOrderItemsParams{
		OrderID:          orderID,
		ProductIds:       productIDs,
		Quantities:       quantities,
		PricesAtSale:     prices,
		Subtotals:        subtotals,
		NetSubtotals:     netSubs,
		CostPricesAtSale: costPrices,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to batch insert items: %w", err)
	}

	var batchOptionParams []orders_repo.BatchCreateOrderItemOptionsParams
	for i, pi := range prepared.Items {
		createdItem := createdItems[i]
		for _, opt := range pi.Options {
			batchOptionParams = append(batchOptionParams, orders_repo.BatchCreateOrderItemOptionsParams{
				OrderItemID:     createdItem.ID,
				ProductOptionID: opt.ProductOptionID,
				PriceAtSale:     opt.PriceAtSale,
			})
		}
	}

	if len(batchOptionParams) > 0 {
		_, err = qtx.BatchCreateOrderItemOptions(ctx, batchOptionParams)
		if err != nil {
			return nil, fmt.Errorf("failed to batch insert options: %w", err)
		}
	}

	updatedProducts, err := qtx.BatchDecreaseProductStock(ctx, orders_repo.BatchDecreaseProductStockParams{
		ProductIds: productIDs,
		Quantities: quantities,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to batch update stock: %w", err)
	}

	// Check for low stock and notify
	s.notifyLowStock(updatedProducts)

	for i, pID := range productIDs {
		qty := quantities[i]
		product := prepared.Items[i].Product
		_, err := qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
			ProductID:     pID,
			ChangeAmount:  -qty,
			PreviousStock: product.Stock,
			CurrentStock:  product.Stock - qty,
			ChangeType:    orders_repo.StockChangeTypeSale,
			ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
			Note:          utils.StringPtr("Order Created"),
			CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: ok},
		})
		if err != nil {
			return nil, err
		}
	}

	return createdItems, nil
}

func (s *OrderService) fetchProductCategoriesBatch(ctx context.Context, tx pgx.Tx, productIDs []uuid.UUID) (map[uuid.UUID][]int, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID][]int), nil
	}

	query := "SELECT product_id, category_id FROM product_categories WHERE product_id = ANY($1)"
	rows, err := tx.Query(ctx, query, productIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categoryCache := make(map[uuid.UUID][]int)
	for rows.Next() {
		var pid uuid.UUID
		var cid int
		if err := rows.Scan(&pid, &cid); err != nil {
			return nil, err
		}
		categoryCache[pid] = append(categoryCache[pid], cid)
	}

	return categoryCache, nil
}

func (s *OrderService) notifyLowStock(products []orders_repo.Product) {
	if s.wsHub == nil {
		return
	}

	for _, p := range products {
		if p.Stock <= p.LowStockThreshold {
			s.log.Warnf("Low stock alert: %s (Stock: %d, Threshold: %d)", p.Name, p.Stock, p.LowStockThreshold)
			s.wsHub.BroadcastEvent(ws.EventLowStockAlert, map[string]interface{}{
				"product_id":   p.ID,
				"product_name": p.Name,
				"stock":        p.Stock,
				"threshold":    p.LowStockThreshold,
			})
		}
	}
}
