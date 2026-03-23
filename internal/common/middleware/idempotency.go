package middleware

import (
	"POS-kasir/internal/common"
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
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

func Idempotency(rdb *redis.Client) fiber.Handler {
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

		ctx := context.Background()
		cacheKey := "idempotency:" + key

		// Check if key already exists
		exists, err := rdb.Exists(ctx, cacheKey).Result()
		if err != nil {
			return c.Next() // Continue if Redis fails
		}

		if exists > 0 {
			return c.Status(fiber.StatusConflict).JSON(common.ErrorResponse{
				Message: "Request is already being processed or has been processed",
			})
		}

		// Set key with 24h expiration
		err = rdb.Set(ctx, cacheKey, time.Now().Format(time.RFC3339), 24*time.Hour).Err()
		if err != nil {
			return c.Next()
		}

		err = c.Next()

		// If the request failed with a server error, allow retrying by deleting the key
		if c.Response().StatusCode() >= 500 {
			rdb.Del(ctx, cacheKey)
		}

		return err
	}
}
