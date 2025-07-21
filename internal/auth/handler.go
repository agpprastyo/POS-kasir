package auth

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/middleware"
	"POS-kasir/pkg/validator"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io"
	"mime/multipart"
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

	var req UpdatePasswordRequest
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

	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("Failed to get userID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "No user ID in context",
		})
	}

	err := h.Service.UpdatePassword(ctx, userUUID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to update password",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Password updated successfully",
	})
}

func (h *AthHandler) LoginHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	var req LoginRequest
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Error:   err.Error(),
			Message: "Failed to login",
		})
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
	var req RegisterRequest
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to register",
		})
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
		h.Log.Errorf("Failed to get userID from context, userId : %v", userUUID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Error:   "No user id in context",
			Message: "Failed to get user ID",
		})
	}

	response, err := h.Service.Profile(ctx, userUUID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to fetch profile",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data:    response,
	})

}

func (h *AthHandler) AddUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	var req RegisterRequest
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
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to register user",
			Error:   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User added successfully",
		Data:    resp,
	})
}

func (h *AthHandler) UpdateAvatarHandler(c *fiber.Ctx) error {
	fmt.Println("UpdateAvatarHandler called")
	ctx := c.Context()
	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("Failed to get userID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "No user ID in context",
		})
	}

	fmt.Println("UpdateAvatarHandler called 2")

	file, err := c.FormFile("avatar")
	if err != nil {
		h.Log.Errorf("Failed to get avatar file: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to upload avatar",
		})
	}

	fmt.Println("UpdateAvatarHandler called 3")

	fileData, err := file.Open()
	if err != nil {
		h.Log.Errorf("Failed to open avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to process avatar file",
		})
	}

	fmt.Println("UpdateAvatarHandler called 4")

	defer func(fileData multipart.File) {
		err := fileData.Close()
		if err != nil {
			h.Log.Errorf("Failed to close avatar file: %v", err)
		}
	}(fileData)

	fmt.Println("UpdateAvatarHandler called 5")

	data, err := io.ReadAll(fileData)
	if err != nil {
		h.Log.Errorf("Failed to read avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to read avatar file",
		})
	}

	fmt.Println("UpdateAvatarHandler called 6")

	response, err := h.Service.UploadAvatar(ctx, userUUID, data)
	if err != nil {
		h.Log.Errorf("Failed to upload avatar: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	fmt.Println("UpdateAvatarHandler called 7")

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Avatar updated successfully",
		Data:    response,
	})
}
