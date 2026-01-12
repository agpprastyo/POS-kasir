package payment_methods

import (
	"POS-kasir/internal/common"
	_ "POS-kasir/internal/dto"
	"POS-kasir/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type IPaymentMethodHandler interface {
	ListPaymentMethodsHandler(c *fiber.Ctx) error
}

type PaymentMethodHandler struct {
	service IPaymentMethodService
	log     logger.ILogger
}

func NewPaymentMethodHandler(service IPaymentMethodService, log logger.ILogger) IPaymentMethodHandler {
	return &PaymentMethodHandler{service: service, log: log}
}

// ListPaymentMethodsHandler handles list payment methods requests.
// @Summary List payment methods
// @Description List payment methods
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse{data=[]dto.PaymentMethodResponse} "Success"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /payment-methods [get]
func (h *PaymentMethodHandler) ListPaymentMethodsHandler(c *fiber.Ctx) error {
	methods, err := h.service.ListPaymentMethods(c.Context())
	if err != nil {
		h.log.Error("Failed to get payment methods from service", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{Message: "Failed to retrieve payment methods"})
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Payment methods retrieved successfully",
		Data:    methods,
	})
}
