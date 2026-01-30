package common

import (
	"errors"
)

// Common error variables
var (
	ErrNotFound                = errors.New("resource not found")
	ErrUnauthorized            = errors.New("unauthorized")
	ErrForbidden               = errors.New("forbidden")
	ErrInvalidInput            = errors.New("invalid input")
	ErrInternal                = errors.New("internal server error")
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrUserExists              = errors.New("user already exists")
	ErrUsernameExists          = errors.New("username already exists")
	ErrEmailExists             = errors.New("email already exists")
	ErrCategoryInUse           = errors.New("category is in use and cannot be deleted")
	ErrInvalidID               = errors.New("invalid ID format")
	ErrNotImplemented          = errors.New("not implemented")
	ErrCategoryNotFound        = errors.New("category not found")
	ErrCategoryExists          = errors.New("category already exists")
	ErrOrderNotCancellable     = errors.New("order cannot be cancelled, it might have been paid or already cancelled")
	ErrOrderNotModifiable      = errors.New("order cannot be modified, it might have been paid or already cancelled")
	ErrInvalidStatusTransition = errors.New("invalid status transition for the order")
	ErrPromotionNotApplicable  = errors.New("promotion is not applicable to the order items")
	ErrFileTooLarge            = errors.New("file size exceeds the maximum limit")
	ErrFileTypeNotSupported    = errors.New("file type is not supported")
	ErrImageNotSquare          = errors.New("image must be square")
	ErrImageTooSmall           = errors.New("image dimensions are too small, must be at least 300x300 pixels")
	ErrImageTooLarge           = errors.New("image dimensions are too large, must not exceed 2000x2000 pixels")
	ErrImageUploadFailed       = errors.New("image upload failed")
	ErrImageProcessingFailed   = errors.New("image processing failed")
	ErrUploadFailed            = errors.New("file upload failed")
	ErrPaymentFailed           = errors.New("payment processing failed")
	ErrAvatarNotFound          = errors.New("avatar not found")
	ErrAvatarUploadFailed      = errors.New("avatar upload failed")
	ErrAvatarProcessingFailed  = errors.New("avatar processing failed")
	ErrAvatarTooLarge          = errors.New("avatar size exceeds the maximum limit")
	ErrAvatarNotSquare         = errors.New("avatar must be square")
	ErrUploadAvatar            = errors.New("failed to upload avatar, please try again later")
	ErrAvatarLink              = errors.New("failed to generate avatar link, please try again later")
)

type ErrorResponse struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
