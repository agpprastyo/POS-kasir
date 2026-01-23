package main

import (
	"POS-kasir/config"
	"POS-kasir/internal/repository"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/database/seeder"
	"POS-kasir/pkg/logger"
	"context"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	cfg := config.Load()

	ctx := context.Background()

	logr := logger.New(cfg)

	db, err := database.NewDatabase(cfg, logr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	queries := repository.New(db.GetPool())

	if err := seeder.RunSeeders(ctx, queries, logr); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding completed successfully.")
}
