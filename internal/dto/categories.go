package dto

import "time"

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
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}
