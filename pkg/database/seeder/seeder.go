package seeder

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
)

func RunSeeders(ctx context.Context, q repository.Querier, log logger.ILogger) error {
	log.Info("Running seeders...")
	if err := SeedUsers(ctx, q, log); err != nil {
		log.Error("Failed to seed users", "error", err)
		return err
	}

	if err := SeedCategory(ctx, q, log); err != nil {
		log.Error("Failed to seed categories", "error", err)
		return err
	}

	if err := SeedPaymentMethods(ctx, q, log); err != nil {
		log.Error("Failed to seed payment methods", "error", err)
		return err
	}
	if err := SeedCancellationReasons(ctx, q, log); err != nil {
		log.Error("Failed to seed cancellation reasons", "error", err)
		return err
	}

	log.Info("Seeders completed successfully")
	return nil
}
