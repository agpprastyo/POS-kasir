package auth

import (
	"POS-kasir/internal/repository"
	"POS-kasir/mocks"
	mocksauth "POS-kasir/mocks/auth"
	"POS-kasir/pkg/middleware"
	"bytes"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"
)

type errorFile struct {
	readErr error
}

func (f *errorFile) Read(p []byte) (int, error)                   { return 0, f.readErr }
func (f *errorFile) Close() error                                 { return nil }
func (f *errorFile) Seek(offset int64, whence int) (int64, error) { return 0, nil }
func (f *errorFile) ReadAt(p []byte, off int64) (int, error)      { return 0, f.readErr }

type errorFileHeader struct {
	openErr error
	file    multipart.File
}

func (fh *errorFileHeader) Open() (multipart.File, error) {
	if fh.openErr != nil {
		return nil, fh.openErr
	}
	return fh.file, nil
}

func TestUpdatePasswordHandler(t *testing.T) {
	tests := []struct {
		name       string
		handler    *auth.AthHandler
		setUserID  bool
		body       []byte
		wantStatus int
	}{
		{
			name: "Success",
			handler: &auth.AthHandler{
				Service:   MockIAuthService{},
				Log:       &mocks.MockILogger{},
				Validator: &mocks.MockValidator{},
			},
			setUserID:  true,
			body:       []byte(`{"old_password":"old","new_password":"new"}`),
			wantStatus: fiber.StatusOK,
		},
		{
			name: "BodyParseError",
			handler: &auth.AthHandler{
				Service:   mocksauth.MockIAuthService{},
				Log:       &mocks.MockILogger{},
				Validator: &mocks.MockValidator{},
			},
			setUserID:  true,
			body:       []byte(`{invalid json}`),
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "ValidationError",
			handler: &auth.AthHandler{
				Service:   mocksauth.MockIAuthService{},
				Log:       &mocks.MockILogger{},
				Validator: &mockValidator{validateErr: errors.New("validation failed")},
			},
			setUserID:  true,
			body:       []byte(`{"old_password":"old","new_password":"new"}`),
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "UserIDMissing",
			handler: &auth.AthHandler{
				service:   mocksauth.MockIAuthService{},
				log:       &mocks.MockILogger{},
				validator: &mocks.MockValidator{},
			},
			setUserID:  false,
			body:       []byte(`{"old_password":"old","new_password":"new"}`),
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name: "ServiceError",
			handler: &auth.AthHandler{
				service:   &mockAuthService{uploadErr: errors.New("update failed")},
				log:       &mocks.MockILogger{},
				validator: &mocks.MockValidator{},
			},
			setUserID:  true,
			body:       []byte(`{"old_password":"old","new_password":"new"}`),
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Post("/update-password", func(c *fiber.Ctx) error {
				if tt.setUserID {
					c.Locals("user_id", uuid.New())
				}
				return tt.handler.UpdatePasswordHandler(c)
			})

			req := httptest.NewRequest("POST", "/update-password", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestLoginHandler(t *testing.T) {
	tests := []struct {
		name       string
		service    *mockAuthService
		validator  *mockValidator
		body       []byte
		wantStatus int
	}{
		{
			name:       "Success",
			service:    &mockAuthService{uploadResult: &LoginResponse{Token: "token123", ExpiredAt: time.Now().Add(time.Hour), Profile: ProfileResponse{}}},
			validator:  &mocks.MockValidator{},
			body:       []byte(`{"username":"user","password":"pass"}`),
			wantStatus: fiber.StatusOK,
		},
		{
			name:       "BodyParseError",
			service:    mocksauth.MockIAuthService{},
			validator:  &mocks.MockValidator{},
			body:       []byte(`{invalid json}`),
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:       "ValidationError",
			service:    mocksauth.MockIAuthService{},
			validator:  &mockValidator{validateErr: errors.New("validation failed")},
			body:       []byte(`{"username":"user","password":"pass"}`),
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:       "ServiceError",
			service:    &mockAuthService{uploadErr: errors.New("login failed")},
			validator:  &mocks.MockValidator{},
			body:       []byte(`{"username":"user","password":"pass"}`),
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &AthHandler{
				service:   tt.service,
				log:       &mocks.MockILogger{},
				validator: tt.validator,
			}
			app.Post("/login", handler.LoginHandler)

			req := httptest.NewRequest("POST", "/login", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestProfileHandler(t *testing.T) {
	tests := []struct {
		name       string
		setUserID  bool
		service    *mockAuthService
		wantStatus int
	}{
		{
			name:       "Success",
			setUserID:  true,
			service:    &mockAuthService{uploadResult: new(string)},
			wantStatus: fiber.StatusOK,
		},
		{
			name:       "UserIDMissing",
			setUserID:  false,
			service:    &mockAuthService{uploadResult: new(string)},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:       "ServiceError",
			setUserID:  true,
			service:    &mockAuthService{uploadErr: errors.New("profile error")},
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &AthHandler{
				service:   tt.service,
				log:       &mocks.MockILogger{},
				validator: &mocks.MockValidator{},
			}
			app.Get("/profile", func(c *fiber.Ctx) error {
				if tt.setUserID {
					c.Locals("user_id", uuid.New())
				}
				return handler.ProfileHandler(c)
			})

			req := httptest.NewRequest("GET", "/profile", nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestLogoutHandler(t *testing.T) {
	tests := []struct {
		name       string
		wantStatus int
	}{
		{
			name:       "Success",
			wantStatus: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &AthHandler{
				log:       &mocks.MockILogger{},
				validator: &mocks.MockValidator{},
			}
			app.Post("/logout", handler.LogoutHandler)

			req := httptest.NewRequest("POST", "/logout", nil)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}

			// Check that the cookie is cleared
			cookies := resp.Cookies()
			found := false
			for _, cookie := range cookies {
				if cookie.Name == "access_token" && cookie.Value == "" {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("access_token cookie not cleared")
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       []byte
		validator  *mockValidator
		service    *mockAuthService
		wantStatus int
	}{
		{
			name:       "Success",
			body:       []byte(`{"username":"user","password":"pass"}`),
			validator:  &mocks.MockValidator{},
			service:    mocksauth.MockIAuthService{},
			wantStatus: fiber.StatusOK,
		},
		{
			name:       "BodyParseError",
			body:       []byte(`{invalid json}`),
			validator:  &mocks.MockValidator{},
			service:    mocksauth.MockIAuthService{},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:       "ValidationError",
			body:       []byte(`{"username":"user","password":"pass"}`),
			validator:  &mockValidator{validateErr: errors.New("validation failed")},
			service:    mocksauth.MockIAuthService{},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:       "ServiceError",
			body:       []byte(`{"username":"user","password":"pass"}`),
			validator:  &mocks.MockValidator{},
			service:    &mockAuthService{uploadErr: errors.New("register failed")},
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &AthHandler{
				service:   tt.service,
				log:       &mocks.MockILogger{},
				validator: tt.validator,
			}
			app.Post("/register", handler.RegisterHandler)

			req := httptest.NewRequest("POST", "/register", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestAddUserHandler(t *testing.T) {
	tests := []struct {
		name        string
		body        []byte
		currentRole interface{}
		reqRole     string
		validator   *mockValidator
		service     *mockAuthService
		wantStatus  int
	}{
		{
			name:        "Success",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: repository.UserRole("admin"),
			reqRole:     "cashier",
			validator:   &mocks.MockValidator{},
			service:     mocksauth.MockIAuthService{},
			wantStatus:  fiber.StatusOK,
		},
		{
			name:        "BodyParseError",
			body:        []byte(`{invalid json}`),
			currentRole: repository.UserRole("admin"),
			reqRole:     "cashier",
			validator:   &mocks.MockValidator{},
			service:     mocksauth.MockIAuthService{},
			wantStatus:  fiber.StatusBadRequest,
		},
		{
			name:        "InvalidCurrentUserRole",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: "not_a_role",
			reqRole:     "cashier",
			validator:   &mocks.MockValidator{},
			service:     mocksauth.MockIAuthService{},
			wantStatus:  fiber.StatusForbidden,
		},
		{
			name:        "AssignEqualOrHigherRole",
			body:        []byte(`{"username":"user","password":"pass","role":"admin"}`),
			currentRole: repository.UserRole("admin"),
			reqRole:     "admin",
			validator:   &mocks.MockValidator{},
			service:     mocksauth.MockIAuthService{},
			wantStatus:  fiber.StatusForbidden,
		},
		{
			name:        "ValidationError",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: repository.UserRole("admin"),
			reqRole:     "cashier",
			validator:   &mockValidator{validateErr: errors.New("validation failed")},
			service:     mocksauth.MockIAuthService{},
			wantStatus:  fiber.StatusBadRequest,
		},
		{
			name:        "ServiceError",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: repository.UserRole("admin"),
			reqRole:     "cashier",
			validator:   &mocks.MockValidator{},
			service:     &mockAuthService{uploadErr: errors.New("register failed")},
			wantStatus:  fiber.StatusBadRequest,
		},
	}

	// Setup role levels for test
	middleware.RoleLevel = map[repository.UserRole]int{
		repository.UserRole("admin"):   2,
		repository.UserRole("cashier"): 1,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &AthHandler{
				service:   tt.service,
				log:       &mocks.MockILogger{},
				validator: tt.validator,
			}
			app.Post("/add-user", func(c *fiber.Ctx) error {
				c.Locals("role", tt.currentRole)
				return handler.AddUserHandler(c)
			})

			req := httptest.NewRequest("POST", "/add-user", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestUpdateAvatarHandler(t *testing.T) {
	tests := []struct {
		name       string
		setUserID  bool
		createForm func() (*bytes.Buffer, string)
		service    *mockAuthService
		wantStatus int
	}{
		{
			name:      "Success",
			setUserID: true,
			createForm: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
				part.Write([]byte("avatar data"))
				writer.Close()
				return body, writer.FormDataContentType()
			},
			service:    &mockAuthService{uploadResult: "avatar_url"},
			wantStatus: fiber.StatusOK,
		},
		{
			name:      "UserIDMissing",
			setUserID: false,
			createForm: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
				part.Write([]byte("avatar data"))
				writer.Close()
				return body, writer.FormDataContentType()
			},
			service:    mocksauth.MockIAuthService{},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:      "NoAvatarFile",
			setUserID: true,
			createForm: func() (*bytes.Buffer, string) {
				// Create a multipart form but without the "avatar" file part
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close() // Close it immediately
				return body, writer.FormDataContentType()
			},
			service:    mocksauth.MockIAuthService{},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			// Note: Testing for file open/read errors is complex with httptest
			// as it requires a malformed request at a very low level.
			// The handler's error handling for `c.FormFile` already covers this.
			// A simple "NoAvatarFile" test is often sufficient to check the error path.
			name:      "ServiceError",
			setUserID: true,
			createForm: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				part, _ := writer.CreateFormFile("avatar", "avatar.jpg")
				part.Write([]byte("avatar data"))
				writer.Close()
				return body, writer.FormDataContentType()
			},
			service:    &mockAuthService{uploadErr: errors.New("upload failed")},
			wantStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			handler := &AthHandler{
				service:   tt.service,
				log:       &mocks.MockILogger{},
				validator: &mocks.MockValidator{},
			}
			app.Post("/avatar", func(c *fiber.Ctx) error {
				if tt.setUserID {
					c.Locals("user_id", uuid.New())
				}
				return handler.UpdateAvatarHandler(c)
			})

			body, contentType := tt.createForm()

			req := httptest.NewRequest("POST", "/avatar", body)
			req.Header.Set("Content-Type", contentType)
			resp, _ := app.Test(req)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}
