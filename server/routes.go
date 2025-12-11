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
	api.Get("/categories/count", authMiddleware, container.CategoryHandler.GetCategoryCountHandler)
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

	api.Get("/payment-methods", authMiddleware, container.PaymentMethodHandler.ListPaymentMethodsHandler)
	api.Get("/cancellation-reasons", authMiddleware, container.CancellationReasonHandler.ListCancellationReasonsHandler)

	api.Post("/orders", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.CreateOrderHandler)
	api.Get("/orders", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.ListOrdersHandler)
	api.Get("/orders/:id", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.GetOrderHandler)
	api.Patch("/orders/:id/items", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.UpdateOrderItemsHandler)

	api.Post("/orders/:id/cancel", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.CancelOrderHandler)
	//api.Post("/orders/:id/apply-promotion", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.ApplyPromotionHandler)
	api.Post("/orders/:id/pay", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.ProcessPaymentHandler)
	api.Post("/orders/:id/complete-payment", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.CompleteManualPaymentHandler)
	api.Patch("/orders/:id/status", authMiddleware, middleware.RoleMiddleware(repository.UserRoleCashier), container.OrderHandler.UpdateOperationalStatusHandler)
	api.Post("/payments/midtrans-notification", container.OrderHandler.MidtransNotificationHandler)

	api.Get("/reports/dashboard-summary", authMiddleware, container.ReportHandler.GetDashboardSummaryHandler)
	api.Get("/reports/sales", authMiddleware, container.ReportHandler.GetSalesReportsHandler)
	api.Get("/reports/products?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD", authMiddleware, container.ReportHandler.GetProductPerformanceHandler)
	api.Get("/reports/payment-methods?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD", authMiddleware, container.ReportHandler.GetPaymentMethodPerformanceHandler)
	api.Get("/reports/cashier-performance?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD", authMiddleware, container.ReportHandler.GetCashierPerformanceHandler)
	api.Get("/reports/cancellations?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD", authMiddleware, container.ReportHandler.GetCancellationReportsHandler)

	////api.Get("/reports/sales", authMiddleware, container.ReportHandler.GetSalesSummaryHandler)
	////api.Get("/reports/products", authMiddleware, container.ReportHandler.GetSalesDetailHandler)
	//
	//// --- RUTE TAMBAHAN UNTUK MANAJEMEN DATA MASTER ---
	//// CRUD untuk mengelola data master seperti metode pembayaran dan alasan pembatalan.
	//// Dibutuhkan peran: Admin
	//masterDataGroup := api.Group("/", authMiddleware, middleware.RoleMiddleware(repository.UserRoleAdmin))
	//{
	//	// Payment Methods
	//	masterDataGroup.Post("/payment-methods", container.PaymentMethodHandler.CreatePaymentMethodHandler)
	//	masterDataGroup.Put("/payment-methods/:id", container.PaymentMethodHandler.UpdatePaymentMethodHandler)
	//
	//	// Cancellation Reasons
	//	masterDataGroup.Post("/cancellation-reasons", container.CancellationReasonHandler.CreateCancellationReasonHandler)
	//	masterDataGroup.Put("/cancellation-reasons/:id", container.CancellationReasonHandler.UpdateCancellationReasonHandler)
	//}
	//
	//// --- RUTE TAMBAHAN UNTUK MANAJEMEN PROMOSI ---
	//// CRUD lengkap untuk mengelola promosi.
	//// Dibutuhkan peran: Admin / Manager
	//promotionsGroup := api.Group("/promotions", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager))
	//{
	//	// Membuat promosi baru beserta aturan dan targetnya.
	//	promotionsGroup.Post("/", container.PromotionHandler.CreatePromotionHandler)
	//	// Mendapatkan daftar semua promosi.
	//	promotionsGroup.Get("/", container.PromotionHandler.ListPromotionsHandler)
	//	// Mendapatkan detail satu promosi, termasuk aturan dan targetnya.
	//	promotionsGroup.Get("/:id", container.PromotionHandler.GetPromotionHandler)
	//	// Memperbarui promosi, aturan, dan targetnya.
	//	promotionsGroup.Put("/:id", container.PromotionHandler.UpdatePromotionHandler)
	//	// Menghapus promosi.
	//	promotionsGroup.Delete("/:id", container.PromotionHandler.DeletePromotionHandler)
	//}
	//
	//// --- RUTE TAMBAHAN UNTUK DASHBOARD & REPORTING ---
	//// Rute ini penting untuk insight bisnis dan biasanya hanya untuk Manajer/Admin.
	//reportingGroup := api.Group("/reports", authMiddleware, middleware.RoleMiddleware(repository.UserRoleManager))
	//{
	//	// Mendapatkan data ringkasan untuk dashboard utama.
	//	// Contoh: total penjualan hari ini, jumlah transaksi, produk terlaris.
	//	// GET /api/v1/reports/dashboard-summary
	//	reportingGroup.Get("/dashboard-summary", container.ReportHandler.GetDashboardSummaryHandler)
	//
	//	// Mendapatkan laporan penjualan dengan rentang tanggal.
	//	// GET /api/v1/reports/sales?start_date=2025-07-01&end_date=2025-07-18
	//	reportingGroup.Get("/sales", container.ReportHandler.GetSalesReportsHandler)
	//
	//	// Mendapatkan laporan performa produk.
	//	// GET /api/v1/reports/products
	//	reportingGroup.Get("/products", container.ReportHandler.GetProductPerformanceHandler)
	//}
}
