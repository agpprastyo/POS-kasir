package orders

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/pagination"
	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"sync"
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
	activityService activitylog.Service
	log             *logger.Logger
}

func NewOrderService(store repository.Store, midtransService payment.IMidtrans, activityService activitylog.Service, log *logger.Logger) IOrderService {
	return &OrderService{
		store:           store,
		midtransService: midtransService,
		activityService: activityService,
		log:             log,
	}
}

// Definisikan state machine untuk transisi status yang valid.
var allowedStatusTransitions = map[repository.OrderStatus]map[repository.OrderStatus]bool{
	repository.OrderStatusPaid: {
		repository.OrderStatusInProgress: true,
	},
	repository.OrderStatusInProgress: {
		repository.OrderStatusServed: true,
	},
}

func (s *OrderService) ApplyPromotion(ctx context.Context, orderID uuid.UUID, req dto.ApplyPromotionRequest) (*dto.OrderDetailResponse, error) {
	return nil, common.ErrNotImplemented
	//var finalOrder repository.Order
	//
	//txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
	//	// 1. Ambil data pesanan dan promosi yang relevan
	//	order, err := qtx.GetOrderForUpdate(ctx, orderID)
	//	if err != nil {
	//		return common.ErrNotFound
	//	}
	//	if order.Status != repository.OrderStatusOpen {
	//		return common.ErrOrderNotModifiable
	//	}
	//	if order.AppliedPromotionID.Valid {
	//		return fmt.Errorf("an order can only have one promotion")
	//	}
	//
	//	promo, err := qtx.GetPromotionByID(ctx, req.PromotionID)
	//	if err != nil {
	//		return common.ErrNotFound
	//	}
	//	rules, _ := qtx.GetPromotionRules(ctx, req.PromotionID)
	//
	//	// 2. Validasi Aturan Promosi (Rule Engine Sederhana)
	//	grossTotal := utils.NumericToFloat64(order.GrossTotal)
	//	for _, rule := range rules {
	//		switch rule.RuleType {
	//		case repository.PromotionRuleTypeMINIMUMORDERAMOUNT:
	//			minAmount, _ := strconv.ParseFloat(rule.RuleValue, 64)
	//			if grossTotal < minAmount {
	//				return fmt.Errorf("%w: minimum order amount is %.2f", common.ErrPromotionNotApplicable, minAmount)
	//			}
	//			// TODO: Tambahkan validasi untuk aturan lain (REQUIRED_PRODUCT, dll.)
	//		}
	//	}
	//
	//	// 3. Hitung Diskon
	//	var discountAmount float64
	//	if promo.DiscountType == repository.DiscountTypePercentage {
	//		discountValue := utils.NumericToFloat64(promo.DiscountValue)
	//		discountAmount = grossTotal * (discountValue / 100)
	//		if promo.MaxDiscountAmount.Valid {
	//			maxDiscount := utils.NumericToFloat64(promo.MaxDiscountAmount)
	//			discountAmount = math.Min(discountAmount, maxDiscount)
	//		}
	//	} else { // Fixed Amount
	//		discountAmount = utils.NumericToFloat64(promo.DiscountValue)
	//	}
	//
	//	netTotal := grossTotal - discountAmount
	//	if netTotal < 0 {
	//		netTotal = 0
	//	}
	//
	//	// 4. Update pesanan dengan diskon dan ID promosi
	//	discountNumeric, _ := utils.Float64ToNumeric(discountAmount)
	//	netTotalNumeric, _ := utils.Float64ToNumeric(netTotal)
	//
	//	finalOrder, err = qtx.UpdateOrderTotals(ctx, repository.UpdateOrderTotalsParams{
	//		ID:                orderID,
	//		GrossTotal:        order.GrossTotal, // Gross total tidak berubah
	//		DiscountAmount:    discountNumeric,
	//		NetTotal:          netTotalNumeric,
	//		AppliedPromotionID: pgtype.UUID{Bytes: promo.ID, Valid: true},
	//	})
	//	return err
	//})
	//
	//if txErr != nil {
	//	s.log.Error("Failed to apply promotion transaction", "error", txErr, "orderID", orderID)
	//	return nil, txErr
	//}
	//
	//// 5. Log aktivitas
	//actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	//s.activityService.Log(ctx, actorID, repository.LogActionTypeAPPLY_PROMOTION, repository.LogEntityTypeORDER, orderID.String(), map[string]interface{}{"promotion_id": req.PromotionID})
	//
	//// 6. Ambil data lengkap dan kembalikan
	//return s.GetOrder(ctx, orderID)
}

func (s *OrderService) UpdateOperationalStatus(ctx context.Context, orderID uuid.UUID, req dto.UpdateOrderStatusRequest) (*dto.OrderDetailResponse, error) {
	// 1. Ambil pesanan saat ini untuk validasi
	order, err := s.store.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			s.log.Warn("Order not found for status update", "orderID", orderID)
			return nil, common.ErrNotFound
		}
		s.log.Error("Failed to get order for status update", "error", err)
		return nil, err
	}

	// 2. Validasi transisi status menggunakan state machine
	currentStatus := order.Status
	newStatus := req.Status

	allowed, ok := allowedStatusTransitions[currentStatus][newStatus]
	if !ok || !allowed {
		errMsg := fmt.Sprintf("invalid status transition from '%s' to '%s'", currentStatus, newStatus)
		s.log.Warn(errMsg, "orderID", orderID)
		return nil, fmt.Errorf("%w: %s", common.ErrInvalidStatusTransition, errMsg)
	}

	// 3. Lakukan update status di database
	_, err = s.store.UpdateOrderStatus(ctx, repository.UpdateOrderStatusParams{
		ID:     orderID,
		Status: newStatus,
	})
	if err != nil {
		s.log.Error("Failed to update order status in repository", "error", err, "orderID", orderID)
		return nil, err
	}

	// 4. Log aktivitas
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

	// 5. Ambil kembali data lengkap untuk dikembalikan sebagai respons
	return s.GetOrder(ctx, orderID)
}
func (s *OrderService) CompleteManualPayment(ctx context.Context, orderID uuid.UUID, req dto.CompleteManualPaymentRequest) (*dto.OrderDetailResponse, error) {
	var updatedOrder repository.Order

	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		// 1. Ambil dan kunci pesanan untuk update
		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return common.ErrNotFound
			}
			return err
		}

		// 2. Validasi status pesanan
		if order.Status != repository.OrderStatusOpen {
			s.log.Warn("Attempted to complete payment for an order with invalid status", "orderID", orderID, "status", order.Status)
			return common.ErrOrderNotModifiable
		}

		// 3. Siapkan parameter untuk update
		netTotal := utils.NumericToFloat64(order.NetTotal)
		var changeDue float64 = 0

		// Asumsi ID 1 adalah untuk 'Cash'
		isCashPayment := req.PaymentMethodID == 1

		if isCashPayment {
			if req.CashReceived < netTotal {
				return fmt.Errorf("cash received (%.2f) is less than the net total (%.2f)", req.CashReceived, netTotal)
			}
			changeDue = req.CashReceived - netTotal
		}

		cashReceivedNumeric, _ := utils.Float64ToNumeric(req.CashReceived)
		changeDueNumeric, _ := utils.Float64ToNumeric(changeDue)

		updateParams := repository.UpdateOrderManualPaymentParams{
			ID:              orderID,
			PaymentMethodID: &req.PaymentMethodID,
			CashReceived:    cashReceivedNumeric,
			ChangeDue:       changeDueNumeric,
		}

		// 4. Lakukan update
		updatedOrder, err = qtx.UpdateOrderManualPayment(ctx, updateParams)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				// Ini terjadi jika status pesanan berubah setelah GetOrderForUpdate (sangat jarang, tapi mungkin)
				return common.ErrOrderNotModifiable
			}
			return err
		}
		return nil
	})

	if txErr != nil {
		s.log.Error("Failed to complete manual payment transaction", "error", txErr, "orderID", orderID)
		return nil, txErr
	}

	// 5. Log aktivitas
	actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	logDetails := map[string]interface{}{
		"order_id":          orderID.String(),
		"payment_method_id": req.PaymentMethodID,
		"amount":            utils.NumericToFloat64(updatedOrder.NetTotal),
	}
	s.activityService.Log(
		ctx,
		actorID,
		//repository.LogActionTypePROCESS_PAYMENT,
		repository.LogActionTypePROCESSPAYMENT,
		repository.LogEntityTypeORDER,
		orderID.String(),
		logDetails,
	)

	// 6. Ambil data lengkap dan kembalikan
	fullOrderDetails, err := s.store.GetOrderWithDetails(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return s.buildOrderDetailResponseFromQueryResult(fullOrderDetails)
}

func (s *OrderService) UpdateOrderItems(ctx context.Context, orderID uuid.UUID, reqs []dto.UpdateOrderItemRequest) (*dto.OrderDetailResponse, error) {
	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		// 1. Ambil dan kunci pesanan
		order, err := qtx.GetOrderForUpdate(ctx, orderID)
		if err != nil {
			// ... handle not found
			return err
		}
		if order.Status != repository.OrderStatusOpen {
			return common.ErrOrderNotModifiable
		}

		// 2. Lakukan iterasi pada setiap aksi
		for _, req := range reqs {
			switch req.Action {
			case "delete":
				if req.ID == nil {
					return fmt.Errorf("item id is required for delete action")
				}
				// Ambil item untuk tahu kuantitasnya
				itemToDelete, err := qtx.GetOrderItem(ctx, repository.GetOrderItemParams{ID: *req.ID, OrderID: orderID})
				if err != nil {
					return err
				}
				// Hapus item
				if err := qtx.DeleteOrderItem(ctx, repository.DeleteOrderItemParams{ID: *req.ID, OrderID: orderID}); err != nil {
					return err
				}
				// Kembalikan stok
				if _, err := qtx.AddProductStock(ctx, repository.AddProductStockParams{ID: itemToDelete.ProductID, Stock: itemToDelete.Quantity}); err != nil {
					return err
				}

			case "update":
				if req.ID == nil || req.Quantity == nil {
					return fmt.Errorf("item id and quantity are required for update action")
				}
				// Ambil item lama
				itemToUpdate, err := qtx.GetOrderItem(ctx, repository.GetOrderItemParams{ID: *req.ID, OrderID: orderID})
				if err != nil {
					return err
				}
				// Hitung selisih stok
				stockDifference := itemToUpdate.Quantity - *req.Quantity
				// Update kuantitas dan subtotal
				newPrice := utils.NumericToFloat64(itemToUpdate.PriceAtSale)
				newSubtotal, _ := utils.Float64ToNumeric(newPrice * float64(*req.Quantity))
				if _, err := qtx.UpdateOrderItemQuantity(ctx, repository.UpdateOrderItemQuantityParams{
					ID: *req.ID, OrderID: orderID, Quantity: *req.Quantity, Subtotal: newSubtotal, NetSubtotal: newSubtotal,
				}); err != nil {
					return err
				}
				// Sesuaikan stok (bisa positif atau negatif)
				if _, err := qtx.AddProductStock(ctx, repository.AddProductStockParams{ID: itemToUpdate.ProductID, Stock: stockDifference}); err != nil {
					return err
				}

			case "create":
				if req.ProductID == nil || req.Quantity == nil {
					return fmt.Errorf("product_id and quantity are required for create action")
				}
				// Logika ini mirip dengan yang ada di CreateOrder
				// Ambil produk, validasi stok, buat order_item, kurangi stok
				// ...
			}
		}

		// 3. Hitung ulang total keseluruhan pesanan
		// Ambil semua item yang ada sekarang
		finalItems, err := qtx.GetOrderItemsByOrderID(ctx, orderID)
		if err != nil {
			return err
		}
		// Lakukan loop dan jumlahkan semua `net_subtotal`
		var newGrossTotal float64 = 0
		for _, item := range finalItems {
			newGrossTotal += utils.NumericToFloat64(item.NetSubtotal)
		}
		// Update header pesanan
		grossTotalNumeric, _ := utils.Float64ToNumeric(newGrossTotal)
		if _, err := qtx.UpdateOrderTotals(ctx, repository.UpdateOrderTotalsParams{
			ID: orderID, GrossTotal: grossTotalNumeric, NetTotal: grossTotalNumeric,
		}); err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	// 4. Ambil data pesanan yang sudah final dan kembalikan
	return s.GetOrder(ctx, orderID)
}

func (s *OrderService) buildOrderDetailResponse(order repository.Order, items []repository.OrderItem, itemOptions map[uuid.UUID][]repository.OrderItemOption) *dto.OrderDetailResponse {
	var itemResponses []dto.OrderItemResponse
	for _, item := range items {
		var optionResponses []dto.OrderItemOptionResponse
		if opts, ok := itemOptions[item.ID]; ok {
			for _, opt := range opts {
				optionResponses = append(optionResponses, dto.OrderItemOptionResponse{
					ProductOptionID: opt.ProductOptionID,
					PriceAtSale:     utils.NumericToFloat64(opt.PriceAtSale),
				})
			}
		}
		itemResponses = append(itemResponses, dto.OrderItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			Quantity:    item.Quantity,
			PriceAtSale: utils.NumericToFloat64(item.PriceAtSale),
			Subtotal:    utils.NumericToFloat64(item.Subtotal),
			Options:     optionResponses,
		})
	}

	userIDPtr := utils.NullableUUIDToPointer(order.UserID)

	return &dto.OrderDetailResponse{
		ID:         order.ID,
		UserID:     userIDPtr,
		Type:       order.Type,
		Status:     order.Status,
		GrossTotal: utils.NumericToFloat64(order.GrossTotal),
		NetTotal:   utils.NumericToFloat64(order.NetTotal),
		CreatedAt:  order.CreatedAt.Time,
		Items:      itemResponses,
	}
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
					PriceAtSale:     utils.NumericToFloat64(opt.PriceAtSale),
				})
			}
			itemResponses = append(itemResponses, dto.OrderItemResponse{
				ID:          tempItem.ID,
				ProductID:   tempItem.ProductID,
				Quantity:    tempItem.Quantity,
				PriceAtSale: utils.NumericToFloat64(tempItem.PriceAtSale),
				Subtotal:    utils.NumericToFloat64(tempItem.Subtotal),
				Options:     optionResponses,
			})
		}
	}

	return &dto.OrderDetailResponse{
		ID:                      orderWithDetails.ID,
		UserID:                  utils.NullableUUIDToPointer(orderWithDetails.UserID),
		Type:                    orderWithDetails.Type,
		Status:                  orderWithDetails.Status,
		GrossTotal:              utils.NumericToFloat64(orderWithDetails.GrossTotal),
		DiscountAmount:          utils.NumericToFloat64(orderWithDetails.DiscountAmount),
		NetTotal:                utils.NumericToFloat64(orderWithDetails.NetTotal),
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
			var items []repository.OrderItem
			if err := json.Unmarshal(orderWithDetails.Items.([]byte), &items); err == nil {
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
	var orders []repository.Order
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
		netTotal := utils.NumericToFloat64(order.NetTotal)

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
	var newOrder repository.Order
	var createdItems []repository.OrderItem
	var createdItemOptions = make(map[uuid.UUID][]repository.OrderItemOption)

	actorID, ok := ctx.Value(common.UserIDKey).(uuid.UUID)
	if !ok {
		s.log.Warn("Actor user ID not found in context for order creation")

	}

	txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
		var err error

		newOrder, err = qtx.CreateOrder(ctx, repository.CreateOrderParams{
			UserID: pgtype.UUID{Bytes: actorID, Valid: ok},
			Type:   req.Type,
		})
		if err != nil {
			s.log.Error("Failed to create order header", "error", err)
			return err
		}

		var productIDs []uuid.UUID
		for _, item := range req.Items {
			productIDs = append(productIDs, item.ProductID)
		}

		products, err := qtx.GetProductsByIDs(ctx, productIDs)
		if err != nil {
			s.log.Error("Failed to get products by IDs", "error", err)
			return err
		}

		productMap := make(map[uuid.UUID]repository.Product)
		for _, p := range products {
			productMap[p.ID] = p
		}

		options, err := qtx.GetOptionsForProducts(ctx, productIDs)
		if err != nil {
			s.log.Error("Failed to get options for products", "error", err)
			return err
		}

		optionMap := make(map[uuid.UUID]repository.ProductOption)
		for _, o := range options {
			optionMap[o.ID] = o
		}

		var grossTotal float64 = 0

		for _, itemReq := range req.Items {
			product, exists := productMap[itemReq.ProductID]
			if !exists {
				return fmt.Errorf("product with ID %s not found", itemReq.ProductID)
			}
			if product.Stock < itemReq.Quantity {
				return fmt.Errorf("insufficient stock for product %s: available %d, requested %d", product.Name, product.Stock, itemReq.Quantity)
			}

			itemPrice := utils.NumericToFloat64(product.Price)

			for _, optReq := range itemReq.Options {
				option, optExists := optionMap[optReq.ProductOptionID]
				if !optExists || option.ProductID != product.ID {
					return fmt.Errorf("option with ID %s is not valid for product %s", optReq.ProductOptionID, product.Name)
				}
				itemPrice += utils.NumericToFloat64(option.AdditionalPrice)
			}

			subtotal := itemPrice * float64(itemReq.Quantity)
			grossTotal += subtotal

			priceAtSale, _ := utils.Float64ToNumeric(itemPrice)
			subtotalNumeric, _ := utils.Float64ToNumeric(subtotal)

			createdItem, err := qtx.CreateOrderItem(ctx, repository.CreateOrderItemParams{
				OrderID:     newOrder.ID,
				ProductID:   itemReq.ProductID,
				Quantity:    itemReq.Quantity,
				PriceAtSale: priceAtSale,
				Subtotal:    subtotalNumeric,
				NetSubtotal: subtotalNumeric,
			})
			if err != nil {
				return err
			}
			createdItems = append(createdItems, createdItem)

			for _, optReq := range itemReq.Options {
				option := optionMap[optReq.ProductOptionID]
				priceAtSaleOption, _ := utils.Float64ToNumeric(utils.NumericToFloat64(option.AdditionalPrice))
				createdOpt, err := qtx.CreateOrderItemOption(ctx, repository.CreateOrderItemOptionParams{
					OrderItemID:     createdItem.ID,
					ProductOptionID: optReq.ProductOptionID,
					PriceAtSale:     priceAtSaleOption,
				})
				if err != nil {
					return err
				}
				createdItemOptions[createdItem.ID] = append(createdItemOptions[createdItem.ID], createdOpt)
			}

			_, err = qtx.DecreaseProductStock(ctx, repository.DecreaseProductStockParams{
				ID:    itemReq.ProductID,
				Stock: itemReq.Quantity,
			})
			if err != nil {
				return err
			}
		}

		grossTotalNumeric, _ := utils.Float64ToNumeric(grossTotal)
		_, err = qtx.UpdateOrderTotals(ctx, repository.UpdateOrderTotalsParams{
			ID:         newOrder.ID,
			GrossTotal: grossTotalNumeric,
			NetTotal:   grossTotalNumeric,
		})
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	s.activityService.Log(
		ctx,
		actorID,
		repository.LogActionTypeCREATE,
		repository.LogEntityTypeORDER,
		newOrder.ID.String(),
		map[string]interface{}{"total": newOrder.NetTotal, "items_count": len(createdItems)},
	)

	return s.buildOrderDetailResponse(newOrder, createdItems, createdItemOptions), nil
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
	//TODO implement payment processing logic
	s.log.Info("Processing payment for order", "orderID", orderID)

	return nil, fmt.Errorf("payment processing is not implemented yet")

	//
	//order, err := s.store.GetOrderWithDetails(ctx, orderID)
	//if err != nil {
	//	if errors.Is(err, pgx.ErrNoRows) {
	//		s.log.Warn("Order not found for payment processing", "orderID", orderID)
	//		return nil, common.ErrNotFound
	//	}
	//	s.log.Error("Failed to get order for payment processing", "error", err)
	//	return nil, err
	//}
	//
	//if order.Status != repository.OrderStatusOpen {
	//	s.log.Warn("Attempted to process payment for an order with invalid status", "orderID", orderID, "status", order.Status)
	//	return nil, fmt.Errorf("payment cannot be processed for order with status: %s", order.Status)
	//}
	//
	//netTotal := utils.NumericToFloat64(order.NetTotal)
	//chargeResp, err := s.midtransService.CreateQRISCharge(order.ID.String(), int64(netTotal))
	//if err != nil {
	//	s.log.Error("Failed to create Midtrans charge", "error", err, "orderID", orderID)
	//	return nil, fmt.Errorf("failed to initiate payment gateway transaction")
	//}
	//
	//txErr := s.store.ExecTx(ctx, func(qtx *repository.Queries) error {
	//	updateParams := repository.UpdateOrderPaymentInfoParams{
	//		ID:                      order.ID,
	//		PaymentMethodID:         1,
	//		PaymentGatewayReference: "aa",
	//	}
	//	return qtx.UpdateOrderPaymentInfo(ctx, updateParams)
	//})
	//
	//if txErr != nil {
	//	s.log.Error("Failed to update order with payment info", "error", txErr, "orderID", orderID)
	//	// TODO: handle payment gateway transaction rollback if needed
	//	return nil, txErr
	//}
	//
	//actorID, _ := ctx.Value(common.UserIDKey).(uuid.UUID)
	//s.activityService.Log(
	//	ctx,
	//	actorID,
	//	repository.LogActionTypePROCESSPAYMENT,
	//	repository.LogEntityTypeORDER,
	//	order.ID.String(),
	//	map[string]interface{}{
	//		"payment_gateway": "midtrans",
	//		"transaction_id":  chargeResp.TransactionID,
	//		"amount":          netTotal,
	//	},
	//)
	//
	//response := &dto.QRISResponse{
	//	OrderID:       chargeResp.OrderID,
	//	TransactionID: chargeResp.TransactionID,
	//	GrossAmount:   chargeResp.GrossAmount,
	//	QRString:      chargeResp.QRString,
	//	ExpiryTime:    chargeResp.ExpiryTime,
	//}
	//
	//return response, nil
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
