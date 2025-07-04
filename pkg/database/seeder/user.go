package seeder

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"context"

	"github.com/google/uuid"
)

func SeedUsers(ctx context.Context, q *repository.Queries, log *logger.Logger) error {
	userData := []struct {
		Username string
		Email    string
		Role     repository.UserRole
	}{
		{"admin", "admin@example.com", repository.UserRoleAdmin},
		{"cashier", "cashier@example.com", repository.UserRoleCashier},
		{"manager", "manager@example.com", repository.UserRoleManager},
	}

	hashPassword, err := utils.HashPassword("passwordrahasia")
	if err != nil {
		log.Fatalf("Seeder User | Error hashing password: %v", err)
		return err
	}

	for _, data := range userData {
		userUUID, err := uuid.NewV7()
		if err != nil {
			log.Fatalf("Seeder User | failed to generate UUID: %v", err)
			return err
		}
		params := repository.CreateUserParams{
			ID:           userUUID,
			Username:     data.Username,
			Email:        data.Email,
			PasswordHash: hashPassword,
			Avatar:       nil,
			Role:         data.Role,
		}
		_, err = q.CreateUser(ctx, params)
		if err != nil {
			log.Printf("Seeder User | failed to seed user %s: %v", data.Email, err)
			return err
		}
	}
	return nil
}
