package auth

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/middleware"
	"POS-kasir/pkg/validator"
	"errors"
	"io"
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type IAuthHandler interface {
	UpdatePasswordHandler(c *fiber.Ctx) error
	LoginHandler(c *fiber.Ctx) error
	LogoutHandler(c *fiber.Ctx) error
	RegisterHandler(c *fiber.Ctx) error
	ProfileHandler(c *fiber.Ctx) error
	AddUserHandler(c *fiber.Ctx) error
	UpdateAvatarHandler(c *fiber.Ctx) error
}

// AthHandler handles authentication HTTP requests.
type AthHandler struct {
	Service   IAuthService
	Log       logger.ILogger
	Validator validator.Validator
}

func NewAuthHandler(service IAuthService, log logger.ILogger, validator validator.Validator) IAuthHandler {
	return &AthHandler{
		Service:   service,
		Log:       log,
		Validator: validator,
	}
}

// UpdatePasswordHandler handles password update requests.
func (h *AthHandler) UpdatePasswordHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("UpdatePasswordHandler | Failed to get userID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "No user ID in context",
		})
	}

	var req dto.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
			Error:   err.Error(),
		})
	}

	if err := h.Validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	err := h.Service.UpdatePassword(ctx, userUUID, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "Invalid current password",
			})
		case errors.Is(err, common.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		default:
			h.Log.Errorf("Auth Handler | Failed to update password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to update password",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Password updated successfully",
	})
}

func (h *AthHandler) LoginHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
			Error:   err.Error(),
		})
	}

	if err := h.Validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	resp, err := h.Service.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, common.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "Invalid username or password",
			})
		default:
			h.Log.Errorf("LoginHandler | Failed to login: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to login",
			})
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    resp.Token,
		Path:     "/",
		Expires:  resp.ExpiredAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data: fiber.Map{
			"expired": resp.ExpiredAt,
			"profile": resp.Profile,
		},
	})
}

func (h *AthHandler) LogoutHandler(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:  "access_token",
		Value: "",
	})
	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Successfully logged out",
	})
}

func (h *AthHandler) RegisterHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
		})
	}

	if err := h.Validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	resp, err := h.Service.Register(ctx, req)
	if err != nil {
		switch err {
		case common.ErrUserExists:
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "User already exists",
			})
		case common.ErrUsernameExists:
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Username already exists",
			})
		case common.ErrEmailExists:
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Email already exists",
			})
		default:
			h.Log.Errorf("RegisterHandler | Failed to register user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to register user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data:    resp,
	})
}

func (h *AthHandler) ProfileHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("ProfileHandler | Failed to get userID from context, userId : %v", userUUID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Error:   "No user id in context",
			Message: "Failed to get user ID",
		})
	}

	response, err := h.Service.Profile(ctx, userUUID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, common.ErrAvatarNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "Avatar not found",
			})
		default:
			h.Log.Errorf("ProfileHandler | Failed to get user profile: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to get user profile",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data:    response,
	})

}

func (h *AthHandler) AddUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req dto.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
		})
	}

	currentRoleVal := c.Locals("role")
	currentRole, ok := currentRoleVal.(repository.UserRole)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
			Message: "Invalid current user role",
		})
	}

	if middleware.RoleLevel[req.Role] >= middleware.RoleLevel[currentRole] {
		return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
			Message: "Cannot assign a role equal to or higher than your own",
		})
	}

	if err := h.Validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	resp, err := h.Service.Register(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrUserExists):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "User already exists",
			})
		case errors.Is(err, common.ErrUsernameExists):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Username already exists",
			})
		case errors.Is(err, common.ErrEmailExists):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Email already exists",
			})
		default:
			h.Log.Errorf("AddUserHandler | Failed to register user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to register user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User added successfully",
		Data:    resp,
	})
}

func (h *AthHandler) UpdateAvatarHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("UpdateAvatarHandler | Failed to get userID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "No user ID in context",
		})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		h.Log.Errorf("UpdateAvatarHandler | Failed to get avatar file: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to upload avatar",
		})
	}

	fileData, err := file.Open()
	if err != nil {
		h.Log.Errorf("UpdateAvatarHandler | Failed to open avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to process avatar file",
		})
	}

	defer func(fileData multipart.File) {
		err := fileData.Close()
		if err != nil {
			h.Log.Errorf("UpdateAvatarHandler | Failed to close avatar file: %v", err)
		}
	}(fileData)

	data, err := io.ReadAll(fileData)
	if err != nil {
		h.Log.Errorf("UpdateAvatarHandler | Failed to read avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to read avatar file",
		})
	}

	response, err := h.Service.UploadAvatar(ctx, userUUID, data)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrFileTooLarge):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "File size exceeds the maximum limit",
			})
		case errors.Is(err, common.ErrFileTypeNotSupported):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "File type is not supported",
			})
		case errors.Is(err, common.ErrImageNotSquare):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Image must be square",
			})
		case errors.Is(err, common.ErrImageTooSmall):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Image dimensions are too small, must be at least 300x300 pixels",
			})
		case errors.Is(err, common.ErrImageProcessingFailed):
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Image processing failed",
			})
		case errors.Is(err, common.ErrUploadFailed):
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "File upload failed",
			})

		default:
			h.Log.Errorf("UpdateAvatarHandler | Failed to upload avatar: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to upload avatar",
			})
		}

	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Avatar updated successfully",
		Data:    response,
	})
}
