package middleware

import (
	"POS-kasir/internal/common"
	"time"

	"POS-kasir/pkg/cache"

	"github.com/gofiber/fiber/v3"
)

func RequireIdempotencyKey() fiber.Handler {
	return func(c fiber.Ctx) error {
		key := c.Get("X-Idempotency-Key")
		if key == "" {
			return c.Status(fiber.StatusBadRequest).JSON(common.ErrorResponse{
				Message: "X-Idempotency-Key header is required for this request",
			})
		}
		return c.Next()
	}
}

func Idempotency(cCache cache.Cache) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Only check for POST, PUT, PATCH, DELETE
		method := c.Method()
		if method == fiber.MethodGet || method == fiber.MethodHead || method == fiber.MethodOptions {
			return c.Next()
		}

		key := c.Get("X-Idempotency-Key")
		if key == "" {
			return c.Next()
		}

		cacheKey := "idempotency:" + key

		// Check if key already exists
		exists, err := cCache.Exists(cacheKey)
		if err != nil {
			return c.Next() // Continue if Redis fails
		}

		if exists {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Request is already being processed or has been processed",
			})
		}

		// Set key with 24h expiration
		err = cCache.Set(cacheKey, []byte(time.Now().Format(time.RFC3339)), 24*time.Hour)
		if err != nil {
			return c.Next()
		}

		err = c.Next()

		// If the request failed with a server error, allow retrying by deleting the key
		if c.Response().StatusCode() >= 500 {
			cCache.Delete(cacheKey)
		}

		return err
	}
}
