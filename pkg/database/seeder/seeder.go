package seeder

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
)

func RunSeeders(ctx context.Context, q repository.Querier, log *logger.Logger) error {
	if err := SeedUsers(ctx, q, log); err != nil {
		return err
	}

	if err := SeedCategory(ctx, q, log); err != nil {
		return err
	}

	if err := SeedPaymentMethods(ctx, q, log); err != nil {
		return err
	}
	if err := SeedCancellationReasons(ctx, q, log); err != nil {
		return err
	}

	return nil
}
