package user

import (
	"POS-kasir/internal/auth"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/pagination"
)

type UsersRequest struct {
	Page      *int                        `form:"page" json:"page"`
	Limit     *int                        `form:"limit" json:"limit"`
	SortBy    *repository.UserOrderColumn `form:"sortBy" json:"sortBy"`
	SortOrder *repository.SortOrder       `form:"sortOrder" json:"sortOrder"`
	Search    *string                     `form:"search" json:"search"`
	Role      *repository.UserRole        `form:"role" json:"role,omitempty"`
	IsActive  *bool                       `form:"is_active" json:"is_active,omitempty"`
}

type UsersResponse struct {
	Users      []auth.ProfileResponse `json:"users"`
	Pagination pagination.Pagination  `json:"pagination"`
}

type CreateUserRequest struct {
	Username string              `json:"username" validate:"required,min=3,max=50"`
	Email    string              `json:"email" validate:"required,email,max=100"`
	Password string              `json:"password" validate:"required,min=6,max=100"`
	Role     repository.UserRole `json:"role" validate:"required,oneof=admin manager cashier"`
	IsActive *bool               `json:"is_active,omitempty"`
}

type UpdateUserRequest struct {
	Username *string              `json:"username,omitempty" validate:"omitempty,min=3,max=50"`
	Email    *string              `json:"email,omitempty" validate:"omitempty,email,max=100"`
	Role     *repository.UserRole `json:"role,omitempty" validate:"omitempty,oneof=admin manager cashier"`
	IsActive *bool                `json:"is_active,omitempty"`
}
