package seeder

import (
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/logger"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func SeedPaymentMethods(ctx context.Context, q repository.Querier, log logger.ILogger) error {
	log.Info("Seeding payment methods...")

	defaultMethods := []string{
		"Cash",
		"QRIS Dinamis",
		"QRIS Statis",
	}

	for _, methodName := range defaultMethods {

		_, err := q.GetPaymentMethodByName(ctx, methodName)
		if err == nil {
			log.Infof("Payment method '%s' already exists, skipping.", methodName)
			continue
		}

		if !errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("Failed to check for payment method '%s': %v", methodName, err)
			return err
		}

		_, createErr := q.CreatePaymentMethod(ctx, methodName)
		if createErr != nil {
			log.Errorf("Failed to seed payment method '%s': %v", methodName, createErr)
			return createErr
		}
		log.Infof("Successfully seeded payment method: '%s'", methodName)
	}

	log.Info("Payment methods seeding completed.")
	return nil
}
