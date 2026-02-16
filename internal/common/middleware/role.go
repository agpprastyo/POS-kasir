package middleware

import (
	"github.com/gofiber/fiber/v3"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "admin"
	UserRoleManager UserRole = "manager"
	UserRoleCashier UserRole = "cashier"
)

var RoleLevel = map[UserRole]int{
	UserRoleAdmin:   3,
	UserRoleManager: 2,
	UserRoleCashier: 1,
}

// RoleMiddleware checks if the user's role meets the minimum required role.
func RoleMiddleware(minRole UserRole) fiber.Handler {
	return func(c fiber.Ctx) error {
		roleVal := c.Locals("role")
		var userRole UserRole

		switch v := roleVal.(type) {
		case UserRole:
			userRole = v
		case string:
			userRole = UserRole(v)
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
