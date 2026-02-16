package server

import (
	"POS-kasir/internal/common/middleware"
)

func SetupRoutes(app *App, container *AppContainer) {
	hltHandler := HealthHandler(app)
	app.FiberApp.Get("/healthz", hltHandler)

	api := app.FiberApp.Group("/api/v1")

	authMiddleware := middleware.AuthMiddleware(app.JWT, app.Logger)

	api.Post("/auth/login", container.AuthHandler.LoginHandler)
	api.Post("/auth/refresh", container.AuthHandler.RefreshHandler)
	api.Get("/auth/me", authMiddleware, container.AuthHandler.ProfileHandler)
	api.Post("/auth/add", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.AuthHandler.AddUserHandler)
	api.Put("/auth/me/avatar", authMiddleware, container.AuthHandler.UpdateAvatarHandler)
	api.Put("/auth/me/password", authMiddleware, container.AuthHandler.UpdatePasswordHandler)
	api.Post("/auth/logout", authMiddleware, container.AuthHandler.LogoutHandler)

	api.Get("/users", authMiddleware, container.UserHandler.GetAllUsersHandler)
	api.Post("/users", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.UserHandler.CreateUserHandler)
	api.Get("/users/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.UserHandler.GetUserByIDHandler)
	api.Put("/users/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.UserHandler.UpdateUserHandler)
	api.Post("/users/:id/toggle-status", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.UserHandler.ToggleUserStatusHandler)
	api.Delete("/users/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.UserHandler.DeleteUserHandler)

	api.Get("/categories", authMiddleware, container.CategoryHandler.GetAllCategoriesHandler)
	api.Get("/categories/count", authMiddleware, container.CategoryHandler.GetCategoryCountHandler)
	api.Post("/categories", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.CategoryHandler.CreateCategoryHandler)
	api.Get("/categories/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.CategoryHandler.GetCategoryByIDHandler)
	api.Put("/categories/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.CategoryHandler.UpdateCategoryHandler)
	api.Delete("/categories/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.CategoryHandler.DeleteCategoryHandler)

	api.Post("/products", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.CreateProductHandler)
	api.Post("/products/:id/image", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.UploadProductImageHandler)

	// Deleted Products Management (Admin only)
	api.Get("/products/trash", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.ProductHandler.ListDeletedProductsHandler)
	api.Post("/products/trash/restore-bulk", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.ProductHandler.RestoreProductsBulkHandler)
	api.Get("/products/trash/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.ProductHandler.GetDeletedProductHandler)
	api.Post("/products/trash/:id/restore", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.ProductHandler.RestoreProductHandler)

	api.Get("/products", authMiddleware, container.ProductHandler.ListProductsHandler)
	api.Get("/products/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.ProductHandler.GetProductHandler)
	api.Get("/products/:id/stock-history", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.GetStockHistoryHandler)
	api.Patch("/products/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.UpdateProductHandler)
	api.Delete("/products/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.ProductHandler.DeleteProductHandler)

	api.Post("/products/:product_id/options", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.CreateProductOptionHandler)
	api.Post("/products/:product_id/options/:option_id/image", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.UploadProductOptionImageHandler)
	api.Patch("/products/:product_id/options/:option_id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.UpdateProductOptionHandler)
	api.Delete("/products/:product_id/options/:option_id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager), container.ProductHandler.DeleteProductOptionHandler)

	api.Get("/payment-methods", authMiddleware, container.PaymentMethodHandler.ListPaymentMethodsHandler)
	api.Get("/cancellation-reasons", authMiddleware, container.CancellationReasonHandler.ListCancellationReasonsHandler)
	api.Get("/activity-logs", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleAdmin), container.ActivityLogHandler.GetActivityLogs)

	api.Post("/orders", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), middleware.ShiftMiddleware(container.ShiftRepo, app.Cache), container.OrderHandler.CreateOrderHandler)
	api.Get("/orders", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.ListOrdersHandler)
	api.Get("/orders/:id", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.GetOrderHandler)
	api.Patch("/orders/:id/items", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.UpdateOrderItemsHandler)

	api.Post("/orders/:id/cancel", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.CancelOrderHandler)
	api.Post("/orders/:id/apply-promotion", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.ApplyPromotionHandler)
	api.Post("/orders/:id/pay/midtrans", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.InitiateMidtransPaymentHandler)
	api.Post("/orders/:id/pay/manual", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.ConfirmManualPaymentHandler)
	api.Post("/orders/:id/update-status", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.UpdateOperationalStatusHandler)
	api.Post("/orders/:id/update-status", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.OrderHandler.UpdateOperationalStatusHandler)
	api.Post("/orders/:id/print", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier), container.PrinterHandler.PrintInvoiceHandler)
	api.Post("/payments/midtrans-notification", container.OrderHandler.MidtransNotificationHandler)

	api.Get("/reports/dashboard-summary", authMiddleware, container.ReportHandler.GetDashboardSummaryHandler)
	api.Get("/reports/sales", authMiddleware, container.ReportHandler.GetSalesReportsHandler)
	api.Get("/reports/products", authMiddleware, container.ReportHandler.GetProductPerformanceHandler)
	api.Get("/reports/payment-methods", authMiddleware, container.ReportHandler.GetPaymentMethodPerformanceHandler)
	api.Get("/reports/cashier-performance", authMiddleware, container.ReportHandler.GetCashierPerformanceHandler)
	api.Get("/reports/cancellations", authMiddleware, container.ReportHandler.GetCancellationReportsHandler)
	api.Get("/reports/profit-summary", authMiddleware, container.ReportHandler.GetProfitSummaryHandler)
	api.Get("/reports/profit-products", authMiddleware, container.ReportHandler.GetProductProfitReportsHandler)
	// Public read access for Cashier (and above)
	promotionsReadGroup := api.Group("/promotions", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleCashier))
	{
		promotionsReadGroup.Get("/", container.PromotionHandler.ListPromotionsHandler)
		promotionsReadGroup.Get("/:id", container.PromotionHandler.GetPromotionHandler)
	}

	// Helper for Manager/Admin operations
	promotionsWriteGroup := api.Group("/promotions", authMiddleware, middleware.RoleMiddleware(middleware.UserRoleManager))
	{
		promotionsWriteGroup.Post("/", container.PromotionHandler.CreatePromotionHandler)
		promotionsWriteGroup.Put("/:id", container.PromotionHandler.UpdatePromotionHandler)
		promotionsWriteGroup.Delete("/:id", container.PromotionHandler.DeletePromotionHandler)
		promotionsWriteGroup.Post("/:id/restore", container.PromotionHandler.RestorePromotionHandler)
	}

	settingsGroup := api.Group("/settings", authMiddleware)
	{
		settingsGroup.Get("/branding", container.SettingsHandler.GetBrandingHandler)
		settingsGroup.Put("/branding", middleware.RoleMiddleware(middleware.UserRoleAdmin), container.SettingsHandler.UpdateBrandingHandler)
		settingsGroup.Post("/branding/logo", middleware.RoleMiddleware(middleware.UserRoleAdmin), container.SettingsHandler.UpdateLogoHandler)

		settingsGroup.Get("/printer", container.SettingsHandler.GetPrinterSettingsHandler)
		settingsGroup.Put("/printer", middleware.RoleMiddleware(middleware.UserRoleAdmin), container.SettingsHandler.UpdatePrinterSettingsHandler)
		settingsGroup.Post("/printer/test", middleware.RoleMiddleware(middleware.UserRoleAdmin), container.PrinterHandler.TestPrintHandler)
	}

	shiftGroup := api.Group("/shifts", authMiddleware)
	{
		shiftGroup.Post("/start", container.ShiftHandler.StartShiftHandler)
		shiftGroup.Post("/end", container.ShiftHandler.EndShiftHandler)
		shiftGroup.Get("/current", container.ShiftHandler.GetOpenShiftHandler)
		shiftGroup.Post("/cash-transaction", container.ShiftHandler.CreateCashTransactionHandler)
	}
}
