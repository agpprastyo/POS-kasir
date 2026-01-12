package orders

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"

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
)

type IOrderService interface {
	CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderDetailResponse, error)
	GetOrder(ctx context.Context, orderID uuid.UUID) (*dto.OrderDetailResponse, error)
	ProcessPayment(ctx context.Context, orderID uuid.UUID) (*dto.QRISResponse, error)
	HandleMidtransNotification(ctx context.Context, payload dto.MidtransNotificationPayload) error
	ListOrders(ctx context.Context, req dto.ListOrdersRequest) (*dto.PagedOrderResponse, error)
	CancelOrder(ctx context.Context, orderID uuid.UUID, req dto.CancelOrderRequest) error
	UpdateOrderItems(ctx context.Context, orderID uuid.UUID, reqs []dto.UpdateOrderItemRequest) (*dto.OrderDetailResponse, error)
	CompleteManualPayment(ctx context.Context, orderID uuid.UUID, req dto.CompleteManualPaymentRequest) (*dto.OrderDetailResponse, error)
	UpdateOperationalStatus(ctx context.Context, orderID uuid.UUID, req dto.UpdateOrderStatusRequest) (*dto.OrderDetailResponse, error)
	ApplyPromotion(ctx context.Context, orderID uuid.UUID, req dto.ApplyPromotionRequest) (*dto.OrderDetailResponse, error)
}

type OrderService struct {
	store           repository.Store
	midtransService payment.IMidtrans
	activityService activitylog.IActivityService
	log             logger.ILogger
}

func NewOrderService(store repository.Store, midtransService payment.IMidtrans, activityService activitylog.IActivityService, log logger.ILogger) IOrderService {
	return &OrderService{
		store:           store,
		midtransService: midtransService,
		activityService: activityService,
		log:             log,
	}
}

var allowedStatusTransitions = map[repository.OrderStatus]map[repository.OrderStatus]bool{
	repository.OrderStatusPaid: {
		repository.OrderStatusInProgress: true,
		repository.OrderStatusServed:     true,
	},
	repository.OrderStatusInProgress: {
		repository.OrderStatusServed: true,
	},
}

func (s *OrderService) ApplyPromotion(ctx context.Context, orderID uuid.UUID, req dto.ApplyPromotionRequest) (*dto.OrderDetailResponse, error) {
	var finalOrder repository.GetOrderWithDetailsRow

	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			return common.ErrNotFound
		}
		if order.Status != repository.OrderStatusOpen {
			return common.ErrOrderNotModifiable
		}

		promo, err := qtx.GetPromotionByID(ctx, req.PromotionID)
		if err != nil {
			return common.ErrNotFound
		}

		var discountAmount int64
		grossTotal := order.GrossTotal

		if promo.DiscountType == repository.DiscountTypePercentage {
			percentage := utils.NumericToInt64(promo.DiscountValue)
			discountAmount = (grossTotal * percentage) / 100

			maxDisc := utils.NumericToInt64(promo.MaxDiscountAmount)
			if maxDisc > 0 && discountAmount > maxDisc {
				discountAmount = maxDisc
			}
		} else {
			discountAmount = utils.NumericToInt64(promo.DiscountValue)
		}

		if discountAmount > grossTotal {
			discountAmount = grossTotal
		}

		netTotal := grossTotal - discountAmount

		_, err = qtx.UpdateOrderTotals(ctx, repository.UpdateOrderTotalsParams{
			ID:             orderID,
			GrossTotal:     grossTotal,
			DiscountAmount: discountAmount,
			NetTotal:       netTotal,
		})

		err = qtx.UpdateOrderAppliedPromotion(ctx, repository.UpdateOrderAppliedPromotionParams{
			ID:                 orderID,
			AppliedPromotionID: pgtype.UUID{Bytes: promo.ID, Valid: true},
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

	return s.buildOrderDetailResponseFromQueryResult(finalOrder)
}

func (s *OrderService) UpdateOperationalStatus(ctx context.Context, orderID uuid.UUID, req dto.UpdateOrderStatusRequest) (*dto.OrderDetailResponse, error) {

	order, err := s.store.GetOrderWithDetails(ctx, orderID)
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

	allowed, ok := allowedStatusTransitions[currentStatus][newStatus]
	if !ok || !allowed {
		errMsg := fmt.Sprintf("invalid status transition from '%s' to '%s'", currentStatus, newStatus)
		s.log.Warn(errMsg, "orderID", orderID)
		return nil, fmt.Errorf("%w: %s", common.ErrInvalidStatusTransition, errMsg)
	}

	_, err = s.store.UpdateOrderStatus(ctx, repository.UpdateOrderStatusParams{
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
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	return s.GetOrder(ctx, orderID)
}
func (s *OrderService) CompleteManualPayment(ctx context.Context, orderID uuid.UUID, req dto.CompleteManualPaymentRequest) (*dto.OrderDetailResponse, error) {
	var finalOrder repository.GetOrderWithDetailsRow

	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return common.ErrNotFound
			}
			return err
		}

		if order.Status != repository.OrderStatusOpen {
			return common.ErrOrderNotModifiable
		}

		netTotal := order.NetTotal
		cashReceived := req.CashReceived

		if req.PaymentMethodID == 3 {
			if cashReceived == 0 {
				cashReceived = netTotal
			}
		}

		if cashReceived < netTotal {
			return fmt.Errorf("uang kurang: tagihan %d, diterima %d", netTotal, cashReceived)
		}

		changeDue := cashReceived - netTotal

		_, err = qtx.UpdateOrderManualPayment(ctx, repository.UpdateOrderManualPaymentParams{
			ID:              orderID,
			PaymentMethodID: utils.Int32Ptr(int(req.PaymentMethodID)),
			CashReceived:    &cashReceived,
			ChangeDue:       &changeDue,
		})

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
		repository.LogActionTypePROCESSPAYMENT,
		repository.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	return s.buildOrderDetailResponseFromQueryResult(finalOrder)
}

func (s *OrderService) UpdateOrderItems(ctx context.Context, orderID uuid.UUID, reqs []dto.UpdateOrderItemRequest) (*dto.OrderDetailResponse, error) {
	var finalOrder repository.GetOrderWithDetailsRow

	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {

		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			return common.ErrNotFound
		}
		if order.Status != repository.OrderStatusOpen {
			return common.ErrOrderNotModifiable
		}

		existingItems, err := qtx.GetOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return err
		}

		currentMap := make(map[uuid.UUID]repository.OrderItem)
		for _, item := range existingItems {
			currentMap[item.ProductID] = item
		}

		reqMap := make(map[uuid.UUID]int32)
		for _, req := range reqs {
			reqMap[req.ProductID] = req.Quantity
		}

		var grossTotal int64 = 0

		for _, req := range reqs {
			product, err := qtx.GetProductByID(ctx, req.ProductID)
			if err != nil {
				return err
			}

			price := product.Price

			subtotal := price * int64(req.Quantity)
			grossTotal += subtotal

			if existingItem, exists := currentMap[req.ProductID]; exists {

				qtyDiff := req.Quantity - existingItem.Quantity

				if qtyDiff > 0 {

					if product.Stock < qtyDiff {
						return fmt.Errorf("insufficient stock for update %s", product.Name)
					}
					qtx.DecreaseProductStock(ctx, repository.DecreaseProductStockParams{ID: req.ProductID, Stock: qtyDiff})
				} else if qtyDiff < 0 {

					restoreQty := -qtyDiff
					qtx.AddProductStock(ctx, repository.AddProductStockParams{ID: req.ProductID, Stock: restoreQty})
				}

				qtx.UpdateOrderItemQuantity(ctx, repository.UpdateOrderItemQuantityParams{
					ID:          existingItem.ID,
					OrderID:     orderID,
					Quantity:    req.Quantity,
					Subtotal:    subtotal,
					NetSubtotal: subtotal,
				})

				delete(currentMap, req.ProductID)

			} else {
				if product.Stock < req.Quantity {
					return fmt.Errorf("insufficient stock for new item %s", product.Name)
				}

				qtx.DecreaseProductStock(ctx, repository.DecreaseProductStockParams{ID: req.ProductID, Stock: req.Quantity})

				qtx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
					OrderID:     orderID,
					ProductID:   req.ProductID,
					Quantity:    req.Quantity,
					PriceAtSale: price,
					Subtotal:    subtotal,
					NetSubtotal: subtotal,
				})
			}
		}

		for productID, item := range currentMap {

			qtx.AddProductStock(ctx, repository.AddProductStockParams{ID: productID, Stock: item.Quantity})

			qtx.DeleteOrderItem(ctx, repository.DeleteOrderItemParams{ID: item.ID, OrderID: orderID})
		}

		_, err = qtx.UpdateOrderTotals(ctx, repository.UpdateOrderTotalsParams{
			ID:             orderID,
			GrossTotal:     grossTotal,
			NetTotal:       grossTotal,
			DiscountAmount: 0,
		})

		finalOrder, err = qtx.GetOrderWithDetails(ctx, orderID)
		return err
	})

	if txErr != nil {
		return nil, txErr
	}

	return s.buildOrderDetailResponseFromQueryResult(finalOrder)
}

func (s *OrderService) buildOrderDetailResponseFromQueryResult(orderWithDetails repository.GetOrderWithDetailsRow) (*dto.OrderDetailResponse, error) {
	var itemResponses []dto.OrderItemResponse

	if orderWithDetails.Items != nil {

		itemsJSON, err := json.Marshal(orderWithDetails.Items)
		if err != nil {
			s.log.Error("Failed to re-marshal order items interface", "error", err)
			return nil, fmt.Errorf("could not process order items")
		}

		var tempItems []struct {
			repository.OrderItem
			Options []repository.OrderItemOption `json:"options"`
		}

		if err := json.Unmarshal(itemsJSON, &tempItems); err != nil {
			s.log.Error("Failed to unmarshal order items JSON", "error", err)
			return nil, fmt.Errorf("could not parse order items")
		}

		for _, tempItem := range tempItems {
			var optionResponses []dto.OrderItemOptionResponse
			for _, opt := range tempItem.Options {
				optionResponses = append(optionResponses, dto.OrderItemOptionResponse{
					ProductOptionID: opt.ProductOptionID,
					PriceAtSale:     opt.PriceAtSale,
				})
			}
			itemResponses = append(itemResponses, dto.OrderItemResponse{
				ID:          tempItem.ID,
				ProductID:   tempItem.ProductID,
				Quantity:    tempItem.Quantity,
				PriceAtSale: tempItem.PriceAtSale,
				Subtotal:    tempItem.Subtotal,
				Options:     optionResponses,
			})
		}
	}

	return &dto.OrderDetailResponse{
		ID:                      orderWithDetails.ID,
		UserID:                  utils.NullableUUIDToPointer(orderWithDetails.UserID),
		Type:                    orderWithDetails.Type,
		Status:                  orderWithDetails.Status,
		GrossTotal:              orderWithDetails.GrossTotal,
		DiscountAmount:          orderWithDetails.DiscountAmount,
		NetTotal:                orderWithDetails.NetTotal,
		PaymentMethodID:         orderWithDetails.PaymentMethodID,
		PaymentGatewayReference: orderWithDetails.PaymentGatewayReference,
		CreatedAt:               orderWithDetails.CreatedAt.Time,
		UpdatedAt:               orderWithDetails.UpdatedAt.Time,
		Items:                   itemResponses,
	}, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID uuid.UUID, req dto.CancelOrderRequest) error {
	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		orderWithDetails, err := qtx.GetOrderWithDetails(ctx, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				s.log.Warn("Order not found for cancellation", "orderID", orderID)
				return common.ErrNotFound
			}
			s.log.Error("Failed to get order details for cancellation", "error", err)
			return err
		}

		_, err = qtx.CancelOrder(ctx, repository.CancelOrderParams{
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
			switch v := orderWithDetails.Items.(type) {
			case []byte:
				var items []repository.OrderItem
				if err := json.Unmarshal(v, &items); err != nil {
					for _, item := range items {
						_, stockErr := qtx.AddProductStock(ctx, repository.AddProductStockParams{
							ID:    item.ProductID,
							Stock: item.Quantity,
						})
						if stockErr != nil {
							s.log.Error("Failed to restore stock for product", "error", stockErr, "productID", item.ProductID)
							return stockErr
						}
					}
				}
			case []interface{}:
				for _, item := range v {
					if itemMap, ok := item.(map[string]interface{}); ok {
						productID, ok := itemMap["product_id"].(uuid.UUID)
						if !ok {
							s.log.Error("Invalid product ID in order items", "item", item)
							continue
						}
						quantity, ok := itemMap["quantity"].(float64)
						if !ok {
							s.log.Error("Invalid quantity in order items", "item", item)
							continue
						}
						_, stockErr := qtx.AddProductStock(ctx, repository.AddProductStockParams{
							ID:    productID,
							Stock: int32(quantity),
						})
						if stockErr != nil {
							s.log.Error("Failed to restore stock for product", "error", stockErr, "productID", productID)
							return stockErr
						}
					} else {
						s.log.Error("Unexpected item type in order items", "item", item)
					}
				}
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"cancelled_order_id": orderID.String(),
		"reason_id":          req.CancellationReasonID,
		"notes":              req.CancellationNotes,
	}
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeCANCEL,
		repository.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	s.log.Info("Order cancelled successfully", "orderID", orderID)
	return nil
}

func (s *OrderService) ListOrders(ctx context.Context, req dto.ListOrdersRequest) (*dto.PagedOrderResponse, error) {

	page := 1
	if req.Page != nil {
		page = *req.Page
	}
	limit := 10
	if req.Limit != nil {
		limit = *req.Limit
	}
	offset := (page - 1) * limit

	var nullStatus repository.NullOrderStatus
	if req.Status != nil {
		nullStatus.Valid = true
		nullStatus.OrderStatus = *req.Status
	}

	var nullUserID pgtype.UUID
	if req.UserID != nil {
		nullUserID.Valid = true
		nullUserID.Bytes = *req.UserID
	}

	listParams := repository.ListOrdersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Status: nullStatus,
		UserID: nullUserID,
	}
	countParams := repository.CountOrdersParams{
		Status: nullStatus,
		UserID: nullUserID,
	}
	var wg sync.WaitGroup
	var orders []repository.ListOrdersRow
	var totalData int64
	var listErr, countErr error

	wg.Add(2)

	go func() {
		defer wg.Done()
		orders, listErr = s.store.ListOrders(ctx, listParams)
	}()

	go func() {
		defer wg.Done()
		totalData, countErr = s.store.CountOrders(ctx, countParams)
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

	var ordersResponse []dto.OrderListResponse
	for _, order := range orders {
		netTotal := order.NetTotal

		ordersResponse = append(ordersResponse, dto.OrderListResponse{
			ID:        order.ID,
			UserID:    utils.NullableUUIDToPointer(order.UserID),
			Type:      order.Type,
			Status:    order.Status,
			NetTotal:  netTotal,
			CreatedAt: order.CreatedAt.Time,
		})
	}

	pagedResponse := &dto.PagedOrderResponse{
		Orders: ordersResponse,
		Pagination: pagination.BuildPagination(
			page,
			int(totalData),
			limit,
		),
	}

	return pagedResponse, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, req dto.CreateOrderRequest) (*dto.OrderDetailResponse, error) {
	var newOrderID uuid.UUID
	var finalOrder repository.GetOrderWithDetailsRow

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for order creation")
	}

	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {

		orderHeader, err := qtx.CreateOrder(ctx, repository.CreateOrderParams{
			UserID: pgtype.UUID{Bytes: actorID, Valid: ok},
			Type:   req.Type,
		})
		if err != nil {
			return fmt.Errorf("failed to create order header: %w", err)
		}
		newOrderID = orderHeader.ID

		productIDs := make([]uuid.UUID, len(req.Items))
		for i, item := range req.Items {
			productIDs[i] = item.ProductID
		}

		products, err := qtx.GetProductsForUpdate(ctx, productIDs)
		if err != nil {
			return fmt.Errorf("failed to lock products: %w", err)
		}

		productMap := make(map[uuid.UUID]repository.Product)
		for _, p := range products {
			productMap[p.ID] = p
		}

		var allOptionIDs []uuid.UUID
		for _, item := range req.Items {
			for _, opt := range item.Options {
				allOptionIDs = append(allOptionIDs, opt.ProductOptionID)
			}
		}

		optionMap := make(map[uuid.UUID]repository.ProductOption)
		if len(allOptionIDs) > 0 {
			options, err := qtx.GetProductOptionsByIDs(ctx, allOptionIDs)
			if err != nil {
				return fmt.Errorf("failed to fetch options: %w", err)
			}
			for _, opt := range options {
				optionMap[opt.ID] = opt
			}
		}

		var (
			itemOrderIDs   []uuid.UUID
			itemProductIDs []uuid.UUID
			itemQuantities []int32
			itemPrices     []pgtype.Numeric
			itemSubtotals  []pgtype.Numeric
			itemNetSubs    []pgtype.Numeric
			stockUpdateIDs []uuid.UUID
			stockUpdateQty []int32
		)

		var grossTotal int64 = 0

		for _, itemReq := range req.Items {
			product, exists := productMap[itemReq.ProductID]
			if !exists {
				return fmt.Errorf("product %s not found", itemReq.ProductID)
			}

			if product.Stock < itemReq.Quantity {
				return fmt.Errorf("insufficient stock for %s: available %d, requested %d", product.Name, product.Stock, itemReq.Quantity)
			}

			priceAtSale := product.Price

			for _, optReq := range itemReq.Options {
				option, exists := optionMap[optReq.ProductOptionID]
				if !exists {
					return fmt.Errorf("option %s not found (or belongs to different product)", optReq.ProductOptionID)
				}
				priceAtSale += option.AdditionalPrice
			}

			subtotal := priceAtSale * int64(itemReq.Quantity)
			grossTotal += subtotal

			itemOrderIDs = append(itemOrderIDs, newOrderID)
			itemProductIDs = append(itemProductIDs, itemReq.ProductID)
			itemQuantities = append(itemQuantities, itemReq.Quantity)
			itemPrices = append(itemPrices, utils.Int64ToNumeric(priceAtSale))
			itemSubtotals = append(itemSubtotals, utils.Int64ToNumeric(subtotal))
			itemNetSubs = append(itemNetSubs, utils.Int64ToNumeric(subtotal))

			stockUpdateIDs = append(stockUpdateIDs, itemReq.ProductID)
			stockUpdateQty = append(stockUpdateQty, itemReq.Quantity)
		}

		createdItems, err := qtx.BatchCreateOrderItems(ctx, repository.BatchCreateOrderItemsParams{
			OrderID:      newOrderID,
			ProductIds:   itemProductIDs,
			Quantities:   itemQuantities,
			PricesAtSale: itemPrices,
			Subtotals:    itemSubtotals,
			NetSubtotals: itemNetSubs,
		})
		if err != nil {
			return fmt.Errorf("failed to batch insert items: %w", err)
		}

		var batchOptionParams []repository.BatchCreateOrderItemOptionsParams

		if len(createdItems) != len(req.Items) {
			return fmt.Errorf("mismatch between requested items (%d) and created items (%d)", len(req.Items), len(createdItems))
		}

		for i, reqItem := range req.Items {
			createdItem := createdItems[i]

			for _, optReq := range reqItem.Options {
				option, _ := optionMap[optReq.ProductOptionID]

				batchOptionParams = append(batchOptionParams, repository.BatchCreateOrderItemOptionsParams{
					OrderItemID:     createdItem.ID,
					ProductOptionID: optReq.ProductOptionID,
					PriceAtSale:     option.AdditionalPrice,
				})
			}
		}

		if len(batchOptionParams) > 0 {
			rowCount, err := qtx.BatchCreateOrderItemOptions(ctx, batchOptionParams)
			if err != nil {
				return fmt.Errorf("failed to batch insert options: %w", err)
			}
			if rowCount != int64(len(batchOptionParams)) {
				s.log.Warn("Batch insert options row count mismatch", "expected", len(batchOptionParams), "actual", rowCount)
			}
		}

		err = qtx.BatchDecreaseProductStock(ctx, repository.BatchDecreaseProductStockParams{
			ProductIds: stockUpdateIDs,
			Quantities: stockUpdateQty,
		})
		if err != nil {
			return fmt.Errorf("failed to batch update stock: %w", err)
		}

		_, err = qtx.UpdateOrderTotals(ctx, repository.UpdateOrderTotalsParams{
			ID:             newOrderID,
			GrossTotal:     grossTotal,
			NetTotal:       grossTotal,
			DiscountAmount: 0,
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

	go func() {
		s.activityService.Log(context.Background(), actorID, repository.LogActionTypeCREATE, repository.LogEntityTypeORDER, newOrderID.String(), nil)
	}()

	return s.buildOrderDetailResponseFromQueryResult(finalOrder)
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (*dto.OrderDetailResponse, error) {

	orderWithDetails, err := s.store.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Order not found by ID", "orderID", orderID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get order details from repository", "error", err)
		return nil, err
	}

	return s.buildOrderDetailResponseFromQueryResult(orderWithDetails)
}

func (s *OrderService) ProcessPayment(ctx context.Context, orderID uuid.UUID) (*dto.QRISResponse, error) {
	order, err := s.store.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		return nil, err
	}

	chargeResp, err := s.midtransService.CreateQRISCharge(order.ID.String(), order.NetTotal)
	if err != nil {
		return nil, err
	}

	err = s.store.UpdateOrderPaymentInfo(ctx, repository.UpdateOrderPaymentInfoParams{
		ID:                      order.ID,
		PaymentMethodID:         utils.Int32Ptr(2),
		PaymentGatewayReference: utils.StringPtr(chargeResp.TransactionID),
	})
	if err != nil {
		return nil, err
	}

	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypePROCESSPAYMENT,
		repository.LogEntityTypeORDER,
		order.ID.String(),
		map[string]interface{}{
			"payment_gateway": "midtrans",
			"transaction_id":  chargeResp.TransactionID,
			"amount":          chargeResp.GrossAmount,
		},
	)

	response := &dto.QRISResponse{
		OrderID:       chargeResp.OrderID,
		TransactionID: chargeResp.TransactionID,
		GrossAmount:   chargeResp.GrossAmount,
		QRString:      chargeResp.QRString,
		ExpiryTime:    chargeResp.ExpiryTime,
	}

	return response, nil
}

func (s *OrderService) HandleMidtransNotification(ctx context.Context, payload dto.MidtransNotificationPayload) error {
	s.log.Infof("Handling Midtrans notification for Order ID: %s", payload.OrderID)

	if err := s.midtransService.VerifyNotificationSignature(payload); err != nil {
		s.log.Error("Midtrans notification signature verification failed", "error", err, "orderID", payload.OrderID)
		return fmt.Errorf("signature verification failed")
	}

	order, err := s.store.GetOrderByGatewayRef(ctx, &payload.TransactionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Order not found for Midtrans notification", "transactionID", payload.TransactionID)
			return common.ErrNotFound
		}
		s.log.Error("Failed to get order by gateway reference", "error", err)
		return err
	}

	if order.Status == repository.OrderStatusPaid || order.Status == repository.OrderStatusCancelled {
		s.log.Warn("Received notification for an already finalized order", "orderID", order.ID, "status", order.Status)
		return nil
	}

	var newStatus repository.OrderStatus
	switch payload.TransactionStatus {
	case "settlement", "capture":

		newStatus = repository.OrderStatusPaid
	case "cancel", "deny", "expire":

		newStatus = repository.OrderStatusCancelled
	default:

		s.log.Infof("Ignoring Midtrans notification with status: %s", payload.TransactionStatus)
		return nil
	}

	updatedOrder, err := s.store.UpdateOrderStatusByGatewayRef(ctx, repository.UpdateOrderStatusByGatewayRefParams{
		PaymentGatewayReference: &payload.TransactionID,
		Status:                  newStatus,
	})
	if err != nil {
		s.log.Error("Failed to update order status from notification", "error", err, "orderID", order.ID)
		return err
	}

	userUUID := utils.NullableUUIDToPointer(updatedOrder.UserID)
	s.activityService.Log(
		ctx,
		*userUUID,
		repository.LogActionTypeUPDATE,
		repository.LogEntityTypeORDER,
		updatedOrder.ID.String(),
		map[string]interface{}{
			"status_from":     order.Status,
			"status_to":       newStatus,
			"payment_gateway": "midtrans",
			"gateway_status":  payload.TransactionStatus,
		},
	)

	s.log.Info("Successfully updated order status from notification", "orderID", updatedOrder.ID, "newStatus", newStatus)
	return nil
}
