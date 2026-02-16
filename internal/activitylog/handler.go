package activitylog

import (
	"POS-kasir/internal/common"
	"errors"

	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/validator"

	"github.com/gofiber/fiber/v3"
)

type ActivityLogHandler struct {
	service IActivityService
	log     logger.ILogger
}

func NewActivityLogHandler(service IActivityService, log logger.ILogger) *ActivityLogHandler {
	return &ActivityLogHandler{
		service: service,
		log:     log,
	}
}

// GetActivityLogs godoc
// @Summary      Get activity logs
// @Description  Get a list of activity logs with filtering and pagination (Roles: admin)
// @Tags         ActivityLogs
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number" default(1)
// @Param        limit query int false "Items per page" default(10)
// @Param        search query string false "Search term"
// @Param        start_date query string false "Start date (YYYY-MM-DD)"
// @Param        end_date query string false "End date (YYYY-MM-DD)"
// @Param        user_id query string false "User ID" Format(uuid)
// @Param        entity_type query string false "Entity Type" Enums(PRODUCT, CATEGORY, PROMOTION, ORDER, USER)
// @Param        action_type query string false "Action Type" Enums(CREATE, UPDATE, DELETE, CANCEL, APPLY_PROMOTION, PROCESS_PAYMENT, REGISTER, UPDATE_PASSWORD, UPDATE_AVATAR, LOGIN_SUCCESS, LOGIN_FAILED)
// @Success      200  {object}  common.SuccessResponse{data=ActivityLogListResponse}
// @Failure      400  {object}  common.ErrorResponse
// @Failure      500  {object}  common.ErrorResponse
// @x-roles ["admin"]
// @Router       /activity-logs [get]
func (h *ActivityLogHandler) GetActivityLogs(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	var req GetActivityLogsRequest
	if err := c.Bind().Query(&req); err != nil {
		h.log.Errorf("Handler | GetActivityLogs | %v", err)
		var ve *validator.ValidationErrors
		if errors.As(err, &ve) {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Validation failed",
				Error:   ve.Error(),
				Data: map[string]interface{}{
					"errors": ve.Errors,
				},
			})
		}
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse query parameters",
			Error:   err.Error(),
		})
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
