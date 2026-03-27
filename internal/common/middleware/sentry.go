package middleware

import (
	"time"

	"github.com/gofiber/contrib/v3/sentry"
	"github.com/gofiber/fiber/v3"
)

func NewSentryMiddleware() fiber.Handler {
	return sentry.New(sentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
		Timeout:         2 * time.Second,
	})
}
