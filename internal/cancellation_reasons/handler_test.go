package cancellation_reasons_test

import (
	"POS-kasir/internal/cancellation_reasons"
	"POS-kasir/internal/common"
	"POS-kasir/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(t *testing.T) (*mocks.MockICancellationReasonService, *mocks.MockFieldLogger, *cancellation_reasons.CancellationReasonHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockICancellationReasonService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	handler := cancellation_reasons.NewCancellationReasonHandler(mockService, mockLogger).(*cancellation_reasons.CancellationReasonHandler)
	app := fiber.New()
	return mockService, mockLogger, handler, app
}

func TestCancellationReasonHandler_GetActivityLogs(t *testing.T) {
	mockService, mockLogger, handler, app := setupHandlerTest(t)
	app.Get("/cancellation-reasons", handler.ListCancellationReasonsHandler)

	t.Run("Success", func(t *testing.T) {
		reasons := []cancellation_reasons.CancellationReasonResponse{
			{ID: 1, Reason: "Reason 1", IsActive: true, CreatedAt: time.Now()},
		}
		mockService.EXPECT().ListCancellationReasons(gomock.Any()).Return(reasons, nil)

		req := httptest.NewRequest("GET", "/cancellation-reasons", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result common.SuccessResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Cancellation reasons retrieved successfully", result.Message)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService.EXPECT().ListCancellationReasons(gomock.Any()).Return(nil, errors.New("service error"))
		mockLogger.EXPECT().Error(gomock.Any(), gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/cancellation-reasons", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var result common.ErrorResponse
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "Failed to retrieve cancellation reasons", result.Message)
	})
}
