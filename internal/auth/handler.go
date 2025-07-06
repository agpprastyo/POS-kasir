package auth

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/middleware"

	"io"

	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AthHandler handles authentication HTTP requests.
type AthHandler struct {
	service   IAuthService
	log       *logger.Logger
	validator validator.Validator
}

func NewAuthHandler(service IAuthService, log *logger.Logger, validator validator.Validator) *AthHandler {
	return &AthHandler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

func (h *AthHandler) Loginhandler(c *fiber.Ctx) error {
	ctx := c.Context()

	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	resp, err := h.service.Login(ctx, req)
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
		SameSite: "Lax",
	})

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data: fiber.Map{
			"expired": resp.ExpiredAt,
			"profile": resp.Profile,
		},
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

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	resp, err := h.service.Register(ctx, req)
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
		h.log.Errorf("Failed to get userID from context, userId : %v", userUUID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Error:   "No user id in context",
			Message: "Failed to get user ID",
		})
	}

	//userUUID, err := uuid.Parse(userID)
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
	//		Message: "Invalid user id",
	//	})
	//}

	response, err := h.service.Profile(ctx, userUUID)
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

	if err := h.validator.Validate(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	resp, err := h.service.Register(ctx, req)
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
	ctx := c.Context()
	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.log.Errorf("Failed to get userID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "No user ID in context",
		})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		h.log.Errorf("Failed to get avatar file: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to upload avatar",
		})
	}

	fileData, err := file.Open()
	if err != nil {
		h.log.Errorf("Failed to open avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to process avatar file",
		})
	}
	defer fileData.Close()

	data, err := io.ReadAll(fileData)
	if err != nil {
		h.log.Errorf("Failed to read avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to read avatar file",
		})
	}

	response, err := h.service.UploadAvatar(ctx, userUUID, data)
	if err != nil {
		h.log.Errorf("Failed to upload avatar: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Avatar updated successfully",
		Data:    response,
	})
}
