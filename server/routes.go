package server

import (
	"POS-kasir/config"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/repository"
	"POS-kasir/internal/user"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/middleware"
	"POS-kasir/pkg/minio"
	"POS-kasir/pkg/utils"
	"POS-kasir/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, logger *logger.Logger, db *database.Postgres, cfg *config.AppConfig, jwt utils.Manager, minio *minio.Minio) {
	hltHandler := HealthHandler(&App{
		DB:     db,
		Logger: logger,
		Config: cfg,
	})
	app.Get("/healthz", hltHandler)

	api := app.Group("/api/v1")

	repo := repository.New(db.DB)

	authRepo := auth.NewAuthRepo(logger, minio)
	authService := auth.NewAuthService(*repo, logger, jwt, authRepo)
	val := validator.NewValidator()
	authHandler := auth.NewAuthHandler(authService, logger, val)
	authMiddleware := middleware.AuthMiddleware(jwt, logger)

	api.Post("/auth/login", authHandler.LoginHandler)
	api.Post("/auth/register", authHandler.RegisterHandler)
	api.Get("/auth/me", authMiddleware, authHandler.ProfileHandler)
	api.Post("/auth/add", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), authHandler.AddUserHandler)
	api.Put("/auth/me/avatar", authMiddleware, authHandler.UpdateAvatarHandler)
	api.Put("/auth/me/password", authMiddleware, authHandler.UpdatePasswordHandler)
	api.Post("/auth/logout", authMiddleware, authHandler.LogoutHandler)

	userService := user.NewUsrService(*repo, logger)
	userHandler := user.NewUsrHandler(userService, logger, val)

	api.Get("/users", authMiddleware, userHandler.GetAllUsersHandler)
}
