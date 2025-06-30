package database

import (
	"POS-kasir/config"
	"POS-kasir/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type Postgres struct {
	DB     *pgxpool.Pool
	Config *config.AppConfig
	Log    *logger.Logger
}

func NewPostgresPool(cfg *config.AppConfig, log *logger.Logger) (*Postgres, error) {
	// Build PostgreSQL connection string (pgx format)
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host,
		cfg.DB.Port, cfg.DB.DBName, cfg.DB.SSLMode,
	)

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL connection string: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.DB.MaxOpenConn)

	// Create connection pool
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL connection pool: %w", err)
	}

	// Verify the connection
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL database: %w", err)
	}

	log.Info("Successfully connected to PostgreSQL database")

	return &Postgres{
		DB:     pool,
		Config: cfg,
		Log:    log,
	}, nil
}

// ClosePostgresPool closes the pgxpool.Pool connection.
func ClosePostgresPool(p *Postgres) error {
	defer func() {
		if r := recover(); r != nil {
			p.Log.Errorf("Panic while closing database connection: %v", r)
		}
	}()
	p.DB.Close()
	p.Log.Println("Closed database connection")
	return nil
}

// PingPostgresPool pings the database to check connectivity.
func PingPostgresPool(p *Postgres) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := p.DB.Ping(ctx); err != nil {
		p.Log.Errorf("Failed to ping database: %v", err)
		return err
	}
	p.Log.Println("Successfully pinged database")
	return nil
}
