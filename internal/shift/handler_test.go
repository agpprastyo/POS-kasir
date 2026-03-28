package shift_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/shift"
	"POS-kasir/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func allowAllLoggerCalls(mockLogger *mocks.MockILogger) {
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

func TestShiftHandler_StartShiftHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := shift.NewHandler(mockService, mockLogger)

	app := fiber.New()
	userID := uuid.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	app.Post("/shifts/start", handler.StartShiftHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := shift.StartShiftRequest{Password: "password123", StartCash: 1000}
		mockService.EXPECT().StartShift(gomock.Any(), userID, reqBody).Return(&shift.ShiftResponse{Status: "open"}, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/shifts/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/shifts/start", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := shift.StartShiftRequest{Password: "password123"}
		mockService.EXPECT().StartShift(gomock.Any(), userID, gomock.Any()).Return(nil, errors.New("service error"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/shifts/start", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestShiftHandler_EndShiftHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := shift.NewHandler(mockService, mockLogger)

	app := fiber.New()
	userID := uuid.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	app.Post("/shifts/end", handler.EndShiftHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := shift.EndShiftRequest{Password: "password123", ActualCashEnd: 1500}
		mockService.EXPECT().EndShift(gomock.Any(), userID, reqBody).Return(&shift.ShiftResponse{Status: "closed"}, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/shifts/end", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		reqBody := shift.EndShiftRequest{Password: "password123"}
		mockService.EXPECT().EndShift(gomock.Any(), userID, gomock.Any()).Return(nil, errors.New("service error"))

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/shifts/end", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestShiftHandler_GetOpenShiftHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := shift.NewHandler(mockService, mockLogger)

	app := fiber.New()
	userID := uuid.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	app.Get("/shifts/current", handler.GetOpenShiftHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().GetOpenShift(gomock.Any(), userID).Return(&shift.ShiftResponse{Status: "open"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/shifts/current", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService.EXPECT().GetOpenShift(gomock.Any(), userID).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/shifts/current", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("GenericError", func(t *testing.T) {
		mockService.EXPECT().GetOpenShift(gomock.Any(), userID).Return(nil, errors.New("generic error"))

		req := httptest.NewRequest(http.MethodGet, "/shifts/current", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestShiftHandler_CreateCashTransactionHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := shift.NewHandler(mockService, mockLogger)

	app := fiber.New()
	userID := uuid.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("user_id", userID)
		return c.Next()
	})
	app.Post("/shifts/cash-transaction", handler.CreateCashTransactionHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := shift.CashTransactionRequest{Amount: 500, Type: "cash_in"}
		mockService.EXPECT().CreateCashTransaction(gomock.Any(), userID, reqBody).Return(&shift.CashTransactionResponse{Amount: 500}, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/shifts/cash-transaction", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService.EXPECT().CreateCashTransaction(gomock.Any(), userID, gomock.Any()).Return(nil, errors.New("service error"))

		body, _ := json.Marshal(shift.CashTransactionRequest{Amount: 100})
		req := httptest.NewRequest(http.MethodPost, "/shifts/cash-transaction", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
