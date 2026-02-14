package server

import (
	"POS-kasir/config"
	"POS-kasir/internal/activitylog"
	activitylog_repo "POS-kasir/internal/activitylog/repository"
	"POS-kasir/internal/auth"
	"POS-kasir/internal/cancellation_reasons"
	cancellation_reasons_repo "POS-kasir/internal/cancellation_reasons/repository"
	"POS-kasir/internal/categories"
	categories_repo "POS-kasir/internal/categories/repository"
	"POS-kasir/internal/orders"
	"POS-kasir/internal/payment_methods"
	"POS-kasir/internal/printer"
	"POS-kasir/internal/products"
	products_repo "POS-kasir/internal/products/repository"
	"POS-kasir/internal/promotions"
	"POS-kasir/internal/report"
	"POS-kasir/internal/repository"
	"POS-kasir/internal/settings"
	"POS-kasir/internal/shift"
	"POS-kasir/internal/user"
	"POS-kasir/pkg/cache"
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/payment"
	"POS-kasir/pkg/utils"
	"POS-kasir/pkg/validator"
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	swagger "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	fiberlog "github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"

	_ "embed"
	"html/template"

	"POS-kasir/sqlc/migrations"
)

//go:embed swagger_script.js
var swaggerScript string

type App struct {
	Config          *config.AppConfig
	Logger          logger.ILogger
	DB              database.IDatabase
	FiberApp        *fiber.App
	JWT             utils.Manager
	Store           repository.Store
	Validator       validator.Validator
	MidtransService payment.IMidtrans
	R2              cloudflarer2.IR2
	Cache           *shift.Cache
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
	PromotionHandler          promotions.IPromotionHandler
	ActivityLogHandler        *activitylog.ActivityLogHandler
	SettingsHandler           *settings.SettingsHandler
	PrinterHandler            *printer.PrinterHandler
	ShiftHandler              shift.Handler
}

func InitApp() *App {
	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Println("Warning: .env file not found, using system environment variables")
		}
	}

	cfg := config.Load()
	log := logger.New(cfg)

	db, err := database.NewDatabase(cfg, log, migrations.FS)
	if err != nil {
		log.Fatalf("Failed to initialize PostgreSQL: %v", err)
	}

	fiberApp := fiber.New(fiber.Config{
		AppName:      cfg.Server.AppName,
		ErrorHandler: CustomErrorHandler(log),
	})

	jwtManager := utils.NewJWTManager(cfg)

	store := repository.NewStore(db.GetPool(), log)
	val := validator.NewValidator()
	midtransService := payment.NewMidtransService(cfg, log)

	newR2, err := cloudflarer2.NewCloudflareR2(cfg, log)
	if err != nil {
		log.Errorf("Failed to initialize Cloudflare R2: %v", err)
	}

	memCache := cache.NewMemoryCache()
	shiftCache := shift.NewCache(memCache)

	return &App{
		Config:          cfg,
		Logger:          log,
		DB:              db,
		FiberApp:        fiberApp,
		JWT:             jwtManager,
		Store:           store,
		Validator:       val,
		MidtransService: midtransService,
		R2:              newR2,
		Cache:           shiftCache,
	}
}

func BuildAppContainer(app *App) *AppContainer {
	// Activity Log IActivityService
	activityLogRepo := activitylog_repo.New(app.DB.GetPool())
	activityService := activitylog.NewActivityService(activityLogRepo, app.Logger)
	activityLogHandler := activitylog.NewActivityLogHandler(activityService, app.Logger, app.Validator)

	// Auth Module
	authRepo := auth.NewAuthRepo(app.Logger, app.R2)
	authService := auth.NewAuthService(app.Store, app.Logger, app.JWT, authRepo, activityService)
	authHandler := auth.NewAuthHandler(authService, app.Logger, app.Validator, app.Config)

	// User Module
	userService := user.NewUsrService(app.Store, app.Logger, activityService, authRepo)
	userHandler := user.NewUsrHandler(userService, app.Logger, app.Validator)

	// Category Module
	categoryRepo := categories_repo.New(app.DB.GetPool())
	categoryService := categories.NewCtgService(categoryRepo, app.Logger, activityService)
	categoryHandler := categories.NewCtgHandler(categoryService, app.Logger, app.Validator)

	// Product Module
	prdRepo := products.NewPrdRepo(app.R2, app.Logger)
	productsRepo := products_repo.New(app.DB.GetPool())
	prdService := products.NewPrdService(app.Store, productsRepo, app.Logger, prdRepo, activityService)
	prdHandler := products.NewPrdHandler(prdService, app.Logger, app.Validator)

	// Order & Payment Module
	orderService := orders.NewOrderService(app.Store, app.MidtransService, activityService, app.Logger)
	orderHandler := orders.NewOrderHandler(orderService, app.Logger, app.Validator)

	// Payment Method Module
	paymentMethodService := payment_methods.NewPaymentMethodService(app.Store, app.Logger)
	paymentMethodHandler := payment_methods.NewPaymentMethodHandler(paymentMethodService, app.Logger)

	// Cancellation Reason Module
	cancellationRepo := cancellation_reasons_repo.New(app.DB.GetPool())
	cancellationReasonService := cancellation_reasons.NewCancellationReasonService(cancellationRepo, app.Logger)
	cancellationReasonHandler := cancellation_reasons.NewCancellationReasonHandler(cancellationReasonService, app.Logger)

	// report module
	reportService := report.NewRptService(app.Store, activityService, app.Logger)
	reportHandler := report.NewRptHandler(reportService, app.Logger)

	// Promotion Module
	promotionService := promotions.NewPromotionService(app.Store, app.Logger)
	promotionHandler := promotions.NewPromotionHandler(promotionService, app.Logger, app.Validator)

	// Settings Module
	settingsService := settings.NewSettingsService(app.Store, app.R2, app.Logger)
	settingsHandler := settings.NewSettingsHandler(settingsService)

	// Printer Module
	printerService := printer.NewPrinterService(orderService, settingsService, app.Store, app.Logger)
	printerHandler := printer.NewPrinterHandler(printerService)

	// Shift Module
	shiftService := shift.NewService(app.Store, app.Logger, app.Cache)
	shiftHandler := shift.NewHandler(shiftService, app.Logger, app.Validator)

	return &AppContainer{
		AuthHandler:               authHandler,
		UserHandler:               userHandler,
		CategoryHandler:           categoryHandler,
		ProductHandler:            prdHandler,
		OrderHandler:              orderHandler,
		PaymentMethodHandler:      paymentMethodHandler,
		CancellationReasonHandler: cancellationReasonHandler,
		ReportHandler:             reportHandler,
		PromotionHandler:          promotionHandler,
		ActivityLogHandler:        activityLogHandler,
		SettingsHandler:           settingsHandler,
		PrinterHandler:            printerHandler,
		ShiftHandler:              shiftHandler,
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

	origins := strings.TrimSpace(app.Config.Server.CorsAllowOrigins)

	if origins == "" {
		log.Fatal("CORS_ALLOW_ORIGINS is empty or invalid")
	}
	app.FiberApp.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(origins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	app.FiberApp.Use(fiberlog.New())
	app.FiberApp.Use(recover.New())

	app.FiberApp.Get("/swagger/*", swagger.New(
		swagger.Config{
			URL:            "/swagger/doc.json",
			CustomScript:   template.JS(swaggerScript),
			ShowExtensions: true,
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
	return func(c fiber.Ctx, err error) error {
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
