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
		{"user1", "user1@example.com", repository.UserRoleCashier},
		{"user2", "user2@example.com", repository.UserRoleCashier},
		{"user3", "user3@example.com", repository.UserRoleManager},
		{"user4", "user4@example.com", repository.UserRoleManager},
		{"user5", "user5@example.com", repository.UserRoleAdmin},
		{"user6", "user6@example.com", repository.UserRoleCashier},
		{"user7", "user7@example.com", repository.UserRoleManager},
		{"user8", "user8@example.com", repository.UserRoleAdmin},
		{"user9", "user9@example.com", repository.UserRoleCashier},
		{"user10", "user10@example.com", repository.UserRoleManager},
		{"user11", "user11@example.com", repository.UserRoleAdmin},
		{"user12", "user12@example.com", repository.UserRoleCashier},
		{"user13", "user13@example.com", repository.UserRoleManager},
		{"user14", "user14@example.com", repository.UserRoleAdmin},
		{"user15", "user15@example.com", repository.UserRoleCashier},
		{"user16", "user16@example.com", repository.UserRoleManager},
		{"user17", "user17@example.com", repository.UserRoleAdmin},
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
			continue // skip this user and continue
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
			log.Printf("Seeder User | failed to seed user %s: %v", data.Email, err)
			continue // skip error and continue with next user
		}
	}
	return nil
}
