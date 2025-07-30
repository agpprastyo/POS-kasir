package user

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

func NewUsrHandler(service IUsrService, log logger.ILogger, validator validator.Validator) IUsrHandler {
	return &UsrHandler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

type IUsrHandler interface {
	GetAllUsersHandler(c *fiber.Ctx) error
	CreateUserHandler(c *fiber.Ctx) error
	GetUserByIDHandler(c *fiber.Ctx) error
	UpdateUserHandler(c *fiber.Ctx) error
	ToggleUserStatusHandler(c *fiber.Ctx) error
	DeleteUserHandler(c *fiber.Ctx) error
}

type UsrHandler struct {
	service   IUsrService
	log       logger.ILogger
	validator validator.Validator
}

func (h *UsrHandler) DeleteUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		h.log.Errorf("DeleteUserHandler | User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("DeleteUserHandler | Invalid user ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	if err := h.service.DeleteUser(ctx, idParsed); err != nil {
		h.log.Errorf("DeleteUserHandler | Failed to delete user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User deleted successfully",
	})
}

// GetAllUsersHandler handles the request to get all users
func (h *UsrHandler) GetAllUsersHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(UsersRequest)
	if err := c.QueryParser(req); err != nil {
		h.log.Errorf("GetAllUsersHandler | Failed to parse query parameters: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
	}

	isActiveStr := c.Query("is_active")
	if isActiveStr != "" {
		parsedBool, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			req.IsActive = &parsedBool
		}
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Errorf("GetAllUsersHandler | Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	response, err := h.service.GetAllUsers(ctx, *req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			h.log.Warnf("GetAllUsersHandler | No users found")
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "No users found",
			})
		default:
			h.log.Errorf("GetAllUsersHandler | Failed to get users: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to get users",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    response,
	})
}

// CreateUserHandler handles the request to create a new user
func (h *UsrHandler) CreateUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(CreateUserRequest)
	if err := c.BodyParser(req); err != nil {
		h.log.Errorf("CreateUserHandler | Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
			Error:   err.Error(),
		})
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Errorf("CreateUserHandler | Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	user, err := h.service.CreateUser(ctx, *req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrUserExists):
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "User already exists",
			})
		case errors.Is(err, common.ErrUsernameExists):
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Username already exists",
			})
		case errors.Is(err, common.ErrEmailExists):
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Email already exists",
			})
		default:
			h.log.Errorf("CreateUserHandler | Failed to create user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to create user",
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "User created successfully",
		Data:    user,
	})
}

// GetUserByIDHandler handles the request to get a user by ID
func (h *UsrHandler) GetUserByIDHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		h.log.Errorf("GetUserByIDHandler | Failed to get user by ID: %v", id)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "User ID is required",
		})
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("GetUserByIDHandler | Invalid user ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid user ID format",
		})
	}

	user, err := h.service.GetUserByID(ctx, idParsed)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):

			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		default:

			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to get user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// UpdateUserHandler handles the request to update a user by ID
func (h *UsrHandler) UpdateUserHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	role := c.Locals("role")
	if id == "" {
		h.log.Errorf("UpdateUserHandler | User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "User ID is required",
		})
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("UpdateUserHandler | Invalid user ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid user ID format",
		})
	}

	req := new(UpdateUserRequest)
	if err := c.BodyParser(req); err != nil {
		h.log.Errorf("UpdateUserHandler | Failed to parse request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid request body",
		})
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Errorf("UpdateUserHandler | Validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Validation failed",
			Error:   err.Error(),
		})
	}

	if req.Role != nil {
		if role != repository.UserRoleAdmin {
			h.log.Warnf("UpdateUserHandler | Unauthorized attempt to change user role by non-admin user")
			return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
				Message: "You are not authorized to change user roles",
			})
		}
	}

	user, err := h.service.UpdateUser(ctx, idParsed, *req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			h.log.Warnf("UpdateUserHandler | User not found: %v", id)
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, common.ErrUsernameExists):
			h.log.Warnf("UpdateUserHandler | Username already exists: %v", *req.Username)
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Username already exists",
			})
		case errors.Is(err, common.ErrEmailExists):
			h.log.Warnf("UpdateUserHandler | Email already exists: %v", *req.Email)
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Email already exists",
			})
		default:
			h.log.Errorf("UpdateUserHandler | Failed to update user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to update user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User updated successfully",
		Data:    user,
	})
}

// ToggleUserStatusHandler handles the request to toggle a user's active status
func (h *UsrHandler) ToggleUserStatusHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	id := c.Params("id")
	if id == "" {
		h.log.Errorf("ToggleUserStatusHandler | User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "User ID is required",
		})
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("ToggleUserStatusHandler | Invalid user ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid user ID format",
		})
	}

	if err := h.service.ToggleUserStatus(ctx, idParsed); err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			h.log.Warnf("ToggleUserStatusHandler | User not found: %v", id)
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, common.ErrInternal):
			h.log.Errorf("ToggleUserStatusHandler | Failed to toggle user status: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to toggle user status",
			})
		default:
			h.log.Errorf("ToggleUserStatusHandler | Unexpected error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Unexpected error occurred",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User status toggled successfully",
	})
}
