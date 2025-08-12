package dto

import (
	"POS-kasir/pkg/pagination"
	"time"

	"github.com/google/uuid"
)

type CreateProductOptionRequest struct {
	Name            string  `json:"name" validate:"required,min=1,max=100"`
	AdditionalPrice float64 `json:"additional_price" validate:"gte=0"`
}

type CreateProductRequest struct {
	Name       string                       `json:"name" validate:"required,min=3,max=100"`
	CategoryID int32                        `json:"category_id" validate:"required,gt=0"`
	Price      float64                      `json:"price" validate:"required,gt=0"`
	Stock      int32                        `json:"stock" validate:"required,gte=0"`
	Options    []CreateProductOptionRequest `json:"options" validate:"dive"`
}

type UpdateProductRequest struct {
	Name       *string  `json:"name" validate:"omitempty,min=3,max=100"`
	CategoryID *int32   `json:"category_id" validate:"omitempty,gt=0"`
	Price      *float64 `json:"price" validate:"omitempty,gt=0"`
	Stock      *int32   `json:"stock" validate:"omitempty,gte=0"`
}

type ListProductsRequest struct {
	Page       *int    `form:"page" validate:"omitempty,gte=1"`
	Limit      *int    `form:"limit" validate:"omitempty,gte=1,lte=100"`
	CategoryID *int32  `form:"category_id" json:"category_id" query:"category_id"  validate:"omitempty,gt=0"`
	Search     *string `form:"search" validate:"omitempty,min=1,max=100"`
}

type CreateProductOptionRequestStandalone struct {
	Name            string  `json:"name" validate:"required,min=1,max=100"`
	AdditionalPrice float64 `json:"additional_price" validate:"gte=0"`
}

type UpdateProductOptionRequest struct {
	Name            *string  `json:"name" validate:"omitempty,min=1,max=100"`
	AdditionalPrice *float64 `json:"additional_price" validate:"omitempty,gte=0"`
}

type ProductOptionResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	AdditionalPrice float64   `json:"additional_price"`
	ImageURL        *string   `json:"image_url,omitempty"` // Pointer karena bisa NULL
}

type ProductResponse struct {
	ID           uuid.UUID               `json:"id"`
	Name         string                  `json:"name"`
	CategoryID   *int32                  `json:"category_id,omitempty"`
	CategoryName *string                 `json:"category_name,omitempty"`
	ImageURL     *string                 `json:"image_url,omitempty"`
	Price        float64                 `json:"price"`
	Stock        int32                   `json:"stock"`
	CreatedAt    time.Time               `json:"created_at"`
	UpdatedAt    time.Time               `json:"updated_at"`
	Options      []ProductOptionResponse `json:"options,omitempty"`
}

type ProductListResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CategoryID   *int32    `json:"category_id,omitempty"`
	CategoryName *string   `json:"category_name,omitempty"`
	ImageURL     *string   `json:"image_url,omitempty"`
	Price        float64   `json:"price"`
	Stock        int32     `json:"stock"`
}

type ListProductsResponse struct {
	Products   []ProductListResponse `json:"products"`
	Pagination pagination.Pagination `json:"pagination"`
}
