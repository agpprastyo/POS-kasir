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
	Role      *repository.UserRole        `form:"filter[role]" json:"role,omitempty"`
	IsActive  *bool                       `form:"filter[isActive]" json:"is_active,omitempty"`
}

type UsersResponse struct {
	Users      []auth.ProfileResponse `json:"users"`
	Pagination pagination.Pagination  `json:"pagination"`
}
