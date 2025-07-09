package user

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type UsrHandler struct {
	service   IUsrService
	log       *logger.Logger
	validator validator.Validator
}

func NewUsrHandler(service IUsrService, log *logger.Logger, validator validator.Validator) *UsrHandler {
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
		h.log.Error("Failed to parse query parameters", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if err := h.validator.Validate(req); err != nil {
		h.log.Error("Validation failed", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	response, err := h.service.GetAllUsers(ctx, *req)
	if err != nil {
		h.log.Error("Failed to get all users", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Users retrieved successfully",
		Data:    response,
	})

}
