package activitylog

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"

	"github.com/gofiber/fiber/v2"
)

type ActivityLogHandler struct {
	service   IActivityService
	log       logger.ILogger
	validator validator.Validator
}

func NewActivityLogHandler(service IActivityService, log logger.ILogger, validator validator.Validator) *ActivityLogHandler {
	return &ActivityLogHandler{
		service:   service,
		log:       log,
		validator: validator,
	}
}

// GetActivityLogs godoc
// @Summary      Get activity logs
// @Description  Get a list of activity logs with filtering and pagination
// @Tags         ActivityLogs
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        search query string false "Search term"
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Param        user_id query string false "User ID"
// @Param        entity_type query string false "Entity Type"
// @Param        action_type query string false "Action Type"
// @Success      200  {object}  common.SuccessResponse{data=dto.ActivityLogListResponse}
// @Failure      400  {object}  common.ErrorResponse
// @Failure      500  {object}  common.ErrorResponse
// @x-roles ["admin"]
// @Router       /activity-logs [get]
func (h *ActivityLogHandler) GetActivityLogs(c *fiber.Ctx) error {
	ctx := c.Context()

	var req dto.GetActivityLogsRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse query parameters",
			Error:   err.Error(),
		})
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 10
	}

	if done, err := common.ValidateAndRespond(c, h.validator, h.log, &req); done {
		return err
	}

	result, err := h.service.GetActivityLogs(ctx, req)
	if err != nil {
		h.log.Errorf("Handler | GetActivityLogs | %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to retrieve activity logs",
		})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data:    result,
	})
}
