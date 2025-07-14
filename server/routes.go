package server

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/middleware"
)

func SetupRoutes(app *App, container *AppContainer) {
	hltHandler := HealthHandler(app)
	app.FiberApp.Get("/healthz", hltHandler)

	api := app.FiberApp.Group("/api/v1")

	authMiddleware := middleware.AuthMiddleware(app.JWT, app.Logger)

	api.Post("/auth/login", container.AuthHandler.LoginHandler)
	api.Post("/auth/register", container.AuthHandler.RegisterHandler)
	api.Get("/auth/me", authMiddleware, container.AuthHandler.ProfileHandler)
	api.Post("/auth/add", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.AuthHandler.AddUserHandler)
	api.Put("/auth/me/avatar", authMiddleware, container.AuthHandler.UpdateAvatarHandler)
	api.Put("/auth/me/password", authMiddleware, container.AuthHandler.UpdatePasswordHandler)
	api.Post("/auth/logout", authMiddleware, container.AuthHandler.LogoutHandler)

	api.Get("/users", authMiddleware, container.UserHandler.GetAllUsersHandler)
	api.Post("/users", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin), container.UserHandler.CreateUserHandler)
	api.Get("/users/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.UserHandler.GetUserByIDHandler)
	api.Put("/users/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.UserHandler.UpdateUserHandler)
	api.Post("/users/:id/toggle-status", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin), container.UserHandler.ToggleUserStatusHandler)
	api.Delete("/users/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin), container.UserHandler.DeleteUserHandler)

	api.Get("/categories", authMiddleware, container.CategoryHandler.GetAllCategoriesHandler)
	api.Post("/categories", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.CategoryHandler.CreateCategoryHandler)
	api.Get("/categories/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.CategoryHandler.GetCategoryByIDHandler)
	api.Put("/categories/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.CategoryHandler.UpdateCategoryHandler)
	api.Delete("/categories/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.CategoryHandler.DeleteCategoryHandler)

	api.Post("/products", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.ProductHandler.CreateProductHandler)
	api.Post("/products/:id/image", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.ProductHandler.UploadProductImageHandler)
	api.Get("/products", authMiddleware, container.ProductHandler.ListProductsHandler)
	api.Get("/products/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.ProductHandler.GetProductHandler)
	api.Patch("/products/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.ProductHandler.UpdateProductHandler)
	api.Delete("/products/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin), container.ProductHandler.DeleteProductHandler)

	api.Post("/products/:product_id/options", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.ProductHandler.CreateProductOptionHandler)
	api.Post("/products/:product_id/options/:option_id/image", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.ProductHandler.UploadProductOptionImageHandler)
	api.Patch("/products/:product_id/options/:option_id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.ProductHandler.UpdateProductOptionHandler)
	api.Delete("/products/:product_id/options/:option_id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager), container.ProductHandler.DeleteProductOptionHandler)

	api.Post("/orders", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.CreateOrderHandler)
	api.Get("/orders/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.GetOrderHandler)
	api.Post("/orders/:id/pay", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.ProcessPaymentHandler)
	api.Post("/payments/midtrans-notification", container.OrderHandler.MidtransNotificationHandler)

}
