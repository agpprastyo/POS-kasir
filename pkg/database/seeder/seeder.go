package seeder

import (
	cancellation_reasons_repo "POS-kasir/internal/cancellation_reasons/repository"
	categories_repo "POS-kasir/internal/categories/repository"
	payment_methods_repo "POS-kasir/internal/payment_methods/repository"
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"

	user_repo "POS-kasir/internal/user/repository"
	"POS-kasir/pkg/logger"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunSeeders(ctx context.Context, pool *pgxpool.Pool, userRepo user_repo.Querier, catRepo categories_repo.Querier, pmRepo payment_methods_repo.Querier, cancelRepo cancellation_reasons_repo.Querier, r2Client cloudflarer2.IR2, log logger.ILogger) error {
	log.Info("Running seeders...")
	if err := SeedUsers(ctx, userRepo, log); err != nil {
		log.Error("Failed to seed users", "error", err)
		return err
	}

	if err := SeedCategory(ctx, catRepo, log); err != nil {
		log.Error("Failed to seed categories", "error", err)
		return err
	}

	if err := SeedPaymentMethods(ctx, pmRepo, log); err != nil {
		log.Error("Failed to seed payment methods", "error", err)
		return err
	}
	if err := SeedCancellationReasons(ctx, cancelRepo, log); err != nil {
		log.Error("Failed to seed cancellation reasons", "error", err)
		return err
	}

	// Seed production data (products, variants, customers, promotions, orders)
	if err := SeedProductionData(ctx, pool, r2Client, log); err != nil {
		log.Error("Failed to seed production data", "error", err)
		return err
	}

	log.Info("Seeders completed successfully")
	return nil
}
