package user

import (
	"POS-kasir/internal/common"
	user_repo "POS-kasir/internal/user/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func NewUsrHandler(service IUsrService, log logger.ILogger, validator validator.Validator) IUsrHandler {
	return &UsrHandler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

type IUsrHandler interface {
	GetAllUsersHandler(c fiber.Ctx) error
	CreateUserHandler(c fiber.Ctx) error
	GetUserByIDHandler(c fiber.Ctx) error
	UpdateUserHandler(c fiber.Ctx) error
	ToggleUserStatusHandler(c fiber.Ctx) error
	DeleteUserHandler(c fiber.Ctx) error
}

type UsrHandler struct {
	service   IUsrService
	log       logger.ILogger
	validator validator.Validator
}

// DeleteUserHandler handles the request to delete a user by ID
// @Summary      Delete user
// @Description  Hard delete a user from the system by their ID (Roles: admin)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "User deleted successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid user ID format"
// @Failure      404 {object} common.ErrorResponse "User not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin"]
// @Router       /users/{id} [delete]
func (h *UsrHandler) DeleteUserHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	id := c.Params("id")
	if id == "" {
		h.log.Errorf("DeleteUserHandler | User ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "User ID is required",
		})
	}

	idParsed, err := uuid.Parse(id)
	if err != nil {
		h.log.Errorf("DeleteUserHandler | Invalid user ID format: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid user ID format",
		})
	}

	if err := h.service.DeleteUser(ctx, idParsed); err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			h.log.Warnf("DeleteUserHandler | User not found: %v", id)
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		default:
			h.log.Errorf("DeleteUserHandler | Failed to delete user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to delete user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User deleted successfully",
	})
}

// GetAllUsersHandler handles the request to get all users
// @Summary      Get all users
// @Description  Retrieve a list of users with pagination, filtering, and sorting (Roles: admin, manager, cashier)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        page      query int     false "Page number (default 1)" default(1)
// @Param        limit     query int     false "Items per page (default 10)" default(10)
// @Param        search    query string  false "Search by username or email"
// @Param        role      query string  false "Filter by User Role" Enums(admin, cashier, manager)
// @Param        is_active query boolean false "Filter by Active Status"
// @Param        status    query string  false "Filter by Account Status" Enums(active, deleted, all)
// @Param        sortBy    query string  false "Sort by column" Enums(created_at, username)
// @Param        sortOrder query string  false "Sort direction" Enums(asc, desc)
// @Success      200 {object} common.SuccessResponse{data=UsersResponse} "Users retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /users [get]
func (h *UsrHandler) GetAllUsersHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	req := new(UsersRequest)
	if err := c.Bind().Query(req); err != nil {
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
			return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
				Message: "Users retrieved successfully",
				Data:    response,
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
// @Summary      Create user
// @Description  Create a new user account (Roles: admin)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "New user details"
// @Success      201 {object} common.SuccessResponse{data=ProfileResponse} "User created successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body or validation failed"
// @Failure      409 {object} common.ErrorResponse "User, username, or email already exists"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin"]
// @Router       /users [post]
func (h *UsrHandler) CreateUserHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	req := new(CreateUserRequest)
	if err := c.Bind().Body(req); err != nil {
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
// @Summary      Get user by ID
// @Description  Retrieve detailed profile for a specific user by their ID (Roles: admin, manager)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse{data=ProfileResponse} "User retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid user ID format"
// @Failure      404 {object} common.ErrorResponse "User not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager"]
// @Router       /users/{id} [get]
func (h *UsrHandler) GetUserByIDHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
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
// @Summary      Update user
// @Description  Update details of an existing user account (Roles: admin)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id   path string true "User ID" Format(uuid)
// @Param        user body UpdateUserRequest true "User update details"
// @Success      200 {object} common.SuccessResponse{data=ProfileResponse} "User updated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request body or validation failed"
// @Failure      403 {object} common.ErrorResponse "Forbidden - higher role assignment attempt"
// @Failure      404 {object} common.ErrorResponse "User not found"
// @Failure      409 {object} common.ErrorResponse "Username or email already exists"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin"]
// @Router       /users/{id} [put]
func (h *UsrHandler) UpdateUserHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
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
	if err := c.Bind().Body(req); err != nil {
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
		if role != user_repo.UserRoleAdmin {
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

	h.log.Infof("UpdateUserHandler | User updated successfully: %v", user.Username)

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User updated successfully",
		Data:    user,
	})
}

// ToggleUserStatusHandler handles the request to toggle a user's active status
// @Summary      Toggle user status
// @Description  Toggle the is_active status of a user (Roles: admin)
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "User status toggled successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid user ID format"
// @Failure      404 {object} common.ErrorResponse "User not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin"]
// @Router       /users/{id}/toggle-status [post]
func (h *UsrHandler) ToggleUserStatusHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
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

// fiber:context-methods migrated
