package dto

import (
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/repository"

	"time"

	"github.com/google/uuid"
)

type ApplyPromotionRequest struct {
	PromotionID uuid.UUID `json:"promotion_id" validate:"required"`
}

type CreateOrderItemOptionRequest struct {
	ProductOptionID uuid.UUID `json:"product_option_id" validate:"required"`
}

type CreateOrderItemRequest struct {
	ProductID uuid.UUID                      `json:"product_id" validate:"required"`
	Quantity  int32                          `json:"quantity" validate:"required,gt=0"`
	Options   []CreateOrderItemOptionRequest `json:"options" validate:"dive"`
}

type CreateOrderRequest struct {
	Type  repository.OrderType     `json:"type" validate:"required,oneof=dine_in takeaway"`
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type ListOrdersRequest struct {
	Page     *int                     `query:"page"`
	Limit    *int                     `query:"limit"`
	Statuses []repository.OrderStatus `query:"statuses" validate:"dive,oneof=open in_progress served paid cancelled"`
	UserID   *uuid.UUID               `query:"user_id"`
}

type CancelOrderRequest struct {
	CancellationReasonID int32  `json:"cancellation_reason_id" validate:"required,gt=0"`
	CancellationNotes    string `json:"cancellation_notes" validate:"omitempty,max=255"`
}

type UpdateOrderItemRequest struct {
	ProductID uuid.UUID                      `json:"product_id" validate:"required"`
	Quantity  int32                          `json:"quantity" validate:"required,gt=0"`
	Options   []CreateOrderItemOptionRequest `json:"options" validate:"dive"`
}

type CompleteManualPaymentRequest struct {
	PaymentMethodID int32 `json:"payment_method_id" validate:"required,gt=0"`
	CashReceived    int64 `json:"cash_received" validate:"omitempty,gte=0"`
}

type UpdateOrderStatusRequest struct {
	Status repository.OrderStatus `json:"status" validate:"required,oneof=in_progress served paid"`
}

type MidtransNotificationPayload struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

type OrderItemOptionResponse struct {
	ProductOptionID uuid.UUID `json:"product_option_id"`
	PriceAtSale     int64     `json:"price_at_sale"`
}

type OrderItemResponse struct {
	ID          uuid.UUID                 `json:"id"`
	ProductID   uuid.UUID                 `json:"product_id"`
	ProductName string                    `json:"product_name,omitempty"`
	Quantity    int32                     `json:"quantity"`
	PriceAtSale int64                     `json:"price_at_sale"`
	Subtotal    int64                     `json:"subtotal"`
	Options     []OrderItemOptionResponse `json:"options,omitempty"`
}

type OrderDetailResponse struct {
	ID                      uuid.UUID              `json:"id"`
	UserID                  *uuid.UUID             `json:"user_id,omitempty"`
	Type                    repository.OrderType   `json:"type"`
	Status                  repository.OrderStatus `json:"status"`
	GrossTotal              int64                  `json:"gross_total"`
	DiscountAmount          int64                  `json:"discount_amount"`
	NetTotal                int64                  `json:"net_total"`
	PaymentMethodID         *int32                 `json:"payment_method_id,omitempty"`
	PaymentGatewayReference *string                `json:"payment_gateway_reference,omitempty"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
	Items                   []OrderItemResponse    `json:"items"`
}

type OrderListResponse struct {
	ID          uuid.UUID              `json:"id"`
	UserID      *uuid.UUID             `json:"user_id,omitempty"`
	Type        repository.OrderType   `json:"type"`
	Status      repository.OrderStatus `json:"status"`
	NetTotal    int64                  `json:"net_total"`
	CreatedAt   time.Time              `json:"created_at"`
	Items       []OrderItemResponse    `json:"items,omitempty"`
	QueueNumber string                 `json:"queue_number,omitempty"`
	IsPaid      bool                   `json:"is_paid"`
}

type PagedOrderResponse struct {
	Orders     []OrderListResponse   `json:"orders"`
	Pagination pagination.Pagination `json:"pagination"`
}

type QRISResponse struct {
	OrderID       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	GrossAmount   string `json:"gross_amount"`
	QRString      string `json:"qr_string"`
	ExpiryTime    string `json:"expiry_time"`
}
