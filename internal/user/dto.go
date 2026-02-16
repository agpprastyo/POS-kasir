package user

import (
	"POS-kasir/internal/common/pagination"
	"POS-kasir/internal/user/repository"
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ExpiredAt    time.Time       `json:"expired_at"`
	Token        string          `json:"Token"`
	RefreshToken string          `json:"refresh_token"`
	Profile      ProfileResponse `json:"profile"`
}

type RegisterRequest struct {
	Username string              `json:"username" validate:"required,min=3,max=32"`
	Email    string              `json:"email" validate:"required,email"`
	Password string              `json:"password" validate:"required,min=8,max=32"`
	Role     repository.UserRole `json:"role" validate:"required"`
}

type ProfileResponse struct {
	ID        uuid.UUID           `json:"id"`
	Username  string              `json:"username"`
	Email     string              `json:"email"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	DeletedAt *time.Time          `json:"deleted_at,omitempty"`
	Avatar    *string             `json:"avatar"`
	Role      repository.UserRole `json:"role"`
	IsActive  bool                `json:"is_active"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=8,max=32"`
	NewPassword string `json:"new_password" validate:"required,min=8,max=32"`
}

type UsersRequest struct {
	Page      *int                        `form:"page" json:"page"`
	Limit     *int                        `form:"limit" json:"limit"`
	SortBy    *repository.UserOrderColumn `form:"sortBy" json:"sortBy"`
	SortOrder *repository.SortOrder       `form:"sortOrder" json:"sortOrder"`
	Search    *string                     `form:"search" json:"search"`
	Role      *repository.UserRole        `form:"role" json:"role,omitempty"`
	IsActive  *bool                       `form:"is_active" json:"is_active,omitempty"`
	Status    *string                     `form:"status" json:"status" validate:"omitempty,oneof=active deleted all"`
}

type UsersResponse struct {
	Users      []ProfileResponse     `json:"users"`
	Pagination pagination.Pagination `json:"pagination"`
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
