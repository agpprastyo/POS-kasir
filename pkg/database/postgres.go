// File: pkg/database/postgres.go
package database

import (
	"POS-kasir/config"
	"POS-kasir/pkg/logger"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type IDatabase interface {
	GetPool() *pgxpool.Pool
	Ping(ctx context.Context) error
	Close()
}

type postgresService struct {
	DB     *pgxpool.Pool
	Config *config.AppConfig
	Log    logger.ILogger
}

func NewDatabase(cfg *config.AppConfig, log logger.ILogger) (IDatabase, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host,
		cfg.DB.Port, cfg.DB.DBName, cfg.DB.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PostgreSQL connection string: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.DB.MaxOpenConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create PostgreSQL connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL database: %w", err)
	}

	log.Infof("Successfully connected to PostgreSQL database")

	return &postgresService{
		DB:     pool,
		Config: cfg,
		Log:    log,
	}, nil
}

func (s *postgresService) GetPool() *pgxpool.Pool {
	return s.DB
}

func (s *postgresService) Close() {
	defer func() {
		if r := recover(); r != nil {
			s.Log.Errorf("Panic while closing database connection: %v", r)
		}
	}()
	s.DB.Close()
	s.Log.Println("Closed database connection")
}

func (s *postgresService) Ping(ctx context.Context) error {
	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := s.DB.Ping(pingCtx); err != nil {
		s.Log.Errorf("Failed to ping database: %v", err)
		return err
	}
	s.Log.Println("Successfully pinged database")
	return nil
}
