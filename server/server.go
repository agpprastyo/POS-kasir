package server

import (
	"POS-kasir/config"
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/categories"
	"POS-kasir/internal/orders"
	"POS-kasir/internal/products"
	"POS-kasir/internal/repository"
	"POS-kasir/internal/user"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"
	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/utils"
	"POS-kasir/pkg/validator"
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"

	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type App struct {
	Config   *config.AppConfig
	Logger   *logger.Logger
	DB       *database.Postgres
	FiberApp *fiber.App
	JWT      utils.Manager
	Minio    *minio.Minio
}

type AppContainer struct {
	AuthHandler     auth.AthHandler
	UserHandler     user.UsrHandler
	CategoryHandler categories.ICtgHandler
	ProductHandler  products.IPrdHandler
	OrderHandler    orders.IOrderHandler
}

func BuildAppContainer(app *App) *AppContainer {
	val := validator.NewValidator()

	store := repository.NewStore(app.DB.DB, app.Logger)

	// Activity Log Service
	activityService := activitylog.NewService(store, app.Logger)

	// Auth Module
	authRepo := auth.NewAuthRepo(app.Logger, app.Minio)
	authService := auth.NewAuthService(store, app.Logger, app.JWT, authRepo, activityService)
	authHandler := auth.NewAuthHandler(authService, app.Logger, val)

	// User Module
	userService := user.NewUsrService(store, app.Logger, activityService, authRepo)
	userHandler := user.NewUsrHandler(userService, app.Logger, val)

	// Category Module
	categoryService := categories.NewCtgService(store, app.Logger, activityService)
	categoryHandler := categories.NewCtgHandler(categoryService, app.Logger)

	// Product Module
	prdRepo := products.NewPrdRepo(app.Minio, app.Logger)
	prdService := products.NewPrdService(store, app.Logger, prdRepo, activityService)
	prdHandler := products.NewPrdHandler(prdService, app.Logger, val)

	// Order & Payment Module
	midtransService := payment.NewMidtransService(app.Config, app.Logger)
	orderService := orders.NewOrderService(store, midtransService, activityService, app.Logger)
	orderHandler := orders.NewOrderHandler(orderService, app.Logger, val)

	return &AppContainer{
		AuthHandler:     *authHandler,
		UserHandler:     *userHandler,
		CategoryHandler: categoryHandler,
		ProductHandler:  prdHandler,
		OrderHandler:    orderHandler,
	}
}

func InitApp() *App {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.New(cfg)

	db, err := database.NewPostgresPool(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	// Create Fiber app
	fiberApp := fiber.New(fiber.Config{
		AppName:      cfg.Server.AppName,
		ErrorHandler: CustomErrorHandler(log),
	})
	// Initialize JWT manager
	jwtManager := utils.NewJWTManager(cfg)

	newMinio, err := minio.NewMinio(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize Minio: %v", err)
	}

	return &App{
		Config:   cfg,
		Logger:   log,
		DB:       db,
		FiberApp: fiberApp,
		JWT:      jwtManager,
		Minio:    newMinio,
	}
}

func StartServer(app *App) {
	// Setup middleware
	SetupMiddleware(app)

	container := BuildAppContainer(app)

	SetupRoutes(app, container, app.JWT)

	// Start app
	app.Logger.Infof("Starting app on port %s...", app.Config.Server.Port)
	if err := app.FiberApp.Listen(":" + app.Config.Server.Port); err != nil {
		app.Logger.Fatalf("Error starting app: %v", err)
	}
}

func SetupMiddleware(app *App) {
	app.FiberApp.Use(fiberlog.New())
	app.FiberApp.Use(recover.New())

}

func Cleanup(app *App) {
	if app.DB != nil {
		if err := database.ClosePostgresPool(app.DB); err != nil {
			app.Logger.Errorf("Error closing database connection: %v", err)
		}
	}

}

func WaitForShutdown(app *App) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.Logger.Info("Shutting down app...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := app.FiberApp.ShutdownWithContext(ctx); err != nil {
		app.Logger.Fatalf("Server shutdown failed: %v", err)
	}
}

func CustomErrorHandler(logger *logger.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		logger.Errorf("Error 1: %v", err)

		var e *fiber.Error
		if errors.As(err, &e) {
			logger.Errorf("Fiber error 1: %v", e)
			return c.Status(e.Code).JSON(fiber.Map{
				"error": e.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})

	}
}
