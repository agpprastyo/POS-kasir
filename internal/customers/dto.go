package customers

import (
	"POS-kasir/internal/common/pagination"
	"time"

	"github.com/google/uuid"
)

type CreateCustomerRequest struct {
	Name    string  `json:"name" validate:"required,max=100"`
	Phone   *string `json:"phone" validate:"omitempty,max=20"`
	Email   *string `json:"email" validate:"omitempty,email,max=255"`
	Address *string `json:"address" validate:"omitempty"`
}

type UpdateCustomerRequest struct {
	Name    string  `json:"name" validate:"required,max=100"`
	Phone   *string `json:"phone" validate:"omitempty,max=20"`
	Email   *string `json:"email" validate:"omitempty,email,max=255"`
	Address *string `json:"address" validate:"omitempty"`
}

type ListCustomersRequest struct {
	pagination.PaginationRequest
}

type CustomerResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	Phone     *string    `json:"phone,omitempty"`
	Email     *string    `json:"email,omitempty"`
	Address   *string    `json:"address,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type PagedCustomerResponse struct {
	Customers  []CustomerResponse    `json:"customers"`
	Pagination pagination.Pagination `json:"pagination"`
}
