package middleware

import (
	"time"

	"POS-kasir/pkg/cache"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
)

func RateLimiter(c cache.Cache) fiber.Handler {
	return limiter.New(limiter.Config{
		Max:               100,
		Expiration:        1 * time.Minute,
		Storage:           c,
		KeyGenerator: func(c fiber.Ctx) string {
			return c.IP()
		},
	})
}
