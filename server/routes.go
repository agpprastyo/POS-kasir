package server

import (
	"POS-kasir/config"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, logger *logger.Logger, db *database.Postgres, cfg *config.AppConfig) {
	app.Get("/healthz", HealthHandler(&App{
		DB:     db,
		Logger: logger,
		Config: cfg,
	}))
}
