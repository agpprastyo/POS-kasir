package printer

import (
	"POS-kasir/internal/common"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PrinterHandler struct {
	service IPrinterService
}

func NewPrinterHandler(service IPrinterService) *PrinterHandler {
	return &PrinterHandler{
		service: service,
	}
}

// PrintInvoiceHandler godoc
// @Summary Print invoice for an order
// @Description Trigger printing of invoice for a specific order
// @Tags Printer
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} common.SuccessResponse
// @Router /orders/{id}/print [post]
func (h *PrinterHandler) PrintInvoiceHandler(c *fiber.Ctx) error {
	idParam := c.Params("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid order ID",
			Error:   err.Error(),
		})
	}

	ctx := c.Context()
	err = h.service.PrintInvoice(ctx, orderID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to print invoice",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Invoice sent to printer",
	})
}

// TestPrintHandler godoc
// @Summary Test printer connection
// @Description Send a test print command to the configured printer
// @Tags Printer
// @Produce json
// @Success 200 {object} common.SuccessResponse
// @Router /settings/printer/test [post]
func (h *PrinterHandler) TestPrintHandler(c *fiber.Ctx) error {
	ctx := c.Context()
	err := h.service.TestPrint(ctx)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to send test print",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Test print command sent associated with configured printer",
	})
}
