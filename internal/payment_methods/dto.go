package payment_methods

import "time"

type PaymentMethodResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
