package server

import (
	"POS-kasir/config"
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/categories"
	"POS-kasir/internal/products"
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

	store := repository.NewStore(db.DB, logger)

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

	// Initialize product service and handler
	prdRepo := products.NewPrdRepo(minio, logger)
	prdService := products.NewPrdService(store, logger, prdRepo, activityLog)
	prdHandler := products.NewPrdHandler(prdService, logger, val)

	api.Post("/products", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), prdHandler.CreateProductHandler)

	api.Post("/products/:id/image", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), prdHandler.UploadProductImageHandler)

	api.Get("/products", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), prdHandler.ListProductsHandler)

	api.Get("/products/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), prdHandler.GetProductHandler)

	//api.Patch("/products/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), prdHandler.UpdateProductHandler)
	//
	//api.Delete("/products/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin), prdHandler.DeleteProductHandler)
	//
	//api.Post("/products/:product_id/options", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), prdHandler.CreateProductOptionHandler)
	//
	//api.Post("/products/:product_id/options/:option_id/image", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), prdHandler.UploadProductOptionImageHandler)
	//
	//api.Patch("/products/:product_id/options/:option_id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), prdHandler.UpdateProductOptionHandler)
	//
	//api.Delete("/products/:product_id/options/:option_id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), prdHandler.DeleteProductOptionHandler)

}
