package printer

import (
	"POS-kasir/internal/common"
	"net/http"

	"github.com/gofiber/fiber/v3"
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
// @Summary      Print invoice for an order
// @Description  Trigger printing of invoice for a specific order (Roles: admin, manager, cashier)
// @Tags         Printer
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "Invoice sent to printer"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID"
// @Failure      500 {object} common.ErrorResponse "Failed to print invoice"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/print [post]
func (h *PrinterHandler) PrintInvoiceHandler(c fiber.Ctx) error {
	idParam := c.Params("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid order ID",
			Error:   err.Error(),
		})
	}

	ctx := c.RequestCtx()
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
// @Summary      Test printer connection
// @Description  Send a test print command to the configured printer (Roles: admin)
// @Tags         Printer
// @Accept       json
// @Produce      json
// @Success      200 {object} common.SuccessResponse "Test print command sent associated with configured printer"
// @Failure      500 {object} common.ErrorResponse "Failed to send test print"
// @x-roles      ["admin"]
// @Router       /settings/printer/test [post]
func (h *PrinterHandler) TestPrintHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
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

// GetInvoiceDataHandler godoc
// @Summary      Get invoice print data
// @Description  Get raw invoice print data for FE printing (Roles: admin, manager, cashier)
// @Tags         Printer
// @Accept       json
// @Produce      json
// @Param        id path string true "Order ID" Format(uuid)
// @Success      200 {object} common.SuccessResponse "Invoice print data"
// @Failure      400 {object} common.ErrorResponse "Invalid order ID"
// @Failure      500 {object} common.ErrorResponse "Failed to generate print data"
// @x-roles      ["admin", "manager", "cashier"]
// @Router       /orders/{id}/print-data [get]
func (h *PrinterHandler) GetInvoiceDataHandler(c fiber.Ctx) error {
	idParam := c.Params("id")
	orderID, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Invalid order ID",
			Error:   err.Error(),
		})
	}

	ctx := c.RequestCtx()
	data, filename, err := h.service.GetInvoiceData(ctx, orderID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to generate print data",
			Error:   err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(common.SuccessResponse{
		Message: "Invoice print data generated",
		Data: fiber.Map{
			"data":     data,
			"filename": filename,
		},
	})
}
