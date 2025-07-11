package server

import (
	"POS-kasir/config"
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/categories"
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
	activityLog := activitylog.NewService(repo, logger)

	authRepo := auth.NewAuthRepo(logger, minio)
	authService := auth.NewAuthService(repo, logger, jwt, authRepo, activityLog)
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

	userService := user.NewUsrService(repo, logger, activityLog, authRepo)
	userHandler := user.NewUsrHandler(userService, logger, val)

	api.Get("/users", authMiddleware, userHandler.GetAllUsersHandler)
	api.Post("/users", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin), userHandler.CreateUserHandler)
	api.Get("/users/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), userHandler.GetUserByIDHandler)
	api.Put("/users/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), userHandler.UpdateUserHandler)
	api.Post("/users/:id/toggle-status", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin), userHandler.ToggleUserStatusHandler)

	categoryService := categories.NewCtgService(repo, logger, activityLog)
	categoryHandler := categories.NewCtgHandler(categoryService, logger)

	api.Get("/categories", authMiddleware, categoryHandler.GetAllCategoriesHandler)
	api.Post("/categories", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), categoryHandler.CreateCategoryHandler)
	api.Get("/categories/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), categoryHandler.GetCategoryByIDHandler)
	api.Put("/categories/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), categoryHandler.UpdateCategoryHandler)
	api.Delete("/categories/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), categoryHandler.DeleteCategoryHandler)
}
