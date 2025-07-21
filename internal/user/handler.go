package user

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

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
		h.log.Errorf("User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Parse the ID to UUID
	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("Invalid user ID format", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	if err := h.service.DeleteUser(ctx, idParsed); err != nil {
		h.log.Errorf("Failed to delete user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User deleted successfully",
	})
}

func NewUsrHandler(service IUsrService, log logger.ILogger, validator validator.Validator) IUsrHandler {
	return &UsrHandler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

// GetAllUsersHandler handles the request to get all users
func (h *UsrHandler) GetAllUsersHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	req := new(UsersRequest)
	if err := c.QueryParser(req); err != nil {
		h.log.Errorf("Failed to parse query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	isActiveStr := c.Query("is_active")
	if isActiveStr != "" {
		parsedBool, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			req.IsActive = &parsedBool
		}
	}

	h.log.Infof("DTO after manual parse:", "dto", req)

	h.log.Infof("Get all users 1", "dto", req)

	if err := h.validator.Validate(req); err != nil {
		h.log.Errorf("Validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	h.log.Infof("Get all users 2", "dto", req)

	response, err := h.service.GetAllUsers(ctx, *req)
	if err != nil {
		h.log.Errorf("Failed to get all users", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	h.log.Infof("Get all users 3", "response", response)

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
		h.log.Errorf("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Errorf("Validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	user, err := h.service.CreateUser(ctx, *req)
	if err != nil {
		h.log.Errorf("Failed to create user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
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
		h.log.Errorf("User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}
	// Parse the ID to UUID
	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("Invalid user ID format", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	user, err := h.service.GetUserByID(ctx, idParsed)
	if err != nil {
		h.log.Errorf("Failed to get user by ID", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get user",
		})
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
		h.log.Errorf("User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Parse the ID to UUID
	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("Invalid user ID format", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	req := new(UpdateUserRequest)
	if err := c.BodyParser(req); err != nil {
		h.log.Errorf("Failed to parse request body", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Errorf("Validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	if req.Role != nil {
		if role != repository.UserRoleAdmin {
			h.log.Errorf("Unauthorized role change attempt", "role", role)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Must be an admin to change user roles",
			})
		}
	}

	user, err := h.service.UpdateUser(ctx, idParsed, *req)
	if err != nil {
		h.log.Errorf("Failed to update user", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
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
		h.log.Errorf("User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Parse the ID to UUID
	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("Invalid user ID format", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	if err := h.service.ToggleUserStatus(ctx, idParsed); err != nil {
		h.log.Errorf("Failed to toggle user status", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to toggle user status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User status toggled successfully",
	})
}
