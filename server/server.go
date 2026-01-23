package server

import (
	"POS-kasir/config"
	"POS-kasir/internal/activitylog"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/cancellation_reasons"
	"POS-kasir/internal/categories"
	"POS-kasir/internal/orders"
	"POS-kasir/internal/payment_methods"
	"POS-kasir/internal/products"
	"POS-kasir/internal/report"
	"POS-kasir/internal/repository"
	"POS-kasir/internal/user"
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/minio"
	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/utils"
	"POS-kasir/pkg/validator"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberlog "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"

	"github.com/gofiber/swagger"
)

type App struct {
	Config          *config.AppConfig
	Logger          logger.ILogger
	DB              database.IDatabase
	FiberApp        *fiber.App
	JWT             utils.Manager
	Minio           minio.IMinio
	Store           repository.Store
	Validator       validator.Validator
	MidtransService payment.IMidtrans
	R2              cloudflarer2.IR2
}

type AppContainer struct {
	AuthHandler               auth.IAuthHandler
	UserHandler               user.IUsrHandler
	CategoryHandler           categories.ICtgHandler
	ProductHandler            products.IPrdHandler
	OrderHandler              orders.IOrderHandler
	PaymentMethodHandler      payment_methods.IPaymentMethodHandler
	CancellationReasonHandler cancellation_reasons.ICancellationReasonHandler
	ReportHandler             report.IRptHandler
}

func InitApp() *App {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	cfg := config.Load()
	log := logger.New(cfg)

	db, err := database.NewDatabase(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	fiberApp := fiber.New(fiber.Config{
		AppName:      cfg.Server.AppName,
		ErrorHandler: CustomErrorHandler(log),
	})

	jwtManager := utils.NewJWTManager(cfg)
	newMinio, err := minio.NewMinio(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize Minio: %v", err)
	}

	store := repository.NewStore(db.GetPool(), log)
	val := validator.NewValidator()
	midtransService := payment.NewMidtransService(cfg, log)

	newR2, err := cloudflarer2.NewCloudflareR2(cfg, log)
	if err != nil {
		log.Errorf("Failed to initialize Cloudflare R2 (continuing without it): %v", err)
	}

	if err != nil {
		log.Warnf("Failed to initialize Cloudflare R2: %v", err)
	}

	return &App{
		Config:          cfg,
		Logger:          log,
		DB:              db,
		FiberApp:        fiberApp,
		JWT:             jwtManager,
		Minio:           newMinio,
		Store:           store,
		Validator:       val,
		MidtransService: midtransService,
		R2:              newR2,
	}
}

func BuildAppContainer(app *App) *AppContainer {
	// Activity Log IActivityService
	activityService := activitylog.NewActivityService(app.Store, app.Logger)

	// Auth Module
	authRepo := auth.NewAuthRepo(app.Logger, app.Minio, app.R2)
	authService := auth.NewAuthService(app.Store, app.Logger, app.JWT, authRepo, activityService)
	authHandler := auth.NewAuthHandler(authService, app.Logger, app.Validator, app.Config)

	// User Module
	userService := user.NewUsrService(app.Store, app.Logger, activityService, authRepo)
	userHandler := user.NewUsrHandler(userService, app.Logger, app.Validator)

	// Category Module
	categoryService := categories.NewCtgService(app.Store, app.Logger, activityService)
	categoryHandler := categories.NewCtgHandler(categoryService, app.Logger)

	// Product Module
	prdRepo := products.NewPrdRepo(app.Minio, app.Logger)
	prdService := products.NewPrdService(app.Store, app.Logger, prdRepo, activityService)
	prdHandler := products.NewPrdHandler(prdService, app.Logger, app.Validator)

	// Order & Payment Module
	orderService := orders.NewOrderService(app.Store, app.MidtransService, activityService, app.Logger)
	orderHandler := orders.NewOrderHandler(orderService, app.Logger, app.Validator)

	// Payment Method Module
	paymentMethodService := payment_methods.NewPaymentMethodService(app.Store, app.Logger)
	paymentMethodHandler := payment_methods.NewPaymentMethodHandler(paymentMethodService, app.Logger)

	// Cancellation Reason Module
	cancellationReasonService := cancellation_reasons.NewCancellationReasonService(app.Store, app.Logger)
	cancellationReasonHandler := cancellation_reasons.NewCancellationReasonHandler(cancellationReasonService, app.Logger)

	// report module
	reportService := report.NewRptService(app.Store, activityService, app.Logger)
	reportHandler := report.NewRptHandler(reportService, app.Logger)

	return &AppContainer{
		AuthHandler:               authHandler,
		UserHandler:               userHandler,
		CategoryHandler:           categoryHandler,
		ProductHandler:            prdHandler,
		OrderHandler:              orderHandler,
		PaymentMethodHandler:      paymentMethodHandler,
		CancellationReasonHandler: cancellationReasonHandler,
		ReportHandler:             reportHandler,
	}
}

func StartServer(app *App) {

	SetupMiddleware(app)
	container := BuildAppContainer(app)

	SetupRoutes(app, container)

	app.Logger.Infof("Starting app on port %s...", app.Config.Server.Port)
	if err := app.FiberApp.Listen(":" + app.Config.Server.Port); err != nil {
		app.Logger.Fatalf("Error starting app: %v", err)
	}
}

func SetupMiddleware(app *App) {
	app.FiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://127.0.0.1:3000, http://100.72.137.78:3000, http://192.168.0.100:3000, http://172.18.0.1:3000",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Access-Control-Allow-Methods",
		AllowCredentials: true,
		MaxAge:           300,
	}))
	app.FiberApp.Use(fiberlog.New())
	app.FiberApp.Use(recover.New())

	app.FiberApp.Get("/swagger/*", swagger.New(
		swagger.Config{
			URL: "/swagger/doc.json",
		}))

	app.Logger.Infof("Swagger UI available at http://localhost:%s/swagger/index.html", app.Config.Server.Port)
}

func Cleanup(app *App) {
	if app.DB != nil {
		app.DB.Close()
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

func CustomErrorHandler(logger logger.ILogger) fiber.ErrorHandler {
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
