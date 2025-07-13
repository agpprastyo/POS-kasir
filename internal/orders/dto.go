package orders

import (
	"POS-kasir/internal/repository"
	"github.com/google/uuid"
	"time"
)

// --- Request DTOs ---

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

// --- Response DTOs ---

type OrderItemOptionResponse struct {
	ProductOptionID uuid.UUID `json:"product_option_id"`
	PriceAtSale     float64   `json:"price_at_sale"`
}

type OrderItemResponse struct {
	ID          uuid.UUID                 `json:"id"`
	ProductID   uuid.UUID                 `json:"product_id"`
	Quantity    int32                     `json:"quantity"`
	PriceAtSale float64                   `json:"price_at_sale"`
	Subtotal    float64                   `json:"subtotal"`
	Options     []OrderItemOptionResponse `json:"options"`
}

type OrderResponse struct {
	ID         uuid.UUID              `json:"id"`
	UserID     uuid.UUID              `json:"user_id"`
	Type       repository.OrderType   `json:"type"`
	Status     repository.OrderStatus `json:"status"`
	GrossTotal float64                `json:"gross_total"`
	NetTotal   float64                `json:"net_total"`
	CreatedAt  time.Time              `json:"created_at"`
	Items      []OrderItemResponse    `json:"items"`
}

type QRISResponse struct {
	OrderID       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	GrossAmount   string `json:"gross_amount"`
	QRString      string `json:"qr_string"`
	ExpiryTime    string `json:"expiry_time"`
}
