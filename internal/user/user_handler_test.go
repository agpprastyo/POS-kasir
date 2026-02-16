package user_test

import (
	"POS-kasir/internal/common"
	"POS-kasir/internal/user"
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
	user_repo "POS-kasir/internal/user/repository"
)

func setupUserHandlerTest(t *testing.T) (*mocks.MockIUsrService, *mocks.MockFieldLogger, *mocks.MockValidator, user.IUsrHandler, *fiber.App) {
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockIUsrService(ctrl)
	mockLogger := mocks.NewMockFieldLogger(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)

	handler := user.NewUsrHandler(mockService, mockLogger, mockValidator)
	app := fiber.New()

	return mockService, mockLogger, mockValidator, handler, app
}

func TestUsrHandler_DeleteUserHandler(t *testing.T) {
	mockService, mockLogger, _, handler, app := setupUserHandlerTest(t)
	app.Delete("/users/:id", handler.DeleteUserHandler)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().DeleteUser(gomock.Any(), userID).Return(nil)

		req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest("DELETE", "/users/invalid-id", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().DeleteUser(gomock.Any(), userID).Return(common.ErrNotFound)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().DeleteUser(gomock.Any(), userID).Return(errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("DELETE", "/users/"+userID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestUsrHandler_ToggleUserStatusHandler(t *testing.T) {
	mockService, mockLogger, _, handler, app := setupUserHandlerTest(t)
	app.Post("/users/:id/toggle-status", handler.ToggleUserStatusHandler)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().ToggleUserStatus(gomock.Any(), userID).Return(nil)

		req := httptest.NewRequest("POST", "/users/"+userID.String()+"/toggle-status", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest("POST", "/users/invalid-id/toggle-status", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().ToggleUserStatus(gomock.Any(), userID).Return(common.ErrNotFound)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/users/"+userID.String()+"/toggle-status", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().ToggleUserStatus(gomock.Any(), userID).Return(common.ErrInternal)
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/users/"+userID.String()+"/toggle-status", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestUsrHandler_GetUserByIDHandler(t *testing.T) {
	mockService, mockLogger, _, handler, app := setupUserHandlerTest(t)
	app.Get("/users/:id", handler.GetUserByIDHandler)

	t.Run("Success", func(t *testing.T) {
		userID := uuid.New()
		respUser := &user.ProfileResponse{ID: userID, Username: "test"}
		mockService.EXPECT().GetUserByID(gomock.Any(), userID).Return(respUser, nil)

		req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest("GET", "/users/invalid-id", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, common.ErrNotFound)

		req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		userID := uuid.New()
		mockService.EXPECT().GetUserByID(gomock.Any(), userID).Return(nil, errors.New("db error"))

		req := httptest.NewRequest("GET", "/users/"+userID.String(), nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestUsrHandler_CreateUserHandler(t *testing.T) {
	mockService, mockLogger, mockValidator, handler, app := setupUserHandlerTest(t)
	app.Post("/users", handler.CreateUserHandler)

	reqBody := user.CreateUserRequest{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password",
	}
	body, _ := json.Marshal(reqBody)

	t.Run("Success", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		respUser := &user.ProfileResponse{Username: reqBody.Username}
		mockService.EXPECT().CreateUser(gomock.Any(), reqBody).Return(respUser, nil)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest("POST", "/users", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationFailed", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UserExists", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().CreateUser(gomock.Any(), reqBody).Return(nil, common.ErrUserExists)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("UsernameExists", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().CreateUser(gomock.Any(), reqBody).Return(nil, common.ErrUsernameExists)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("EmailExists", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().CreateUser(gomock.Any(), reqBody).Return(nil, common.ErrEmailExists)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().CreateUser(gomock.Any(), reqBody).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
func TestUsrHandler_GetAllUsersHandler(t *testing.T) {
	mockService, mockLogger, mockValidator, handler, app := setupUserHandlerTest(t)
	app.Get("/users", handler.GetAllUsersHandler)

	t.Run("Success", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().GetAllUsers(gomock.Any(), gomock.Any()).Return(&user.UsersResponse{}, nil)

		req := httptest.NewRequest("GET", "/users", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("InvalidQueryParams", func(t *testing.T) {
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		// Send param that cannot be bound, strictly speaking difficult with fiber params
		// but providing weird boolean might work or just invalid query string
		req := httptest.NewRequest("GET", "/users?page=invalid", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ValidationFailed", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation failed"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/users", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().GetAllUsers(gomock.Any(), gomock.Any()).Return(nil, common.ErrNotFound)
		mockLogger.EXPECT().Warnf(gomock.Any()).Times(1)

		req := httptest.NewRequest("GET", "/users", nil)
		resp, _ := app.Test(req)
		// Logic in handler says if not found, return 200 with empty list or error logic?
		// Handler code: Warnf("No users found") -> return 200 OK with message "Users retrieved successfully" and data response
		// Let's check handler code again
		// case errors.Is(err, common.ErrNotFound): h.log.Warnf("GetAllUsersHandler | No users found"); return c.Status(fiber.StatusOK)...
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestUsrHandler_UpdateUserHandler(t *testing.T) {
	mockService, mockLogger, mockValidator, handler, _ := setupUserHandlerTest(t)
	// We need to customize app to inject Locals
	app := fiber.New()

	// Helper to set role in context
	setupRoute := func(role string) {
		app.Put("/users/:id", func(c fiber.Ctx) error {
			if role != "" {
				c.Locals("role", role)
			}
			return handler.UpdateUserHandler(c)
		})
	}

	userID := uuid.New()
	reqBody := user.UpdateUserRequest{
		Username: func() *string { s := "newuser"; return &s }(),
	}
	body, _ := json.Marshal(reqBody)

	t.Run("Success", func(t *testing.T) {
		setupRoute("admin") // Authorized
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(&user.ProfileResponse{Username: "newuser"}, nil)
		mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("UnauthorizedRoleChange", func(t *testing.T) {
		setupRoute("cashier") // Not admin
		role := user_repo.UserRoleManager
		unauthBody := user.UpdateUserRequest{Role: &role}
		b, _ := json.Marshal(unauthBody)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockLogger.EXPECT().Warnf(gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		setupRoute("")
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)
		req := httptest.NewRequest("PUT", "/users/invalid-id", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound", func(t *testing.T) {
		setupRoute("")
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil, common.ErrNotFound)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("UsernameExists", func(t *testing.T) {
		setupRoute("")
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil, common.ErrUsernameExists)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("EmailExists", func(t *testing.T) {
		setupRoute("")
		email := "exists@example.com"
		emailBody := user.UpdateUserRequest{Email: &email}
		b, _ := json.Marshal(emailBody)

		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil, common.ErrEmailExists)
		mockLogger.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
	})

	t.Run("InternalError", func(t *testing.T) {
		setupRoute("")
		mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
		mockService.EXPECT().UpdateUser(gomock.Any(), userID, gomock.Any()).Return(nil, errors.New("db error"))
		mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).Times(1)

		req := httptest.NewRequest("PUT", "/users/"+userID.String(), bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}
