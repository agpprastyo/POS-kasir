package middleware

import (
	"time"

	"POS-kasir/pkg/cache"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func NewRateLimiter(c cache.Cache, max int, expiration time.Duration) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        max,
		Expiration: expiration,
		Storage:    c,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
	})
}

// GlobalRateLimiter sets a default limit of 100 requests per minute
func GlobalRateLimiter(c cache.Cache) fiber.Handler {
	return NewRateLimiter(c, 100, 1*time.Minute)
}

// StrictRateLimiter sets a strict limit of 20 requests per minute (e.g. for login)
func StrictRateLimiter(c cache.Cache) fiber.Handler {
	return NewRateLimiter(c, 20, 1*time.Minute)
}

// WebhookRateLimiter sets a limit of 30 requests per minute (e.g. for Midtrans)
func WebhookRateLimiter(c cache.Cache) fiber.Handler {
	return NewRateLimiter(c, 30, 1*time.Minute)
}

// MetricsRateLimiter sets a limit of 60 requests per minute (e.g. for Web Vitals)
func MetricsRateLimiter(c cache.Cache) fiber.Handler {
	return NewRateLimiter(c, 60, 1*time.Minute)
}
