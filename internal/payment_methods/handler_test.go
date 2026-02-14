package payment_methods

import (
	"POS-kasir/internal/dto"
	"POS-kasir/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPaymentMethodHandler_ListPaymentMethodsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockIPaymentMethodService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	handler := NewPaymentMethodHandler(mockService, mockLogger)

	app := fiber.New()
	app.Get("/payment-methods", handler.ListPaymentMethodsHandler)

	t.Run("Success", func(t *testing.T) {
		methods := []dto.PaymentMethodResponse{
			{
				ID:        1,
				Name:      "Cash",
				IsActive:  true,
				CreatedAt: time.Now().Truncate(time.Second),
			},
			{
				ID:        2,
				Name:      "QRIS",
				IsActive:  true,
				CreatedAt: time.Now().Truncate(time.Second),
			},
		}

		mockService.EXPECT().ListPaymentMethods(gomock.Any()).Return(methods, nil)

		req := httptest.NewRequest("GET", "/payment-methods", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)

		assert.Equal(t, "Payment methods retrieved successfully", result["message"])
		data := result["data"].([]interface{})
		assert.Len(t, data, 2)
		assert.Equal(t, "Cash", data[0].(map[string]interface{})["name"])
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService.EXPECT().ListPaymentMethods(gomock.Any()).Return(nil, errors.New("service error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/payment-methods", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
