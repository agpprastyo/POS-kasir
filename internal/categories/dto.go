package categories

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
	Limit  int32 `json:"limit" query:"limit" validate:"omitempty,gte=1,lte=100"`
	Offset int32 `json:"offset" query:"offset" validate:"omitempty,gte=0"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=3,max=100"`
}
