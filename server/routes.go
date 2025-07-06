package server

import (
	"POS-kasir/config"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/middleware"
	"POS-kasir/pkg/utils"
	"POS-kasir/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, logger *logger.Logger, db *database.Postgres, cfg *config.AppConfig, jwt utils.Manager) {
	hltHandler := HealthHandler(&App{
		DB:     db,
		Logger: logger,
		Config: cfg,
	})
	app.Get("/healthz", hltHandler)

	api := app.Group("/api/v1")

	repo := repository.New(db.DB)
	authService := auth.NewAuthService(*repo, logger, jwt)
	val := validator.NewValidator()
	authHandler := auth.NewAuthHandler(authService, logger, val)
	authMiddleware := middleware.AuthMiddleware(jwt, logger)

	api.Post("/auth/login", authHandler.Loginhandler)
	api.Post("/auth/register", authHandler.RegisterHandler)
	api.Get("/auth/me", authMiddleware, authHandler.ProfileHandler)
	api.Post("/auth/add", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), authHandler.AddUserHandler)
	api.Put("/auth/avatar", authMiddleware, authHandler.UpdateAvatarHandler)
}
