package common

import (
	"errors"
	"fmt"
)

// Common error variables
var (
	ErrNotFound           = errors.New("resource not found")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInternal           = errors.New("internal server error")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
	ErrUsernameExists     = errors.New("username already exists")
	ErrEmailExists        = errors.New("email already exists")
	ErrCategoryInUse      = errors.New("category is in use and cannot be deleted")
	ErrInvalidID          = errors.New("invalid ID format")
)

// WrapError adds context to an error
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}
