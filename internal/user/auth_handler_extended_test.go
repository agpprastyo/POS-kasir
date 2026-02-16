package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"POS-kasir/config"
	"POS-kasir/internal/common"
	"POS-kasir/internal/user"
	"POS-kasir/mocks"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupExtendedAuthHandlerTest(t *testing.T) (*mocks.MockIAuthService, *mocks.MockFieldLogger, *mocks.MockValidator, user.IAuthHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIAuthService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)
	cfg := &config.AppConfig{}

	handler := user.NewAuthHandler(mockService, mockLogger, mockValidator, cfg)
	app := fiber.New()

	return mockService, mockLogger, mockValidator, handler, app
}

func TestAuthHandler_Login_Extended(t *testing.T) {
	reqBody := user.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}
	body, _ := json.Marshal(reqBody)

	t.Run("InvalidBody", func(t *testing.T) {
		_, _, _, handler, app := setupExtendedAuthHandlerTest(t)
		app.Post("/auth/login", handler.LoginHandler)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationFailed", func(t *testing.T) {
		_, _, mockValidator, handler, app := setupExtendedAuthHandlerTest(t)
		app.Post("/auth/login", handler.LoginHandler)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation failed"))

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_Internal", func(t *testing.T) {
		mockService, mockLogger, mockValidator, handler, app := setupExtendedAuthHandlerTest(t)
		app.Post("/auth/login", handler.LoginHandler)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().Login(gomock.Any(), reqBody).Return(nil, errors.New("internal error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("ServiceError_InvalidCredentials", func(t *testing.T) {
		mockService, _, mockValidator, handler, app := setupExtendedAuthHandlerTest(t)
		app.Post("/auth/login", handler.LoginHandler)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().Login(gomock.Any(), reqBody).Return(nil, common.ErrInvalidCredentials)

		req := httptest.NewRequest("POST", "/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestAuthHandler_UpdateAvatar_Extended(t *testing.T) {
	userID := uuid.New()

	setupApp := func(t *testing.T) (*mocks.MockIAuthService, *mocks.MockFieldLogger, *fiber.App) {
		mockService, mockLogger, _, handler, app := setupExtendedAuthHandlerTest(t)
		app.Post("/auth/avatar", func(c fiber.Ctx) error {
			c.Locals("user_id", userID)
			return handler.UpdateAvatarHandler(c)
		})
		return mockService, mockLogger, app
	}

	t.Run("NoFile", func(t *testing.T) {
		_, mockLogger, app := setupApp(t)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		// Create request without file
		req := httptest.NewRequest("POST", "/auth/avatar", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_FileTooLarge", func(t *testing.T) {
		mockService, _, app := setupApp(t)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("avatar", "test.jpg")
		part.Write([]byte("image content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/auth/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		// Mock expectations
		mockService.EXPECT().UploadAvatar(gomock.Any(), userID, gomock.Any()).Return(nil, common.ErrFileTooLarge)

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_Internal", func(t *testing.T) {
		mockService, mockLogger, app := setupApp(t)

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("avatar", "test.jpg")
		part.Write([]byte("image content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/auth/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		mockService.EXPECT().UploadAvatar(gomock.Any(), userID, gomock.Any()).Return(nil, errors.New("processing failed"))
		// Fixed expectation: 2 args (format, arg)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
