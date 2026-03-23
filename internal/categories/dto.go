package categories

import (
	"POS-kasir/internal/common/pagination"
	"time"
)

type CategoryResponse struct {
	ID        int32     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CategoryWithCountResponse struct {
	ID           int32     `json:"id"`
	Name         string    `json:"name"`
	ProductCount int32     `json:"product_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ListCategoryRequest struct {
	pagination.PaginationRequest
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}
