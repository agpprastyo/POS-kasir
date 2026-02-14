package seeder

import (
	"POS-kasir/internal/cancellation_reasons/repository"
	cancellation_reasons_repo "POS-kasir/internal/cancellation_reasons/repository"
	"POS-kasir/pkg/logger"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func SeedCancellationReasons(ctx context.Context, q cancellation_reasons_repo.Querier, log logger.ILogger) error {
	log.Info("Seeding cancellation reasons...")

	defaultReasons := []struct {
		Reason      string
		Description string
	}{
		{Reason: "Stok Habis", Description: "Produk yang dipesan tidak tersedia atau stoknya habis."},
		{Reason: "Permintaan Pelanggan", Description: "Pelanggan meminta untuk membatalkan pesanan."},
		{Reason: "Kesalahan Input", Description: "Kasir salah memasukkan item atau jumlah pesanan."},
		{Reason: "Masalah Pembayaran", Description: "Pembayaran gagal atau tidak dapat diproses."},
		{Reason: "Lainnya", Description: "Alasan lain yang tidak tercakup dalam pilihan yang ada."},
	}

	for _, item := range defaultReasons {

		_, err := q.GetCancellationReasonByReason(ctx, item.Reason)
		if err == nil {

			log.Infof("Cancellation reason '%s' already exists, skipping.", item.Reason)
			continue
		}

		if !errors.Is(err, pgx.ErrNoRows) {
			log.Errorf("Failed to check for cancellation reason '%s': %v", item.Reason, err)
			return err
		}

		params := repository.CreateCancellationReasonParams{
			Reason:      item.Reason,
			Description: &item.Description,
		}
		_, createErr := q.CreateCancellationReason(ctx, params)
		if createErr != nil {
			log.Errorf("Failed to seed cancellation reason '%s': %v", item.Reason, createErr)
			return createErr
		}
		log.Infof("Successfully seeded cancellation reason: '%s'", item.Reason)
	}

	log.Info("Cancellation reasons seeding completed.")
	return nil
}
