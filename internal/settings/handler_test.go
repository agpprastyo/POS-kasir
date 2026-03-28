package settings_test

import (
	"POS-kasir/internal/settings"
	"POS-kasir/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func allowAllLoggerCalls(mockLogger *mocks.MockILogger) {
	mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Info(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Error(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
}

func TestSettingsHandler_GetBrandingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISettingsService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := settings.NewSettingsHandler(mockService, mockLogger)

	app := fiber.New()
	app.Get("/settings/branding", handler.GetBrandingHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().GetBranding(gomock.Any()).Return(&settings.BrandingSettingsResponse{AppName: "Test"}, nil)

		req := httptest.NewRequest(http.MethodGet, "/settings/branding", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Error", func(t *testing.T) {
		mockService.EXPECT().GetBranding(gomock.Any()).Return(nil, errors.New("error"))

		req := httptest.NewRequest(http.MethodGet, "/settings/branding", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestSettingsHandler_UpdateBrandingHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISettingsService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := settings.NewSettingsHandler(mockService, mockLogger)

	app := fiber.New()
	app.Put("/settings/branding", handler.UpdateBrandingHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := settings.UpdateBrandingRequest{AppName: "New Name"}
		mockService.EXPECT().UpdateBranding(gomock.Any(), reqBody).Return(&settings.BrandingSettingsResponse{AppName: "New Name"}, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/settings/branding", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/settings/branding", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSettingsHandler_UpdateLogoHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISettingsService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := settings.NewSettingsHandler(mockService, mockLogger)

	app := fiber.New(fiber.Config{
		BodyLimit: 10 * 1024 * 1024,
	})
	app.Post("/settings/branding/logo", handler.UpdateLogoHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().UpdateLogo(gomock.Any(), gomock.Any(), "logo.png", gomock.Any()).Return("http://url", nil)

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("logo", "logo.png")
		part.Write([]byte("fake image data"))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/settings/branding/logo", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NoFile", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/settings/branding/logo", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("FileTooLarge", func(t *testing.T) {
		
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("logo", "large.png")
		// 6MB > 5MB limit
		part.Write(make([]byte, 6*1024*1024))
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/settings/branding/logo", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, err := app.Test(req)
		
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSettingsHandler_GetPrinterSettingsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISettingsService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := settings.NewSettingsHandler(mockService, mockLogger)

	app := fiber.New()
	app.Get("/settings/printer", handler.GetPrinterSettingsHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().GetPrinterSettings(gomock.Any()).Return(&settings.PrinterSettingsResponse{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/settings/printer", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestSettingsHandler_UpdatePrinterSettingsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mocks.NewMockISettingsService(ctrl)
	mockLogger := mocks.NewMockILogger(ctrl)
	allowAllLoggerCalls(mockLogger)
	handler := settings.NewSettingsHandler(mockService, mockLogger)

	app := fiber.New()
	app.Put("/settings/printer", handler.UpdatePrinterSettingsHandler)

	t.Run("Success", func(t *testing.T) {
		reqBody := settings.UpdatePrinterSettingsRequest{Connection: "new"}
		mockService.EXPECT().UpdatePrinterSettings(gomock.Any(), reqBody).Return(&settings.PrinterSettingsResponse{}, nil)

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPut, "/settings/printer", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
