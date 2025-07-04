package auth

import (
	"POS-kasir/internal/repository"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	ExpiredAt time.Time       `json:"expired_at"`
	Token     string          `json:"token"`
	Profile   ProfileResponse `json:"profile"`
}

type RegisterRequest struct {
	Username string              `json:"username" validate:"required,min=3,max=32"`
	Email    string              `json:"email" validate:"required,email"`
	Password string              `json:"password" validate:"required,min=8,max=32"`
	Role     repository.UserRole `json:"role" validate:"required"`
}

type ProfileResponse struct {
	Username  string              `json:"username"`
	Email     string              `json:"email"`
	CreatedAt time.Time           `json:"created_at"`
	UpdatedAt time.Time           `json:"updated_at"`
	Avatar    *string             `json:"avatar"`
	Role      repository.UserRole `json:"role"`
}
