package auth

import (
	"POS-kasir/config"
	"POS-kasir/internal/common"
	"POS-kasir/internal/common/middleware"
	"POS-kasir/internal/dto"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"

	"POS-kasir/pkg/validator"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type IAuthHandler interface {
	UpdatePasswordHandler(c fiber.Ctx) error
	LoginHandler(c fiber.Ctx) error
	LogoutHandler(c fiber.Ctx) error

	ProfileHandler(c fiber.Ctx) error
	AddUserHandler(c fiber.Ctx) error
	UpdateAvatarHandler(c fiber.Ctx) error
	RefreshHandler(c fiber.Ctx) error
}

type AthHandler struct {
	cfg       *config.AppConfig
	Service   IAuthService
	Log       logger.ILogger
	Validator validator.Validator
}

func NewAuthHandler(service IAuthService, log logger.ILogger, validator validator.Validator, cfg *config.AppConfig) IAuthHandler {
	return &AthHandler{
		Service:   service,
		Log:       log,
		Validator: validator,
		cfg:       cfg,
	}
}

// UpdatePasswordHandler handles password update requests.
// @Summary Update password
// @Description Update password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.UpdatePasswordRequest true "Update password request"
// @Success 200 {object} common.SuccessResponse{data=string} "Success"
// @Failure 400 {object} common.ErrorResponse "Bad Request"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Not Found"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /auth/me/update-password [post]
func (h *AthHandler) UpdatePasswordHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("UpdatePasswordHandler | Failed to get userID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "No user ID in context",
		})
	}

	var req dto.UpdatePasswordRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
			Error:   err.Error(),
		})
	}

	if done, err := common.ValidateAndRespond(c, h.Validator, h.Log, &req); done {
		return err
	}

	err := h.Service.UpdatePassword(ctx, userUUID, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "Invalid current password",
			})
		case errors.Is(err, common.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		default:
			h.Log.Errorf("Auth Handler | Failed to update password: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to update password",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Password updated successfully",
	})
}

// LoginHandler handles login requests.
// @Summary Login
// @Description Login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} common.SuccessResponse{data=dto.LoginResponse} "Success"
// @Failure 400 {object} common.ErrorResponse "Bad Request"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Not Found"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /auth/login [post]
func (h *AthHandler) LoginHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	var req dto.LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
			Error:   err.Error(),
		})
	}

	if done, err := common.ValidateAndRespond(c, h.Validator, h.Log, &req); done {
		return err
	}

	resp, err := h.Service.Login(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, common.ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "Invalid username or password",
			})
		default:
			h.Log.Errorf("LoginHandler | Failed to login: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to login",
			})
		}
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    resp.Token,
		Path:     "/",
		Domain:   h.cfg.Server.CookieDomain,
		Expires:  resp.ExpiredAt,
		MaxAge:   int(time.Until(resp.ExpiredAt).Seconds()),
		HTTPOnly: true,
		Secure:   h.cfg.Server.Env == "production",
		SameSite: func() string {
			if h.cfg.Server.Env == "production" {
				return fiber.CookieSameSiteNoneMode
			}
			return fiber.CookieSameSiteLaxMode
		}(),
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		Domain:   h.cfg.Server.CookieDomain,
		Expires:  time.Now().Add(h.cfg.JWT.RefreshTokenDuration),
		MaxAge:   int(h.cfg.JWT.RefreshTokenDuration.Seconds()),
		HTTPOnly: true,
		Secure:   h.cfg.Server.Env == "production",
		SameSite: func() string {
			if h.cfg.Server.Env == "production" {
				return fiber.CookieSameSiteNoneMode
			}
			return fiber.CookieSameSiteLaxMode
		}(),
	}

	if h.cfg.Server.WebFrontendCrossOrigin {
		cookie.SameSite = fiber.CookieSameSiteNoneMode
		cookie.Secure = true
		refreshCookie.SameSite = fiber.CookieSameSiteNoneMode
		refreshCookie.Secure = true
	}
	c.Cookie(cookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data: fiber.Map{
			"expired": resp.ExpiredAt,
			"profile": resp.Profile,
		},
	})
}

// LogoutHandler handles logout requests.
// @Summary Logout
// @Description Logout
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse{data=string} "Success"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Not Found"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /auth/logout [post]
func (h *AthHandler) LogoutHandler(c fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		Domain:   h.cfg.Server.CookieDomain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.cfg.Server.Env == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		Domain:   h.cfg.Server.CookieDomain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   h.cfg.Server.Env == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	})
	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Successfully logged out",
	})
}

// ProfileHandler handles requests to retrieve the profile of the currently authenticated user.
// @Summary Get profile
// @Description Get profile
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse{data=dto.ProfileResponse} "Success"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Not Found"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /auth/me [get]
func (h *AthHandler) ProfileHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("ProfileHandler | Failed to get userID from context, userId : %v", userUUID)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Error:   "No user id in context",
			Message: "Failed to get user ID",
		})
	}

	response, err := h.Service.Profile(ctx, userUUID)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "User not found",
			})
		case errors.Is(err, common.ErrAvatarNotFound):
			return c.Status(fiber.StatusNotFound).JSON(common.ErrorResponse{
				Message: "Avatar not found",
			})
		default:
			h.Log.Errorf("ProfileHandler | Failed to get user profile: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to get user profile",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Success",
		Data:    response,
	})

}

// AddUserHandler handles requests to add a new user to the system.
// @Summary Add user
// @Description Add user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @x-roles ["admin"]
// @Success 200 {object} common.SuccessResponse{data=dto.ProfileResponse} "Success"
// @Failure 400 {object} common.ErrorResponse "Bad Request"
// @Failure 409 {object} common.ErrorResponse "Conflict"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
func (h *AthHandler) AddUserHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	var req dto.RegisterRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to parse request body",
		})
	}

	currentRoleVal := c.Locals("role")
	currentRole, ok := currentRoleVal.(repository.UserRole)
	if !ok {
		return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
			Message: "Invalid current user role",
		})
	}

	if middleware.RoleLevel[req.Role] >= middleware.RoleLevel[currentRole] {
		return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
			Message: "Cannot assign a role equal to or higher than your own",
		})
	}

	if done, err := common.ValidateAndRespond(c, h.Validator, h.Log, &req); done {
		return err
	}

	resp, err := h.Service.Register(ctx, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrUserExists):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "User already exists",
			})
		case errors.Is(err, common.ErrUsernameExists):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Username already exists",
			})
		case errors.Is(err, common.ErrEmailExists):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Email already exists",
			})
		default:
			h.Log.Errorf("AddUserHandler | Failed to register user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to register user",
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "User added successfully",
		Data:    resp,
	})
}

// UpdateAvatarHandler handles requests to update the avatar of the currently authenticated user.
// @Summary Update avatar
// @Description Update avatar
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar file"
// @Success 200 {object} common.SuccessResponse{data=dto.ProfileResponse} "Success"
// @Failure 400 {object} common.ErrorResponse "Bad Request"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 404 {object} common.ErrorResponse "Not Found"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /auth/me/avatar [put]
func (h *AthHandler) UpdateAvatarHandler(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userUUID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		h.Log.Errorf("UpdateAvatarHandler | Failed to get userID from context")
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "No user ID in context",
		})
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		h.Log.Errorf("UpdateAvatarHandler | Failed to get avatar file: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
			Message: "Failed to upload avatar",
		})
	}

	fileData, err := file.Open()
	if err != nil {
		h.Log.Errorf("UpdateAvatarHandler | Failed to open avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to process avatar file",
		})
	}

	defer func(fileData multipart.File) {
		err := fileData.Close()
		if err != nil {
			h.Log.Errorf("UpdateAvatarHandler | Failed to close avatar file: %v", err)
		}
	}(fileData)

	data, err := io.ReadAll(fileData)
	if err != nil {
		h.Log.Errorf("UpdateAvatarHandler | Failed to read avatar file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
			Message: "Failed to read avatar file",
		})
	}

	response, err := h.Service.UploadAvatar(ctx, userUUID, data)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrFileTooLarge):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "File size exceeds the maximum limit",
			})
		case errors.Is(err, common.ErrFileTypeNotSupported):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "File type is not supported",
			})
		case errors.Is(err, common.ErrImageNotSquare):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Image must be square",
			})
		case errors.Is(err, common.ErrImageTooSmall):
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "Image dimensions are too small, must be at least 300x300 pixels",
			})
		case errors.Is(err, common.ErrImageProcessingFailed):
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Image processing failed",
			})
		case errors.Is(err, common.ErrUploadFailed):
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "File upload failed",
			})

		default:
			h.Log.Errorf("UpdateAvatarHandler | Failed to upload avatar: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Failed to upload avatar",
			})
		}

	}

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Avatar updated successfully",
		Data:    response,
	})
}

// RefreshHandler handles token refresh requests.
// @Summary Refresh token
// @Description Refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} common.SuccessResponse{data=dto.LoginResponse} "Success"
// @Failure 401 {object} common.ErrorResponse "Unauthorized"
// @Failure 500 {object} common.ErrorResponse "Internal Server Error"
// @Router /auth/refresh [post]
func (h *AthHandler) RefreshHandler(c fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
			Message: "Refresh token missing",
		})
	}

	resp, err := h.Service.RefreshToken(c.RequestCtx(), refreshToken)
	if err != nil {
		h.Log.Errorf("RefreshHandler | Failed to refresh token: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
			Message: "Invalid or expired session",
		})
	}

	cookie := &fiber.Cookie{
		Name:     "access_token",
		Value:    resp.Token,
		Path:     "/",
		Domain:   h.cfg.Server.CookieDomain,
		Expires:  resp.ExpiredAt,
		MaxAge:   int(time.Until(resp.ExpiredAt).Seconds()),
		HTTPOnly: true,
		Secure:   h.cfg.Server.Env == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	}

	refreshCookie := &fiber.Cookie{
		Name:     "refresh_token",
		Value:    resp.RefreshToken,
		Path:     "/",
		Domain:   h.cfg.Server.CookieDomain,
		Expires:  time.Now().Add(h.cfg.JWT.RefreshTokenDuration),
		MaxAge:   int(h.cfg.JWT.RefreshTokenDuration.Seconds()),
		HTTPOnly: true,
		Secure:   h.cfg.Server.Env == "production",
		SameSite: fiber.CookieSameSiteLaxMode,
	}

	if h.cfg.Server.WebFrontendCrossOrigin {
		cookie.SameSite = fiber.CookieSameSiteNoneMode
		cookie.Secure = true
		refreshCookie.SameSite = fiber.CookieSameSiteNoneMode
		refreshCookie.Secure = true
	}
	c.Cookie(cookie)
	c.Cookie(refreshCookie)

	return c.Status(fiber.StatusOK).JSON(common.SuccessResponse{
		Message: "Token refreshed successfully",
		Data: fiber.Map{
			"expired": resp.ExpiredAt,
			"profile": resp.Profile,
		},
	})
}

// fiber:context-methods migrated
