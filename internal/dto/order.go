// File: internal/orders/order.go
package dto

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/pagination"
	"github.com/google/uuid"
	"time"
)

// --- DTO untuk Request (Input dari Klien) ---

// CreateOrderItemOptionRequest merepresentasikan satu varian yang dipilih untuk sebuah item.
type CreateOrderItemOptionRequest struct {
	ProductOptionID uuid.UUID `json:"product_option_id" validate:"required"`
}

// CreateOrderItemRequest merepresentasikan satu item produk di dalam keranjang pesanan.
type CreateOrderItemRequest struct {
	ProductID uuid.UUID                      `json:"product_id" validate:"required"`
	Quantity  int32                          `json:"quantity" validate:"required,gt=0"`
	Options   []CreateOrderItemOptionRequest `json:"options" validate:"dive"`
}

// CreateOrderRequest adalah DTO utama untuk membuat pesanan baru.
// Digunakan pada endpoint `POST /orders`.
type CreateOrderRequest struct {
	Type  repository.OrderType     `json:"type" validate:"required,oneof=dine_in takeaway"`
	Items []CreateOrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

// ListOrdersRequest menampung parameter query untuk filtering dan pagination.
// Digunakan pada endpoint `GET /orders`.
type ListOrdersRequest struct {
	Page   *int                    `form:"page"`
	Limit  *int                    `form:"limit"`
	Status *repository.OrderStatus `form:"status" validate:"omitempty,oneof=open in_progress served paid cancelled"`
	UserID *uuid.UUID              `form:"user_id"`
}

// CancelOrderRequest adalah DTO untuk membatalkan pesanan.
// Digunakan pada endpoint `POST /orders/{id}/cancel`.
type CancelOrderRequest struct {
	CancellationReasonID int32  `json:"cancellation_reason_id" validate:"required,gt=0"`
	CancellationNotes    string `json:"cancellation_notes" validate:"omitempty,max=255"`
}

// MidtransNotificationPayload adalah struktur untuk menerima notifikasi webhook dari Midtrans.
// Digunakan pada endpoint `POST /payments/midtrans-notification`.
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

// --- DTO untuk Response (Output ke Klien) ---

// OrderItemOptionResponse adalah struktur data untuk satu varian di dalam item pesanan.
type OrderItemOptionResponse struct {
	ProductOptionID uuid.UUID `json:"product_option_id"`
	PriceAtSale     float64   `json:"price_at_sale"`
}

// OrderItemResponse adalah struktur data untuk satu item di dalam detail pesanan.
type OrderItemResponse struct {
	ID          uuid.UUID                 `json:"id"`
	ProductID   uuid.UUID                 `json:"product_id"`
	Quantity    int32                     `json:"quantity"`
	PriceAtSale float64                   `json:"price_at_sale"`
	Subtotal    float64                   `json:"subtotal"`
	Options     []OrderItemOptionResponse `json:"options,omitempty"`
}

// OrderDetailResponse adalah struktur data lengkap untuk satu pesanan.
// Digunakan sebagai respons untuk `GET /orders/{id}` dan setelah membuat pesanan.
type OrderDetailResponse struct {
	ID                      uuid.UUID              `json:"id"`
	UserID                  *uuid.UUID             `json:"user_id,omitempty"`
	Type                    repository.OrderType   `json:"type"`
	Status                  repository.OrderStatus `json:"status"`
	GrossTotal              float64                `json:"gross_total"`
	DiscountAmount          float64                `json:"discount_amount"`
	NetTotal                float64                `json:"net_total"`
	PaymentMethodID         *int32                 `json:"payment_method_id,omitempty"`
	PaymentGatewayReference *string                `json:"payment_gateway_reference,omitempty"`
	CreatedAt               time.Time              `json:"created_at"`
	UpdatedAt               time.Time              `json:"updated_at"`
	Items                   []OrderItemResponse    `json:"items"`
}

// OrderListResponse adalah struktur data yang lebih ringkas untuk daftar pesanan.
type OrderListResponse struct {
	ID        uuid.UUID              `json:"id"`
	UserID    *uuid.UUID             `json:"user_id,omitempty"`
	Type      repository.OrderType   `json:"type"`
	Status    repository.OrderStatus `json:"status"`
	NetTotal  float64                `json:"net_total"`
	CreatedAt time.Time              `json:"created_at"`
}

// PagedOrderResponse adalah struktur respons untuk daftar pesanan yang terpaginasi.
type PagedOrderResponse struct {
	Orders     []OrderListResponse   `json:"orders"`
	Pagination pagination.Pagination `json:"pagination"`
}

// QRISResponse adalah struktur respons saat meminta pembayaran QRIS.
type QRISResponse struct {
	OrderID       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	GrossAmount   string `json:"gross_amount"`
	QRString      string `json:"qr_string"`
	ExpiryTime    string `json:"expiry_time"`
}
