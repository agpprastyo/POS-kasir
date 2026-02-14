package activitylog_test

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(t *testing.T) (*mocks.MockIActivityService, *mocks.MockFieldLogger, *mocks.MockValidator, *activitylog.ActivityLogHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)
	handler := activitylog.NewActivityLogHandler(mockService, mockLogger, mockValidator)
	app := fiber.New()
	return mockService, mockLogger, mockValidator, handler, app
}

func TestActivityLogHandler_GetActivityLogs(t *testing.T) {
	setup := func(t *testing.T) (*mocks.MockIActivityService, *mocks.MockFieldLogger, *mocks.MockValidator, *activitylog.ActivityLogHandler, *fiber.App) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockIActivityService(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		mockValidator := mocks.NewMockValidator(ctrl)
		handler := activitylog.NewActivityLogHandler(mockService, mockLogger, mockValidator)
		app := fiber.New()
		app.Get("/activity-logs", handler.GetActivityLogs)
		return mockService, mockLogger, mockValidator, handler, app
	}

	t.Run("Success", func(t *testing.T) {
		mockService, _, mockValidator, _, app := setup(t)
		req := activitylog.GetActivityLogsRequest{Page: 1, Limit: 10}
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)

		expectedData := &activitylog.ActivityLogListResponse{
			Logs: []activitylog.ActivityLogResponse{
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
	})

	t.Run("BadRequest_QueryParser", func(t *testing.T) {
		_, mockLogger, _, _, app := setup(t)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?page=abc", nil)
		resp, err := app.Test(reqHTTP)

		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("BadRequest_Validation", func(t *testing.T) {
		_, _, mockValidator, _, app := setup(t)
		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation failed"))

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?page=0", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		_, _, mockValidator, _, app := setup(t)
		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("invalid uuid"))

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?user_id=not-a-uuid", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("BadRequest_InvalidDate", func(t *testing.T) {
		_, _, mockValidator, _, app := setup(t)
		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("invalid date"))

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?start_date=2023/01/01", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService, mockLogger, mockValidator, _, app := setup(t)
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().GetActivityLogs(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		reqHTTP := httptest.NewRequest("GET", "/activity-logs", nil)
		resp, _ := app.Test(reqHTTP)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
