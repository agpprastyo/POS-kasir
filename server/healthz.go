package server

import (
	"github.com/gofiber/fiber/v2"
)

func HealthHandler(app *App) fiber.Handler {
	return func(c *fiber.Ctx) error {

		if err := app.DB.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "fail",
				"error":  "PostgreSQL unavailable",
			})
		}

		exists, err := app.Minio.BucketExists(c.Context())
		if err != nil || !exists {
			if err != nil {
				app.Logger.Errorf("Error checking bucket existence: %v", err)
			}
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "fail",
				"error":  "Minio bucket unavailable",
			})
		}

		return c.JSON(fiber.Map{
			"status": "ok",
		})
	}
}
