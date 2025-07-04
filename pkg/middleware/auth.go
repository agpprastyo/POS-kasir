package middleware

import (
	"POS-kasir/internal/common"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(tokenManager utils.Manager, log *logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("access_token")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "unauthorized",
			})
		}
		claims, err := tokenManager.VerifyToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "unauthorized",
			})
		}
		c.Locals("user", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("email", claims.Email)
		c.Locals("user_id", claims.UserID)

		log.Infof("current user is %v", claims.Username)
		log.Infof("current role is %v", claims.Role)
		log.Infof("current email is %v", claims.Email)
		log.Infof("current user ID is %v", claims.UserID)

		return c.Next()
	}
}
