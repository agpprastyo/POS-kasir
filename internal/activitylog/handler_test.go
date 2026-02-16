package activitylog_test

import (
	"POS-kasir/internal/activitylog"
	"POS-kasir/mocks"
	"POS-kasir/pkg/validator"
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

func setupHandlerTest(t *testing.T) (*mocks.MockIActivityService, *mocks.MockFieldLogger, *activitylog.ActivityLogHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIActivityService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	handler := activitylog.NewActivityLogHandler(mockService, mockLogger)
	app := fiber.New(fiber.Config{
		StructValidator: validator.NewValidator(),
	})
	return mockService, mockLogger, handler, app
}

func TestActivityLogHandler_GetActivityLogs(t *testing.T) {
	setup := func(t *testing.T) (*mocks.MockIActivityService, *mocks.MockFieldLogger, *activitylog.ActivityLogHandler, *fiber.App) {
		ctrl := gomock.NewController(t)
		mockService := mocks.NewMockIActivityService(ctrl)
		mockLogger := mocks.NewMockFieldLogger(ctrl)
		handler := activitylog.NewActivityLogHandler(mockService, mockLogger)
		app := fiber.New(fiber.Config{
			StructValidator: validator.NewValidator(),
		})
		app.Get("/activity-logs", handler.GetActivityLogs)
		return mockService, mockLogger, handler, app
	}

	t.Run("Success", func(t *testing.T) {
		mockService, _, _, app := setup(t)
		page := 1
		limit := 10
		req := activitylog.GetActivityLogsRequest{Page: &page, Limit: &limit}

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
		_, mockLogger, _, app := setup(t)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?page=abc", nil)
		resp, err := app.Test(reqHTTP)

		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("BadRequest_Validation", func(t *testing.T) {
		_, mockLogger, _, app := setup(t)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?limit=101", nil)
		resp, err := app.Test(reqHTTP)

		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("BadRequest_EmptyQueryParam", func(t *testing.T) {
		_, mockLogger, _, app := setup(t)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		// page= without value should be invalid if we want to enforce presence or valid format
		reqHTTP := httptest.NewRequest("GET", "/activity-logs?page=", nil)
		resp, err := app.Test(reqHTTP)

		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("BadRequest_InvalidUUID", func(t *testing.T) {
		_, mockLogger, _, app := setup(t)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?user_id=not-a-uuid", nil)
		resp, err := app.Test(reqHTTP)
		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("BadRequest_InvalidDate", func(t *testing.T) {
		_, mockLogger, _, app := setup(t)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()

		reqHTTP := httptest.NewRequest("GET", "/activity-logs?start_date=2023/01/01", nil)
		resp, err := app.Test(reqHTTP)
		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		}
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService, mockLogger, _, app := setup(t)
		mockService.EXPECT().GetActivityLogs(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		reqHTTP := httptest.NewRequest("GET", "/activity-logs", nil)
		resp, err := app.Test(reqHTTP)
		assert.NoError(t, err)
		if resp != nil {
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		}
	})
}
