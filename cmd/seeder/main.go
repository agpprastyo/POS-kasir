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

	// Initialize logger
	logr := logger.New(cfg)

	// Initialize database connection (replace with your actual DB connection code)
	db, err := database.NewDatabase(cfg, logr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize sqlc Queries
	queries := repository.New(db.GetPool())

	// Run seeders
	if err := seeder.RunSeeders(ctx, queries, logr); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	log.Println("Seeding completed successfully.")
}
