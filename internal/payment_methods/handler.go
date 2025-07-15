package payment_methods

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type IPaymentMethodHandler interface {
	ListPaymentMethodsHandler(c *fiber.Ctx) error
}

type PaymentMethodHandler struct {
	service IPaymentMethodService
	log     *logger.Logger
}

func NewPaymentMethodHandler(service IPaymentMethodService, log *logger.Logger) IPaymentMethodHandler {
	return &PaymentMethodHandler{service: service, log: log}
}

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
