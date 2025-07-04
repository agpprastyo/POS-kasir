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
		userRole, ok := roleVal.(repository.UserRole)
		if !ok {
			// If stored as string, convert to UserRole
			roleStr, ok := roleVal.(string)
			if !ok {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "invalid role"})
			}
			userRole = repository.UserRole(roleStr)
		}
		if RoleLevel[userRole] < RoleLevel[minRole] {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "insufficient role"})
		}
		return c.Next()
	}
}
