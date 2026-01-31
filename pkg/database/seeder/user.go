package seeder

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"POS-kasir/pkg/utils"
	"context"

	"github.com/google/uuid"
)

func SeedUsers(ctx context.Context, q repository.Querier, log logger.ILogger) error {
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
			continue 
		}
		params := repository.CreateUserParams{
			ID:           userUUID,
			Username:     data.Username,
			Email:        data.Email,
			PasswordHash: hashPassword,
			Avatar:       nil,
			Role:         data.Role,
			IsActive:     true,
		}
		_, err = q.CreateUser(ctx, params)
		if err != nil {
			log.Infof("Seeder User | failed to seed user %s: %v", data.Email, err)
			continue
		}
	}
	return nil
}
