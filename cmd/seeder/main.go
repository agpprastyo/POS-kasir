package main

import (
	"POS-kasir/config"
	cancellation_reasons_repo "POS-kasir/internal/cancellation_reasons/repository"
	categories_repo "POS-kasir/internal/categories/repository"
	payment_methods_repo "POS-kasir/internal/payment_methods/repository"
	user_repo "POS-kasir/internal/user/repository"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/database/seeder"
	"POS-kasir/pkg/logger"
	"POS-kasir/sqlc/migrations"
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	envFile := ".env"
	if len(os.Args) > 1 {
		envFile = os.Args[1]
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Printf("Warning: Error loading %s file: %v", envFile, err)
	}

	cfg := config.Load()

	ctx := context.Background()

	logr := logger.New(cfg)

	db, err := database.NewDatabase(cfg, logr, migrations.FS)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := user_repo.New(db.GetPool())
	catRepo := categories_repo.New(db.GetPool())
	paymentMethodRepo := payment_methods_repo.New(db.GetPool())
	cancelRepo := cancellation_reasons_repo.New(db.GetPool())

	if err := seeder.RunSeeders(ctx, userRepo, catRepo, paymentMethodRepo, cancelRepo, logr); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding completed successfully.")
}
