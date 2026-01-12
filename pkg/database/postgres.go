package database

import (
	"POS-kasir/config"
	"POS-kasir/pkg/logger"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/jackc/pgx/v5/stdlib"
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

	if cfg.AutoMigrate {
		if err := runMigrations(cfg, dsn, log); err != nil {

			pool.Close()
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

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

func runMigrations(cfg *config.AppConfig, dsn string, log logger.ILogger) error {
	log.Infof("AutoMigrate enabled — running migrations from: %s", cfg.MigrationsPath)

	if cfg.MigrationsPath == "" {
		return fmt.Errorf("migrations path is empty")
	}
	// Ensure the path has the correct scheme (file://)
	migrationSource := cfg.MigrationsPath
	if !strings.HasPrefix(migrationSource, "file://") {
		log.Warnf("MigrationsPath does not start with file:// — automatically adding prefix. Original: %s", migrationSource)
		migrationSource = "file://" + migrationSource
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database/sql connection for migrations: %w", err)
	}
	defer func(sqlDB *sql.DB) {
		err := sqlDB.Close()
		if err != nil {
			log.Errorf("Failed to close database/sql connection for migrations: %v", err)
		}
	}(sqlDB)

	sqlDB.SetMaxOpenConns(cfg.DB.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.DB.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(cfg.DB.MaxLifetime)

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationSource,
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrations up failed: %w", err)
	}

	log.Infof("Database migrations applied (or no change).")
	return nil
}
