package seeder

import (
	categories_repo "POS-kasir/internal/categories/repository"
	"POS-kasir/pkg/logger"
	"context"
)

func SeedCategory(ctx context.Context, q categories_repo.Querier, log logger.ILogger) error {
	kategori := []string{
		"Makanan",
		"Minuman",
		"Camilan",
		"Makanan Penutup",
		"Paket",
	}

	for _, nama := range kategori {
		_, err := q.CreateCategory(ctx, nama)
		if err != nil {
			log.Errorf("gagal menambahkan kategori %s: %v", nama, err)
			continue
		}
	}
	log.Info("seeding kategori selesai")
	return nil
}
