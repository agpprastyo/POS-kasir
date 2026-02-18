package printer_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/printer"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPrinterService
type MockPrinterService struct {
	mock.Mock
}

func (m *MockPrinterService) PrintInvoice(ctx context.Context, orderID uuid.UUID) error {
	args := m.Called(ctx, orderID)
	return args.Error(0)
}

func (m *MockPrinterService) TestPrint(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func TestPrinterHandler_PrintInvoiceHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app := fiber.New()
		mockService := new(MockPrinterService)
		handler := printer.NewPrinterHandler(mockService)
		app.Post("/orders/:id/print", handler.PrintInvoiceHandler)

		orderID := uuid.New()
		mockService.On("PrintInvoice", mock.Anything, orderID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/orders/"+orderID.String()+"/print", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var successResp common.SuccessResponse
		json.NewDecoder(resp.Body).Decode(&successResp)
		assert.Equal(t, "Invoice sent to printer", successResp.Message)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		app := fiber.New()
		mockService := new(MockPrinterService)
		handler := printer.NewPrinterHandler(mockService)
		app.Post("/orders/:id/print", handler.PrintInvoiceHandler)

		req := httptest.NewRequest(http.MethodPost, "/orders/invalid-uuid/print", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		app := fiber.New()
		mockService := new(MockPrinterService)
		handler := printer.NewPrinterHandler(mockService)
		app.Post("/orders/:id/print", handler.PrintInvoiceHandler)

		orderID := uuid.New()
		mockService.On("PrintInvoice", mock.Anything, orderID).Return(errors.New("printer error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/orders/"+orderID.String()+"/print", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestPrinterHandler_TestPrintHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		app := fiber.New()
		mockService := new(MockPrinterService)
		handler := printer.NewPrinterHandler(mockService)
		app.Post("/settings/printer/test", handler.TestPrintHandler)

		mockService.On("TestPrint", mock.Anything).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/settings/printer/test", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		app := fiber.New()
		mockService := new(MockPrinterService)
		handler := printer.NewPrinterHandler(mockService)
		app.Post("/settings/printer/test", handler.TestPrintHandler)

		mockService.On("TestPrint", mock.Anything).Return(errors.New("printer error")).Once()

		req := httptest.NewRequest(http.MethodPost, "/settings/printer/test", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
