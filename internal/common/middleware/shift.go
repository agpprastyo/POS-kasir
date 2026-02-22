package middleware

import (
	"POS-kasir/internal/common"
	shift_repo "POS-kasir/internal/shift/repository"
	"POS-kasir/pkg/logger"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ShiftCache defines the interface for shift caching.
// It is implemented by internal/shift/cache.go
type ShiftCache interface {
	GetOpen(userID uuid.UUID) (bool, bool)
	SetOpen(userID uuid.UUID, open bool)
}

func ShiftMiddleware(queries shift_repo.Querier, cache ShiftCache, log logger.ILogger) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uuid.UUID)
		if !ok {
			log.Error("Shift Middleware | User ID not found in context", "user_id", userID)
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "Unauthorized",
			})
		}

		// Check cache first
		hasShift, found := cache.GetOpen(userID)
		if found {
			if !hasShift {
				log.Error("Shift Middleware | User has no open shift: ", "user_id", userID)
				return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
					Message: "You must have an open shift to perform this action",
				})
			}
			return c.Next()
		}

		// Check DB if not in cache
		_, err := queries.GetOpenShiftByUserID(c.RequestCtx(), userID)
		if err != nil {
			log.Error("Shift Middleware | Error getting open shift", "user_id", userID, "error", err)
			if err == pgx.ErrNoRows {
				cache.SetOpen(userID, false)
				log.Error("Shift Middleware | User has no open shift", "user_id", userID)
				return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
					Message: "You must have an open shift to perform this action",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Internal Server Error",
			})
		}

		cache.SetOpen(userID, true)
		return c.Next()
	}
}
