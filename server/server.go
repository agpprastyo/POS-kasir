package server

import (
	"POS-kasir/config"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"
	"POS-kasir/pkg/utils"
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

	jwt := app.JWT

	// Setup routes
	SetupRoutes(app.FiberApp, app.Logger, app.DB, app.Config, jwt, app.Minio)

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
