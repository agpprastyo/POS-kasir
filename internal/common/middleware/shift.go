package middleware

import (
	"POS-kasir/internal/common"
	shift_repo "POS-kasir/internal/shift/repository"

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

func ShiftMiddleware(queries shift_repo.Querier, cache ShiftCache) fiber.Handler {
	return func(c fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(uuid.UUID)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(common.ErrorResponse{
				Message: "Unauthorized",
			})
		}

		// Check cache first
		hasShift, found := cache.GetOpen(userID)
		if found {
			if !hasShift {
				return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
					Message: "You must have an open shift to perform this action",
				})
			}
			return c.Next()
		}

		// Check DB if not in cache
		_, err := queries.GetOpenShiftByUserID(c.RequestCtx(), userID)
		if err != nil {
			if err == pgx.ErrNoRows {
				cache.SetOpen(userID, false) // Cache the absence of shift
				return c.Status(fiber.StatusForbidden).JSON(common.ErrorResponse{
					Message: "You must have an open shift to perform this action",
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(common.ErrorResponse{
				Message: "Internal Server Error",
			})
		}

		cache.SetOpen(userID, true) // Cache the presence of shift
		return c.Next()
	}
}

// fiber:context-methods migrated
