package middleware

import (
	"POS-kasir/internal/repository"

	"github.com/gofiber/fiber/v2"
)

var RoleLevel = map[repository.UserRole]int{
	repository.UserRoleAdmin:   3,
	repository.UserRoleManager: 2,
	repository.UserRoleCashier: 1,
}

// RoleMiddleware checks if the user's role meets the minimum required role.
func RoleMiddleware(minRole repository.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals("role")
		var userRole repository.UserRole

		switch v := roleVal.(type) {
		case repository.UserRole:
			userRole = v
		case string:
			userRole = repository.UserRole(v)
		default:
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "invalid role"})
		}

		reqLevel, ok1 := RoleLevel[userRole]
		minLevel, ok2 := RoleLevel[minRole]
		if !ok1 || !ok2 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "invalid role"})
		}

		if reqLevel < minLevel {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient role"})
		}
		return c.Next()
	}
}
