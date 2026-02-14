package seeder

import (
	cancellation_reasons_repo "POS-kasir/internal/cancellation_reasons/repository"
	categories_repo "POS-kasir/internal/categories/repository"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
)

func RunSeeders(ctx context.Context, q repository.Querier, catRepo categories_repo.Querier, cancelRepo cancellation_reasons_repo.Querier, log logger.ILogger) error {
	log.Info("Running seeders...")
	if err := SeedUsers(ctx, q, log); err != nil {
		log.Error("Failed to seed users", "error", err)
		return err
	}

	if err := SeedCategory(ctx, catRepo, log); err != nil {
		log.Error("Failed to seed categories", "error", err)
		return err
	}

	if err := SeedPaymentMethods(ctx, q, log); err != nil {
		log.Error("Failed to seed payment methods", "error", err)
		return err
	}
	if err := SeedCancellationReasons(ctx, cancelRepo, log); err != nil {
		log.Error("Failed to seed cancellation reasons", "error", err)
		return err
	}

	log.Info("Seeders completed successfully")
	return nil
}
