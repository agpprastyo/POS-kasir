package middleware

import (
	"POS-kasir/config"
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"time"

	"github.com/gofiber/fiber/v3"
)

func JWTAuthMiddleware(tokenManager utils.Manager, cfg *config.AppConfig, log logger.ILogger, require bool) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Cookies("access_token", "")
		if token == "" {
			if require {
				log.Warnf("unauthorized access attempt: no token provided")
				return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{Message: "unauthorized"})
			}
			return c.Next()
		}

		claims, err := tokenManager.VerifyToken(token)
		if err != nil {
			c.Cookie(&fiber.Cookie{
				Name:     "access_token",
				Value:    "",
				Path:     "/",
				Domain:   cfg.Server.CookieDomain,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HTTPOnly: true,
				Secure:   cfg.Server.Env == "production",
				SameSite: fiber.CookieSameSiteLaxMode,
			})
			log.Warnf("unauthorized access attempt: invalid token - %v", err)
			if require {
				return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{Message: "unauthorized"})
			}
			return c.Next()
		}

		c.Locals("user", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("email", claims.Email)
		c.Locals("user_id", claims.UserID)

		c.RequestCtx().SetUserValue(common.UserIDKey, claims.UserID)

		return c.Next()
	}
}
