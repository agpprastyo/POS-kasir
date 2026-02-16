package user_test

import (
	"POS-kasir/config"
	"POS-kasir/internal/common"
	"POS-kasir/internal/user"
	user_repo "POS-kasir/internal/user/repository"

	"POS-kasir/mocks"
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupHandlerTest(t *testing.T) (*mocks.MockIAuthService, *mocks.MockFieldLogger, *mocks.MockValidator, user.IAuthHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIAuthService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	cfg := &config.AppConfig{
		Server: config.ServerConfig{
			Env:          "test",
			CookieDomain: "localhost",
		},
		JWT: config.JwtConfig{
			RefreshTokenDuration: 24 * time.Hour,
		},
	}

	handler := user.NewAuthHandler(mockService, mockLogger, mockValidator, cfg)
	app := fiber.New()

	return mockService, mockLogger, mockValidator, handler, app
}

func TestAthHandler_LoginHandler(t *testing.T) {
	mockService, _, mockValidator, handler, app := setupHandlerTest(t)
	app.Post("/login", handler.LoginHandler)

	reqBody := user.LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}
	body, _ := json.Marshal(reqBody)

	t.Run("Success", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)

		loginResp := &user.LoginResponse{
			Token:        "access_token",
			RefreshToken: "refresh_token",
			ExpiredAt:    time.Now().Add(time.Hour),
			Profile: user.ProfileResponse{
				Email:    reqBody.Email,
				Username: "testuser",
			},
		}

		mockService.EXPECT().Login(gomock.Any(), reqBody).Return(loginResp, nil)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().Login(gomock.Any(), reqBody).Return(nil, common.ErrInvalidCredentials)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().Login(gomock.Any(), reqBody).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}

func TestAthHandler_ProfileHandler(t *testing.T) {
	mockService, mockLogger, _, handler, app := setupHandlerTest(t)

	// Inject middleware to simulate context locals
	app.Get("/me", func(c fiber.Ctx) error {
		c.Locals("user_id", uuid.MustParse("00000000-0000-0000-0000-000000000001"))
		return handler.ProfileHandler(c)
	})

	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Success", func(t *testing.T) {
		profileResp := &user.ProfileResponse{
			ID:       userID,
			Username: "testuser",
			Email:    "test@example.com",
		}
		mockService.EXPECT().Profile(gomock.Any(), userID).Return(profileResp, nil)

		req := httptest.NewRequest("GET", "/me", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockService.EXPECT().Profile(gomock.Any(), userID).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("GET", "/me", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("MissingUserID", func(t *testing.T) {
		// New app instance without middleware injection
		appMissing := fiber.New()
		appMissing.Get("/me", handler.ProfileHandler)

		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/me", nil)
		resp, _ := appMissing.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAthHandler_UpdatePasswordHandler(t *testing.T) {
	mockService, mockLogger, mockValidator, handler, app := setupHandlerTest(t)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	app.Post("/update-password", func(c fiber.Ctx) error {
		c.Locals("user_id", userID)
		return handler.UpdatePasswordHandler(c)
	})

	reqBody := user.UpdatePasswordRequest{
		OldPassword: "old",
		NewPassword: "new",
	}
	body, _ := json.Marshal(reqBody)

	t.Run("Success", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().UpdatePassword(gomock.Any(), userID, reqBody).Return(nil)

		req := httptest.NewRequest("POST", "/update-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().UpdatePassword(gomock.Any(), userID, reqBody).Return(common.ErrInvalidCredentials)

		req := httptest.NewRequest("POST", "/update-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("MissingUserID", func(t *testing.T) {
		appMissing := fiber.New()
		appMissing.Post("/update-password", handler.UpdatePasswordHandler)

		mockLogger.EXPECT().Errorf(gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/update-password", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := appMissing.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAthHandler_LogoutHandler(t *testing.T) {
	_, _, _, handler, app := setupHandlerTest(t)
	app.Post("/logout", handler.LogoutHandler)

	t.Run("Success", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/logout", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Optional: Verify cookies are cleared (MaxAge < 0)
		// Fiber Test response headers Set-Cookie
	})
}

func TestAthHandler_RefreshHandler(t *testing.T) {
	mockService, mockLogger, _, handler, app := setupHandlerTest(t)
	app.Post("/refresh", handler.RefreshHandler)

	t.Run("Success", func(t *testing.T) {
		mockService.EXPECT().RefreshToken(gomock.Any(), "valid_refresh").Return(&user.LoginResponse{
			Token:        "new_access",
			RefreshToken: "new_refresh",
			ExpiredAt:    time.Now().Add(time.Hour),
		}, nil)

		req := httptest.NewRequest("POST", "/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "valid_refresh"})
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("MissingCookie", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/refresh", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("ServiceFailure", func(t *testing.T) {
		mockService.EXPECT().RefreshToken(gomock.Any(), "invalid_refresh").Return(nil, common.ErrUnauthorized)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/refresh", nil)
		req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "invalid_refresh"})
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})
}

func TestAthHandler_AddUserHandler(t *testing.T) {
	mockService, mockLogger, mockValidator, handler, app := setupHandlerTest(t)

	app.Post("/users", func(c fiber.Ctx) error {
		c.Locals("role", user_repo.UserRoleAdmin)
		return handler.AddUserHandler(c)
	})

	reqBody := user.RegisterRequest{
		Username: "newurer",
		Email:    "new@example.com",
		Password: "password",
		Role:     user_repo.UserRoleCashier,
	}
	body, _ := json.Marshal(reqBody)

	t.Run("Success", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)

		respDto := &user.ProfileResponse{
			ID:       uuid.New(),
			Username: reqBody.Username,
			Email:    reqBody.Email,
			Role:     reqBody.Role,
		}
		mockService.EXPECT().Register(gomock.Any(), reqBody).Return(respDto, nil)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("UnauthorizedRole", func(t *testing.T) {
		appForbidden := fiber.New()
		appForbidden.Post("/users", func(c fiber.Ctx) error {
			c.Locals("role", user_repo.UserRoleCashier)
			return handler.AddUserHandler(c)
		})

		// Attempting to add Admin/Cashier as Cashier -> should fail if adding equal or higher
		// The mock handler logic: if(RoleLevel[req.Role] >= RoleLevel[currentRole]) check
		// Let's assume middleware RoleLevel is set up (but here we rely on the handler's check middleware.RoleLevel map)
		// We need to ensure middleware.RoleLevel is imported or accessible.
		// Actually handler uses `middleware.RoleLevel`. We need to make sure that map is initialized.
		// Usually it's a global var in package middleware.

		reqBodyAdmin := reqBody
		reqBodyAdmin.Role = user_repo.UserRoleAdmin // Higher than Cashier

		reqBodyBytes, _ := json.Marshal(reqBodyAdmin)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader(reqBodyBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := appForbidden.Test(req)

		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("UserExists", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().Register(gomock.Any(), reqBody).Return(nil, common.ErrUserExists)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().Register(gomock.Any(), reqBody).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestAthHandler_UpdateAvatarHandler(t *testing.T) {
	mockService, mockLogger, _, handler, app := setupHandlerTest(t)
	userID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	app.Put("/me/avatar", func(c fiber.Ctx) error {
		c.Locals("user_id", userID)
		return handler.UpdateAvatarHandler(c)
	})

	t.Run("Success", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("avatar", "test.jpg")
		part.Write([]byte("image data"))
		writer.Close()

		mockService.EXPECT().UploadAvatar(gomock.Any(), userID, gomock.Any()).Return(&user.ProfileResponse{
			ID:     userID,
			Avatar: func() *string { s := "url"; return &s }(),
		}, nil)

		req := httptest.NewRequest("PUT", "/me/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("NoUserID", func(t *testing.T) {
		appMissing := fiber.New()
		appMissing.Put("/me/avatar", handler.UpdateAvatarHandler)

		mockLogger.EXPECT().Errorf(gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/me/avatar", nil)
		resp, _ := appMissing.Test(req)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("FileTooLarge", func(t *testing.T) {
		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("avatar", "test.jpg")
		part.Write([]byte("image data"))
		writer.Close()

		mockService.EXPECT().UploadAvatar(gomock.Any(), userID, gomock.Any()).Return(nil, common.ErrFileTooLarge)

		req := httptest.NewRequest("PUT", "/me/avatar", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
