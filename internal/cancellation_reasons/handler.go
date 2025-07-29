package cancellation_reasons

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type ICancellationReasonHandler interface {
	ListCancellationReasonsHandler(c *fiber.Ctx) error
}

type CancellationReasonHandler struct {
	service ICancellationReasonService
	log     logger.ILogger
}

func NewCancellationReasonHandler(service ICancellationReasonService, log logger.ILogger) ICancellationReasonHandler {
	return &CancellationReasonHandler{service: service, log: log}
}

func (h *CancellationReasonHandler) ListCancellationReasonsHandler(c *fiber.Ctx) error {
	reasons, err := h.service.ListCancellationReasons(c.Context())
	if err != nil {
		h.log.Error("ListCancellationReasonsHandler | Failed to get cancellation reasons from service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve cancellation reasons"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Cancellation reasons retrieved successfully",
		Data:    reasons,
	})
}
