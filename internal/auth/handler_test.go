package auth

import (
	"POS-kasir/internal/repository"
	"POS-kasir/mocks"
	"POS-kasir/pkg/middleware"
	"bytes"
	"errors"
	"go.uber.org/mock/gomock"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"mime/multipart"
	"net/http/httptest"
	"testing"
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
		name      string
		setUserID bool
		body      []byte
		// Updated to accept the mocks for setup
		setupMocks func(mockValidator *mocks.MockValidator, mockService *MockIAuthService)
		wantStatus int
	}{
		{
			name:      "Success",
			setUserID: true,
			body:      []byte(`{"old_password":"old","new_password":"new"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockService.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:      "BodyParseError",
			setUserID: true,
			body:      []byte(`{invalid json}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				// No calls expected, so the function is empty.
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:      "ValidationError",
			setUserID: true,
			body:      []byte(`{"old_password":"old","new_password":"new"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:      "UserIDMissing",
			setUserID: false,
			body:      []byte(`{"old_password":"old","new_password":"new"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				// No calls expected. The handler panics/errors before reaching validation or service.
			},
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:      "ServiceError",
			setUserID: true,
			body:      []byte(`{"old_password":"old","new_password":"new"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				mockService.EXPECT().UpdatePassword(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("service error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// **CORRECTION**: Mocks are now created inside the loop for isolation.
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockIAuthService(ctrl)
			mockValidator := mocks.NewMockValidator(ctrl)
			mockLogger := mocks.NewMockILogger(ctrl)
			// Allow both Info and Errorf to be called any number of times
			mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

			// Run the setup function for this specific test case.
			if tt.setupMocks != nil {
				tt.setupMocks(mockValidator, mockService)
			}

			app := fiber.New()
			handler := &AthHandler{
				Service:   mockService,
				Log:       mockLogger,
				Validator: mockValidator,
			}

			app.Post("/update-password", func(c *fiber.Ctx) error {
				if tt.setUserID {
					c.Locals("user_id", uuid.New())
				}
				return handler.UpdatePasswordHandler(c)
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
		body       []byte
		setupMocks func(mockValidator *mocks.MockValidator, mockService *MockIAuthService)
		wantStatus int
	}{
		{
			name: "Success",
			body: []byte(`{"username":"user","password":"pass"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				// CORRECTED: Return a nil pointer for the struct and a nil error.
				mockService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(&LoginResponse{}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "BodyParseError",
			body: []byte(`{invalid json}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				// No calls expected
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "ValidationError",
			body: []byte(`{"username":"user","password":"pass"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "ServiceError",
			body: []byte(`{"username":"user","password":"pass"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				// CORRECTED: Return a nil pointer for the struct and the error.
				mockService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockIAuthService(ctrl)
			mockValidator := mocks.NewMockValidator(ctrl)
			mockLogger := mocks.NewMockILogger(ctrl)
			mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

			if tt.setupMocks != nil {
				tt.setupMocks(mockValidator, mockService)
			}

			app := fiber.New()
			handler := &AthHandler{
				Service:   mockService,
				Log:       mockLogger,
				Validator: mockValidator,
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
		setupMocks func(mockService *MockIAuthService)
		wantStatus int
	}{
		{
			name:      "Success",
			setUserID: true,
			setupMocks: func(mockService *MockIAuthService) {
				// Assuming Profile returns a response struct and nil error
				mockService.EXPECT().Profile(gomock.Any(), gomock.Any()).Return(&ProfileResponse{}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:       "UserIDMissing",
			setUserID:  false,
			setupMocks: nil, // No service call expected
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:      "ServiceError",
			setUserID: true,
			setupMocks: func(mockService *MockIAuthService) {
				mockService.EXPECT().Profile(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockIAuthService(ctrl)
			mockLogger := mocks.NewMockILogger(ctrl)
			mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

			if tt.setupMocks != nil {
				tt.setupMocks(mockService)
			}

			app := fiber.New()
			handler := &AthHandler{
				Service:   mockService,
				Log:       mockLogger,
				Validator: nil, // Validator is not used in ProfileHandler
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
				Log:       &mocks.MockILogger{},
				Validator: &mocks.MockValidator{},
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
		setupMocks func(mockValidator *mocks.MockValidator, mockService *MockIAuthService)
		wantStatus int
	}{
		{
			name: "Success",
			body: []byte(`{"username":"user","password":"pass"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				// CORRECTED: Return a *ProfileResponse as indicated by the error message.
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&ProfileResponse{}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name: "BodyParseError",
			body: []byte(`{invalid json}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				// No calls to mocks are expected
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "ValidationError",
			body: []byte(`{"username":"user","password":"pass"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name: "ServiceError",
			body: []byte(`{"username":"user","password":"pass"}`),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				// CORRECTED: Return a nil *ProfileResponse and the error.
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockIAuthService(ctrl)
			mockValidator := mocks.NewMockValidator(ctrl)
			mockLogger := mocks.NewMockILogger(ctrl)
			mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

			if tt.setupMocks != nil {
				tt.setupMocks(mockValidator, mockService)
			}

			app := fiber.New()
			handler := &AthHandler{
				Service:   mockService,
				Log:       mockLogger,
				Validator: mockValidator,
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
	middleware.RoleLevel = map[repository.UserRole]int{
		repository.UserRole("admin"):   2,
		repository.UserRole("cashier"): 1,
	}

	tests := []struct {
		name        string
		body        []byte
		currentRole interface{}
		setupMocks  func(mockValidator *mocks.MockValidator, mockService *MockIAuthService)
		wantStatus  int
	}{
		{
			name:        "Success",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: repository.UserRole("admin"),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				// CORRECTED: Assuming AddUser also returns a *ProfileResponse.
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(&ProfileResponse{}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:        "BodyParseError",
			body:        []byte(`{invalid json}`),
			currentRole: repository.UserRole("admin"),
			setupMocks:  nil,
			wantStatus:  fiber.StatusBadRequest,
		},
		{
			name:        "InvalidCurrentUserRole",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: "not_a_role",
			setupMocks:  nil,
			wantStatus:  fiber.StatusForbidden,
		},
		{
			name:        "AssignEqualOrHigherRole",
			body:        []byte(`{"username":"user","password":"pass","role":"admin"}`),
			currentRole: repository.UserRole("admin"),
			setupMocks:  nil,
			wantStatus:  fiber.StatusForbidden,
		},
		{
			name:        "ValidationError",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: repository.UserRole("admin"),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(errors.New("validation error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
		{
			name:        "ServiceError",
			body:        []byte(`{"username":"user","password":"pass","role":"cashier"}`),
			currentRole: repository.UserRole("admin"),
			setupMocks: func(mockValidator *mocks.MockValidator, mockService *MockIAuthService) {
				mockValidator.EXPECT().Validate(gomock.Any()).Return(nil)
				// CORRECTED: Return a nil *ProfileResponse and the error.
				mockService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			wantStatus: fiber.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockIAuthService(ctrl)
			mockValidator := mocks.NewMockValidator(ctrl)
			mockLogger := mocks.NewMockILogger(ctrl)
			mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

			if tt.setupMocks != nil {
				tt.setupMocks(mockValidator, mockService)
			}

			app := fiber.New()
			handler := &AthHandler{
				Service:   mockService,
				Log:       mockLogger,
				Validator: mockValidator,
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
		setupMocks func(mockService *MockIAuthService)
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
			setupMocks: func(mockService *MockIAuthService) {
				// CORRECTED: Return a *ProfileResponse and a nil error.
				mockService.EXPECT().UploadAvatar(gomock.Any(), gomock.Any(), gomock.Any()).Return(&ProfileResponse{}, nil)
			},
			wantStatus: fiber.StatusOK,
		},
		{
			name:      "UserIDMissing",
			setUserID: false,
			createForm: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()
				return body, writer.FormDataContentType()
			},
			setupMocks: nil, // No service call expected
			wantStatus: fiber.StatusInternalServerError,
		},
		{
			name:      "NoAvatarFile",
			setUserID: true,
			createForm: func() (*bytes.Buffer, string) {
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				writer.Close()
				return body, writer.FormDataContentType()
			},
			setupMocks: nil, // No service call expected
			wantStatus: fiber.StatusBadRequest,
		},
		{
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
			setupMocks: func(mockService *MockIAuthService) {
				// CORRECTED: Return a nil *ProfileResponse and the error.
				mockService.EXPECT().UploadAvatar(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("service error"))
			},
			wantStatus: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockService := NewMockIAuthService(ctrl)
			mockLogger := mocks.NewMockILogger(ctrl)
			mockLogger.EXPECT().Info(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any()).AnyTimes()
			mockLogger.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()

			if tt.setupMocks != nil {
				tt.setupMocks(mockService)
			}

			app := fiber.New()
			handler := &AthHandler{
				Service:   mockService,
				Log:       mockLogger,
				Validator: nil,
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
