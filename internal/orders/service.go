package orders

import (
	"POS-kasir/internal/activitylog"
	activity_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/common/store"
	orders_repo "POS-kasir/internal/orders/repository"
	products_repo "POS-kasir/internal/products/repository"
	"POS-kasir/pkg/logger"
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
}

type OrderService struct {
	store           store.Store
	ordersRepo      orders_repo.Querier
	productsRepo    products_repo.Querier
	midtransService payment.IMidtrans
	activityService activitylog.IActivityService
	log             logger.ILogger
}

func NewOrderService(store store.Store, ordersRepo orders_repo.Querier, productsRepo products_repo.Querier, midtransService payment.IMidtrans, activityService activitylog.IActivityService, log logger.ILogger) IOrderService {
	return &OrderService{
		store:           store,
		ordersRepo:      ordersRepo,
		productsRepo:    productsRepo,
		midtransService: midtransService,
		activityService: activityService,
		log:             log,
	}
}

var allowedStatusTransitions = map[orders_repo.OrderStatus]map[orders_repo.OrderStatus]bool{
	orders_repo.OrderStatusOpen: {
		orders_repo.OrderStatusInProgress: true,
		orders_repo.OrderStatusServed:     true,
		orders_repo.OrderStatusPaid:       true,
		orders_repo.OrderStatusCancelled:  true,
		orders_repo.OrderStatusOpen:       true,
	},
	orders_repo.OrderStatusInProgress: {
		orders_repo.OrderStatusInProgress: true,
		orders_repo.OrderStatusServed:     true,
		orders_repo.OrderStatusPaid:       true,
		orders_repo.OrderStatusCancelled:  true,
		orders_repo.OrderStatusOpen:       true,
	},
	orders_repo.OrderStatusServed: {
		orders_repo.OrderStatusServed:     true,
		orders_repo.OrderStatusPaid:       true,
		orders_repo.OrderStatusInProgress: true,
		orders_repo.OrderStatusCancelled:  true,
		orders_repo.OrderStatusOpen:       true,
	},
	orders_repo.OrderStatusPaid: {
		orders_repo.OrderStatusPaid:       true,
		orders_repo.OrderStatusServed:     true,
		orders_repo.OrderStatusInProgress: true,
		orders_repo.OrderStatusOpen:       true,
		orders_repo.OrderStatusCancelled:  true,
	},
	orders_repo.OrderStatusCancelled: {
		orders_repo.OrderStatusCancelled: true,
		orders_repo.OrderStatusOpen:      true,
		orders_repo.OrderStatusInProgress: true,
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

		promo, err := qtx.GetPromotionByID(ctx, req.PromotionID)
		if err != nil {
			return common.ErrNotFound
		}

		now := time.Now()
		if !promo.IsActive {
			return fmt.Errorf("%w: promotion is not active", common.ErrPromotionNotApplicable)
		}
		if now.Before(promo.StartDate.Time) {
			return fmt.Errorf("%w: promotion has not started yet", common.ErrPromotionNotApplicable)
		}
		if now.After(promo.EndDate.Time) {
			return fmt.Errorf("%w: promotion has expired", common.ErrPromotionNotApplicable)
		}

		rules, err := qtx.GetPromotionRules(ctx, promo.ID)
		if err != nil {
			return fmt.Errorf("failed to get promotion rules: %w", err)
		}

		productCategoryCache := make(map[uuid.UUID][]int)
		getProductCategories := func(id uuid.UUID) []int {
			if cats, ok := productCategoryCache[id]; ok {
				return cats
			}
			var cats []int
			rows, err := tx.Query(ctx, "SELECT category_id FROM product_categories WHERE product_id = $1", id)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var cid int
					if err := rows.Scan(&cid); err == nil {
						cats = append(cats, cid)
					}
				}
			}
			productCategoryCache[id] = cats
			return cats
		}

		for _, rule := range rules {
			switch rule.RuleType {
			case orders_repo.PromotionRuleTypeMINIMUMORDERAMOUNT:
				minAmount, err := strconv.ParseInt(rule.RuleValue, 10, 64)
				if err != nil {
					s.log.Warnf("Invalid rule value for MINIMUM_ORDER_AMOUNT: %s", rule.RuleValue)
					continue
				}
				if order.GrossTotal < minAmount {
					return fmt.Errorf("%w: minimum order amount not met (min: %d)", common.ErrPromotionNotApplicable, minAmount)
				}

			case orders_repo.PromotionRuleTypeREQUIREDPRODUCT:
				requiredProductID, err := uuid.Parse(rule.RuleValue)
				if err != nil {
					s.log.Warnf("Invalid rule value for REQUIRED_PRODUCT: %s", rule.RuleValue)
					continue
				}
				found := false
				for _, item := range orderItems {
					if item.ProductID == requiredProductID {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("%w: required product not found in order", common.ErrPromotionNotApplicable)
				}

			case orders_repo.PromotionRuleTypeREQUIREDCATEGORY:
				requiredCategoryID, err := strconv.Atoi(rule.RuleValue)
				if err != nil {
					s.log.Warnf("Invalid rule value for REQUIRED_CATEGORY: %s", rule.RuleValue)
					continue
				}
				found := false
				for _, item := range orderItems {
					cats := getProductCategories(item.ProductID)
					for _, cid := range cats {
						if cid == requiredCategoryID {
							found = true
							break
						}
					}
					if found {
						break
					}
				}
				if !found {
					return fmt.Errorf("%w: required category item not found in order", common.ErrPromotionNotApplicable)
				}
			}
		}

		var discountAmount int64
		grossTotal := order.GrossTotal

		if promo.Scope == orders_repo.PromotionScopeITEM {
			targets, err := qtx.GetPromotionTargets(ctx, promo.ID)
			if err != nil {
				return fmt.Errorf("failed to get promotion targets: %w", err)
			}

			var eligibleTotal int64
			for _, item := range orderItems {
				isEligible := false
				for _, target := range targets {
					if target.TargetType == orders_repo.PromotionTargetTypePRODUCT {
						if target.TargetID == item.ProductID.String() {
							isEligible = true
							break
						}
					} else if target.TargetType == orders_repo.PromotionTargetTypeCATEGORY {
						targetCatID, _ := strconv.Atoi(target.TargetID)
						cats := getProductCategories(item.ProductID)
						for _, cid := range cats {
							if cid == targetCatID {
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
					percentage := utils.NumericToInt64(promo.DiscountValue)
					discountAmount = (eligibleTotal * percentage) / 100
				} else {
					discountAmount = utils.NumericToInt64(promo.DiscountValue)
				}
			}

		} else {

			if promo.DiscountType == orders_repo.DiscountTypePercentage {
				percentage := utils.NumericToInt64(promo.DiscountValue)
				discountAmount = (grossTotal * percentage) / 100
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

		taxAmount := order.TaxAmount
		serviceChargeAmount := order.ServiceChargeAmount
		netTotal := grossTotal - discountAmount + taxAmount + serviceChargeAmount

		_, err = qtx.UpdateOrderTotals(ctx, orders_repo.UpdateOrderTotalsParams{
			ID:                  orderID,
			GrossTotal:          grossTotal,
			DiscountAmount:      discountAmount,
			NetTotal:            netTotal,
			TaxAmount:           taxAmount,
			ServiceChargeAmount: serviceChargeAmount,
			Version:             order.Version,
		})
		if err != nil {
			return err
		}

		err = qtx.UpdateOrderAppliedPromotion(ctx, orders_repo.UpdateOrderAppliedPromotionParams{
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

		if req.PaymentMethodID == 3 {
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

		reqMap := make(map[uuid.UUID]int32)
		for _, item := range req.Items {
			reqMap[item.ProductID] = item.Quantity
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
					qPrd.DecreaseProductStock(ctx, products_repo.DecreaseProductStockParams{ID: reqItem.ProductID, Quantity: qtyDiff})

					// Log Stock Decrease
					qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
						ProductID:     reqItem.ProductID,
						ChangeAmount:  -qtyDiff,
						PreviousStock: product.Stock,
						CurrentStock:  product.Stock - qtyDiff,
						ChangeType:    orders_repo.StockChangeTypeSale,
						ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
						Note:          utils.StringPtr("Order Item Qty Increase"),
						CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
					})
				} else if qtyDiff < 0 {

					restoreQty := -qtyDiff
					qtx.AddProductStock(ctx, orders_repo.AddProductStockParams{ID: reqItem.ProductID, Stock: restoreQty})

					// Log Stock Increase
					qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
						ProductID:     reqItem.ProductID,
						ChangeAmount:  restoreQty,
						PreviousStock: product.Stock,
						CurrentStock:  product.Stock + restoreQty,
						ChangeType:    orders_repo.StockChangeTypeReturn, // or Correction? Return seems appropriate for reducing order qty
						ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
						Note:          utils.StringPtr("Order Item Qty Decrease"),
						CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
					})
				}

				qtx.UpdateOrderItemQuantity(ctx, orders_repo.UpdateOrderItemQuantityParams{
					ID:          existingItem.ID,
					OrderID:     orderID,
					Quantity:    reqItem.Quantity,
					Subtotal:    subtotal,
					NetSubtotal: subtotal,
				})

				delete(currentMap, reqItem.ProductID)

			} else {
				if product.Stock < reqItem.Quantity {
					return fmt.Errorf("insufficient stock for new item %s", product.Name)
				}

				qPrd.DecreaseProductStock(ctx, products_repo.DecreaseProductStockParams{ID: reqItem.ProductID, Quantity: reqItem.Quantity})

				// Log Stock Decrease (New Item)
				qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
					ProductID:     reqItem.ProductID,
					ChangeAmount:  -reqItem.Quantity,
					PreviousStock: product.Stock,
					CurrentStock:  product.Stock - reqItem.Quantity,
					ChangeType:    orders_repo.StockChangeTypeSale,
					ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
					Note:          utils.StringPtr("Order Item Added"),
					CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
				})

				costPrice := 0.0
				if product.CostPrice.Valid {
					f, _ := product.CostPrice.Float64Value()
					costPrice = f.Float64
				}
				numericCost := pgtype.Numeric{}
				numericCost.Scan(fmt.Sprintf("%f", costPrice))

				qtx.CreateOrderItem(ctx, orders_repo.CreateOrderItemParams{
					OrderID:         orderID,
					ProductID:       reqItem.ProductID,
					Quantity:        reqItem.Quantity,
					PriceAtSale:     price,
					Subtotal:        subtotal,
					NetSubtotal:     subtotal,
					CostPriceAtSale: numericCost,
				})
			}
		}

		for productID, item := range currentMap {

			params := orders_repo.AddProductStockParams{ID: productID, Stock: item.Quantity}
			qtx.AddProductStock(ctx, params)

			prod, err := qtx.GetProductByID(ctx, productID)
			if err == nil {
				qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
					ProductID:     productID,
					ChangeAmount:  item.Quantity,
					PreviousStock: prod.Stock,
					CurrentStock:  prod.Stock + item.Quantity,
					ChangeType:    orders_repo.StockChangeTypeReturn,
					ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
					Note:          utils.StringPtr("Order Item Removed"),
					CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
				})
			} else {
				s.log.Warn("Failed to fetch product for stock history logging on item delete", "productID", productID)
			}

			qtx.DeleteOrderItem(ctx, orders_repo.DeleteOrderItemParams{ID: item.ID, OrderID: orderID})
		}

		taxAmount := order.TaxAmount
		serviceChargeAmount := order.ServiceChargeAmount
		netTotal := grossTotal - order.DiscountAmount + taxAmount + serviceChargeAmount

		_, err = qtx.UpdateOrderTotals(ctx, orders_repo.UpdateOrderTotalsParams{
			ID:                  orderID,
			GrossTotal:          grossTotal,
			NetTotal:            netTotal,
			DiscountAmount:      order.DiscountAmount,
			TaxAmount:           taxAmount,
			ServiceChargeAmount: serviceChargeAmount,
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
		activity_repo.LogActionTypeAPPLYPROMOTION,
		activity_repo.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

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
		if len(productIDs) > 0 {
			products, err := s.productsRepo.GetProductsByIDs(ctx, productIDs)
			if err == nil {
				for _, p := range products {
					productNameMap[p.ID] = p.Name
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
				ID:          tempItem.ID,
				ProductID:   tempItem.ProductID,
				ProductName: pName,
				Quantity:    tempItem.Quantity,
				PriceAtSale: tempItem.PriceAtSale,
				Subtotal:    tempItem.Subtotal,
				Options:     optionResponses,
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
				// We log the error but we might want to proceed or block.
				// Given safety first: if we cannot cancel the payment, we should probably not cancel the order locally
				// to avoid a state where user pays for a cancelled order.
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
			switch v := orderWithDetails.Items.(type) {
			case []byte:
				var items []orders_repo.OrderItem
				if err := json.Unmarshal(v, &items); err != nil {
					return err
				}
				for _, item := range items {
					// Fetch product to get stock
					prod, err := qPrd.GetProductByID(ctx, item.ProductID)
					if err != nil {
						s.log.Error("Failed to fetch product for stock return", "error", err)
						return err
					}

					_, stockErr := qPrd.AddProductStock(ctx, products_repo.AddProductStockParams{
						ID:       item.ProductID,
						Quantity: item.Quantity,
					})
					if stockErr != nil {
						s.log.Error("Failed to restore stock for product", "error", stockErr, "productID", item.ProductID)
						return stockErr
					}

					// Log Stock History
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

						// Fetch product
						prod, err := qPrd.GetProductByID(ctx, productID)
						if err != nil {
							return err
						}

						_, stockErr := qPrd.AddProductStock(ctx, products_repo.AddProductStockParams{
							ID:       productID,
							Quantity: int32(quantity),
						})
						if stockErr != nil {
							s.log.Error("Failed to restore stock for product", "error", stockErr, "productID", productID)
							return stockErr
						}

						// Log Stock History
						qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
							ProductID:     productID,
							ChangeAmount:  int32(quantity),
							PreviousStock: prod.Stock,
							CurrentStock:  prod.Stock + int32(quantity),
							ChangeType:    orders_repo.StockChangeTypeReturn,
							ReferenceID:   pgtype.UUID{Bytes: orderID, Valid: true},
							Note:          utils.StringPtr("Order Cancelled"),
							CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: userIdOk},
						})
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

	s.log.Info("Order cancelled successfully", "orderID", orderID)
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
		for i, s := range req.Statuses {
			statusStrings[i] = string(s)
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
	for _, order := range orders {
		netTotal := order.NetTotal

		items, err := s.ordersRepo.GetOrderItemsByOrderID(ctx, order.ID)
		if err != nil {
			s.log.Error("Failed to fetch items for order list", "orderID", order.ID, "error", err)
			continue
		}

		var productIDs []uuid.UUID
		for _, item := range items {
			productIDs = append(productIDs, item.ProductID)
		}

		var productMap map[uuid.UUID]string
		if len(productIDs) > 0 {
			products, err := s.productsRepo.GetProductsByIDs(ctx, productIDs)
			if err != nil {
				s.log.Error("Failed to fetch products for order list items", "error", err)

			} else {
				productMap = make(map[uuid.UUID]string)
				for _, p := range products {
					productMap[p.ID] = p.Name
				}
			}
		}

		var itemResponses []OrderItemResponse
		for _, item := range items {
			name := ""
			if productMap != nil {
				if n, ok := productMap[item.ProductID]; ok {
					name = n
				}
			}

			itemResponses = append(itemResponses, OrderItemResponse{
				ID:          item.ID,
				ProductID:   item.ProductID,
				ProductName: name,
				Quantity:    item.Quantity,
				PriceAtSale: item.PriceAtSale,
				Subtotal:    item.Subtotal,
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

		productIDs := make([]uuid.UUID, len(req.Items))
		for i, item := range req.Items {
			productIDs[i] = item.ProductID
		}

		products, err := qPrd.GetProductsForUpdate(ctx, productIDs)
		if err != nil {
			return fmt.Errorf("failed to lock products: %w", err)
		}

		productMap := make(map[uuid.UUID]products_repo.Product)
		for _, p := range products {
			productMap[p.ID] = p
		}

		var allOptionIDs []uuid.UUID
		for _, item := range req.Items {
			for _, opt := range item.Options {
				allOptionIDs = append(allOptionIDs, opt.ProductOptionID)
			}
		}

		optionMap := make(map[uuid.UUID]products_repo.ProductOption)
		if len(allOptionIDs) > 0 {
			options, err := qPrd.GetProductOptionsByIDs(ctx, allOptionIDs)
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
			itemCostPrices []pgtype.Numeric
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

			costPrice := 0.0
			if product.CostPrice.Valid {
				f, _ := product.CostPrice.Float64Value()
				costPrice = f.Float64
			}
			numericCost := pgtype.Numeric{}
			numericCost.Scan(fmt.Sprintf("%f", costPrice))
			itemCostPrices = append(itemCostPrices, numericCost)

			stockUpdateIDs = append(stockUpdateIDs, itemReq.ProductID)
			stockUpdateQty = append(stockUpdateQty, itemReq.Quantity)
		}

		createdItems, err := qtx.BatchCreateOrderItems(ctx, orders_repo.BatchCreateOrderItemsParams{
			OrderID:          newOrderID,
			ProductIds:       itemProductIDs,
			Quantities:       itemQuantities,
			PricesAtSale:     itemPrices,
			Subtotals:        itemSubtotals,
			NetSubtotals:     itemNetSubs,
			CostPricesAtSale: itemCostPrices,
		})
		if err != nil {
			return fmt.Errorf("failed to batch insert items: %w", err)
		}

		var batchOptionParams []orders_repo.BatchCreateOrderItemOptionsParams

		if len(createdItems) != len(req.Items) {
			return fmt.Errorf("mismatch between requested items (%d) and created items (%d)", len(req.Items), len(createdItems))
		}

		for i, reqItem := range req.Items {
			createdItem := createdItems[i]

			for _, optReq := range reqItem.Options {
				option, _ := optionMap[optReq.ProductOptionID]

				batchOptionParams = append(batchOptionParams, orders_repo.BatchCreateOrderItemOptionsParams{
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

		err = qtx.BatchDecreaseProductStock(ctx, orders_repo.BatchDecreaseProductStockParams{
			ProductIds: stockUpdateIDs,
			Quantities: stockUpdateQty,
		})
		if err != nil {
			return fmt.Errorf("failed to batch update stock: %w", err)
		}

		for i, pID := range stockUpdateIDs {
			qty := stockUpdateQty[i]
			product := productMap[pID]
			previousStock := product.Stock
			currentStock := previousStock - qty

			_, err := qtx.CreateStockHistory(ctx, orders_repo.CreateStockHistoryParams{
				ProductID:     pID,
				ChangeAmount:  -qty,
				PreviousStock: previousStock,
				CurrentStock:  currentStock,
				ChangeType:    orders_repo.StockChangeTypeSale,
				ReferenceID:   pgtype.UUID{Bytes: newOrderID, Valid: true},
				Note:          utils.StringPtr("Order Created"),
				CreatedBy:     pgtype.UUID{Bytes: actorID, Valid: ok},
			})
			if err != nil {
				// We log error but maybe valid to strictly fail transaction?
				// For data integrity, strictly failing is better.
				return fmt.Errorf("failed to log stock history for %s: %w", pID, err)
			}
		}

		taxAmount := int64(float64(grossTotal) * 0.11) // Default 11% tax
		serviceChargeAmount := int64(0)
		netTotal := grossTotal + taxAmount + serviceChargeAmount

		_, err = qtx.UpdateOrderTotals(ctx, orders_repo.UpdateOrderTotalsParams{
			ID:                  newOrderID,
			GrossTotal:          grossTotal,
			NetTotal:            netTotal,
			DiscountAmount:      0,
			TaxAmount:           taxAmount,
			ServiceChargeAmount: serviceChargeAmount,
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

	logDetails := map[string]interface{}{
		"created_order_id":     newOrderID,
		"created_order_status": finalOrder.Status,
	}

	s.activityService.Log(
		ctx,
		actorID,
		activity_repo.LogActionTypeCREATE,
		activity_repo.LogEntityTypeORDER,
		newOrderID.String(),
		logDetails,
	)

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
					GrossAmount:   fmt.Sprintf("%d.00", order.NetTotal), // Approximation
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
	s.log.Infof("QRIS charge response: %+v", chargeResp)

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
		PaymentMethodID:         nil, // Updated on settlement
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
		// If order is still 'open', move to 'in_progress' instead of 'paid'
		// This follows the new flow where 'paid' is the final status after 'served'
		if order.Status == orders_repo.OrderStatusOpen {
			newStatus = orders_repo.OrderStatusInProgress
		} else {
			newStatus = order.Status
		}
		paymentMethodID = utils.Int32Ptr(2) // QRIS/Midtrans
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

	s.log.Info("Successfully updated order status from notification", "orderID", updatedOrder.ID, "newStatus", newStatus)
	return nil
}
