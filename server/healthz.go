package server

import (
	"POS-kasir/pkg/database"
	"github.com/gofiber/fiber/v2"
)

// HealthHandler returns 200 if the app and dependencies are healthy.
func HealthHandler(app *App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := database.PingPostgresPool(app.DB); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "fail",
				"error":  "PostgreSQL unavailable",
			})
		}
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	}
}
