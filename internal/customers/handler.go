package customers

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"errors"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ICustomerHandler interface {
	CreateCustomerHandler(c fiber.Ctx) error
	GetCustomerHandler(c fiber.Ctx) error
	UpdateCustomerHandler(c fiber.Ctx) error
	DeleteCustomerHandler(c fiber.Ctx) error
	ListCustomersHandler(c fiber.Ctx) error
}

type CustomerHandler struct {
	service ICustomerService
	log     logger.ILogger
}

func NewCustomerHandler(service ICustomerService, log logger.ILogger) ICustomerHandler {
	return &CustomerHandler{service: service, log: log}
}

// CreateCustomerHandler creates a new customer
// @Summary      Create a customer
// @Description  Create a new customer (Roles: admin, manager, cashier)
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        request body CreateCustomerRequest true "Customer details"
// @Success      201 {object} common.SuccessResponse{data=CustomerResponse} "Customer created successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /customers [post]
func (h *CustomerHandler) CreateCustomerHandler(c fiber.Ctx) error {
	var req CreateCustomerRequest
	if err := c.Bind().Body(&req); err != nil {
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	resp, err := h.service.CreateCustomer(c.RequestCtx(), req)
	if err != nil {
		h.log.Errorf("CreateCustomer failed", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to create customer"})
	}

	return c.Status(fiber.StatusCreated).JSON(common.SuccessResponse{
		Message: "Customer created successfully",
		Data:    resp,
	})
}

// GetCustomerHandler retrieves a customer by ID
// @Summary      Get a customer
// @Description  Get customer by ID
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Success      200 {object} common.SuccessResponse{data=CustomerResponse} "Customer retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format"
// @Failure      404 {object} common.ErrorResponse "Customer not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /customers/{id} [get]
func (h *CustomerHandler) GetCustomerHandler(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid ID format"})
	}

	resp, err := h.service.GetCustomer(c.RequestCtx(), id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to fetch customer"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Customer retrieved successfully",
		Data:    resp,
	})
}

// UpdateCustomerHandler updates a customer by ID
// @Summary      Update a customer
// @Description  Update customer by ID
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Param        request body UpdateCustomerRequest true "Customer details to update"
// @Success      200 {object} common.SuccessResponse{data=CustomerResponse} "Customer updated successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid request"
// @Failure      404 {object} common.ErrorResponse "Customer not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager"]
// @Router       /customers/{id} [put]
func (h *CustomerHandler) UpdateCustomerHandler(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid ID format"})
	}

	var req UpdateCustomerRequest
	if err := c.Bind().Body(&req); err != nil {
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data:    map[string]interface{}{"errors": ve.Errors},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid request body"})
	}

	resp, err := h.service.UpdateCustomer(c.RequestCtx(), id, req)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to update customer"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Customer updated successfully",
		Data:    resp,
	})
}

// DeleteCustomerHandler deletes a customer by ID
// @Summary      Delete a customer
// @Description  Delete customer by ID
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        id path string true "Customer ID"
// @Success      200 {object} common.SuccessResponse "Customer deleted successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid ID format"
// @Failure      404 {object} common.ErrorResponse "Customer not found"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager"]
// @Router       /customers/{id} [delete]
func (h *CustomerHandler) DeleteCustomerHandler(c fiber.Ctx) error {
	id, err := fiber.Convert(c.Params("id"), uuid.Parse)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid ID format"})
	}

	err = h.service.DeleteCustomer(c.RequestCtx(), id)
	if err != nil {
		if errors.Is(err, common.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{Message: "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to delete customer"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Customer deleted successfully",
	})
}

// ListCustomersHandler retrieves a list of customers
// @Summary      List customers
// @Description  List customers with pagination and search
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Page size limit"
// @Param        search query string false "Search by name, phone, or email"
// @Success      200 {object} common.SuccessResponse{data=PagedCustomerResponse} "Customers retrieved successfully"
// @Failure      400 {object} common.ErrorResponse "Invalid query parameters"
// @Failure      500 {object} common.ErrorResponse "Internal server error"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /customers [get]
func (h *CustomerHandler) ListCustomersHandler(c fiber.Ctx) error {
	var req ListCustomersRequest
	if err := c.Bind().Query(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{Message: "Invalid query parameters"})
	}

	resp, err := h.service.ListCustomers(c.RequestCtx(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to list customers"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Customers retrieved successfully",
		Data:    resp,
	})
}
