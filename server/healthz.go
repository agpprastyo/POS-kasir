package server

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

// HealthCheckResult represents the result of a single dependency health check.
type HealthCheckResult struct {
	Status    string  `json:"status"`
	LatencyMs float64 `json:"latency_ms,omitempty"`
	Error     string  `json:"error,omitempty"`
}

// HealthResponse represents the overall health check response.
type HealthResponse struct {
	Status  string                       `json:"status"`
	Uptime  string                       `json:"uptime"`
	Version string                       `json:"version"`
	Checks  map[string]HealthCheckResult `json:"checks"`
}

func HealthHandler(app *App) fiber.Handler {
	return func(c fiber.Ctx) error {
		checks := make(map[string]HealthCheckResult)
		overallStatus := "ok"

		// Check PostgreSQL
		pgCheck := checkPostgres(c.RequestCtx(), app)
		checks["postgres"] = pgCheck
		if pgCheck.Status != "ok" {
			overallStatus = "degraded"
		}

		// Check Redis
		redisCheck := checkRedis(c.RequestCtx(), app)
		checks["redis"] = redisCheck
		if redisCheck.Status != "ok" {
			overallStatus = "degraded"
		}

		// Check R2/MinIO
		r2Check := checkR2(c.RequestCtx(), app)
		checks["r2"] = r2Check
		if r2Check.Status != "ok" {
			overallStatus = "degraded"
		}

		// Determine HTTP status code
		httpStatus := fiber.StatusOK
		if overallStatus == "degraded" {
			httpStatus = fiber.StatusServiceUnavailable
		}

		uptime := time.Since(app.StartTime).Round(time.Second).String()

		return c.Status(httpStatus).JSON(HealthResponse{
			Status:  overallStatus,
			Uptime:  uptime,
			Version: "1.0.0",
			Checks:  checks,
		})
	}
}

func checkPostgres(ctx context.Context, app *App) HealthCheckResult {
	start := time.Now()
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := app.DB.GetPool().Ping(pingCtx); err != nil {
		app.Logger.Errorf("Health check: PostgreSQL ping failed: %v", err)
		return HealthCheckResult{
			Status:    "fail",
			LatencyMs: float64(time.Since(start).Milliseconds()),
			Error:     "PostgreSQL unavailable",
		}
	}

	return HealthCheckResult{
		Status:    "ok",
		LatencyMs: float64(time.Since(start).Milliseconds()),
	}
}

func checkRedis(ctx context.Context, app *App) HealthCheckResult {
	start := time.Now()
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if app.Redis == nil {
		return HealthCheckResult{
			Status: "skip",
			Error:  "Redis client not configured",
		}
	}

	if err := app.Redis.Ping(pingCtx).Err(); err != nil {
		app.Logger.Errorf("Health check: Redis ping failed: %v", err)
		return HealthCheckResult{
			Status:    "fail",
			LatencyMs: float64(time.Since(start).Milliseconds()),
			Error:     "Redis unavailable",
		}
	}

	return HealthCheckResult{
		Status:    "ok",
		LatencyMs: float64(time.Since(start).Milliseconds()),
	}
}

func checkR2(ctx context.Context, app *App) HealthCheckResult {
	start := time.Now()

	if app.R2 == nil {
		return HealthCheckResult{
			Status: "skip",
			Error:  "R2 client not configured",
		}
	}

	r2Exists, err := app.R2.BucketExists(ctx)
	if err != nil || !r2Exists {
		errMsg := "R2 bucket unavailable"
		if err != nil {
			app.Logger.Errorf("Health check: R2 bucket check failed: %v", err)
			errMsg = err.Error()
		}
		return HealthCheckResult{
			Status:    "fail",
			LatencyMs: float64(time.Since(start).Milliseconds()),
			Error:     errMsg,
		}
	}

	return HealthCheckResult{
		Status:    "ok",
		LatencyMs: float64(time.Since(start).Milliseconds()),
	}
}

// fiber:context-methods migrated
