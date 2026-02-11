package activitylog

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/dto"
	"POS-kasir/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(t *testing.T) (*mocks.MockIActivityService, *mocks.MockFieldLogger, *mocks.MockValidator, *ActivityLogHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)
	handler := NewActivityLogHandler(mockService, mockLogger, mockValidator)
	app := fiber.New()
	return mockService, mockLogger, mockValidator, handler, app
}

func TestActivityLogHandler_GetActivityLogs(t *testing.T) {
	mockService, mockLogger, mockValidator, handler, app := setupHandlerTest(t)
	app.Get("/activity-logs", handler.GetActivityLogs)

	t.Run("Success", func(t *testing.T) {
		req := dto.GetActivityLogsRequest{Page: 1, Limit: 10}
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)

		expectedData := &dto.ActivityLogListResponse{
			Logs: []dto.ActivityLogResponse{
				{ID: uuid.New(), UserName: "Admin", ActionType: "CREATE", CreatedAt: time.Now()},
			},
			TotalItems: 1,
			Page:       1,
			Limit:      10,
			TotalPages: 1,
		}

		mockService.EXPECT().GetActivityLogs(gomock.Any(), req).Return(expectedData, nil)

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?page=1&limit=10", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var successResp common.SuccessResponse
		json.NewDecoder(resp.Body).Decode(&successResp)
		assert.Equal(t, "Success", successResp.Message)
	})

	t.Run("BadRequest_QueryParser", func(t *testing.T) {
		reqHTTP := httptest.NewRequest("GET", "/activity-logs?page=abc", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("BadRequest_Validation", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation failed"))

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?page=1", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().GetActivityLogs(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		reqHTTP := httptest.NewRequest("GET", "/activity-logs", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
