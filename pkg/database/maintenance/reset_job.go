package maintenance

import (
	"context"

	cancellation_reasons_repo "POS-kasir/internal/cancellation_reasons/repository"
	categories_repo "POS-kasir/internal/categories/repository"
	payment_methods_repo "POS-kasir/internal/payment_methods/repository"
	user_repo "POS-kasir/internal/user/repository"
	cloudflarer2 "POS-kasir/pkg/cloudflare-r2"
	"POS-kasir/pkg/database"
	"POS-kasir/pkg/database/seeder"
	"POS-kasir/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DatabaseResetJob implements the cron job contract for scheduled database resets.
//
// This job is intended for development/demo environments where the database is periodically
// wiped and re-seeded to maintain a consistent state. It should not be enabled in production.
type DatabaseResetJob struct {
	db                     database.IDatabase
	userRepo               user_repo.Querier
	categoryRepo           categories_repo.Querier
	paymentMethodRepo      payment_methods_repo.Querier
	cancellationReasonRepo cancellation_reasons_repo.Querier
	r2                     cloudflarer2.IR2
	logger                 logger.ILogger
	enabled                bool
	schedule               string
	pool                   *pgxpool.Pool
}

// NewDatabaseResetJob creates a new database maintenance cron job.
func NewDatabaseResetJob(
	db database.IDatabase,
	pool *pgxpool.Pool,
	userRepo user_repo.Querier,
	categoryRepo categories_repo.Querier,
	paymentMethodRepo payment_methods_repo.Querier,
	cancellationReasonRepo cancellation_reasons_repo.Querier,
	r2 cloudflarer2.IR2,
	logger logger.ILogger,
	enabled bool,
	schedule string,
) *DatabaseResetJob {
	return &DatabaseResetJob{
		db:                     db,
		pool:                   pool,
		userRepo:               userRepo,
		categoryRepo:           categoryRepo,
		paymentMethodRepo:      paymentMethodRepo,
		cancellationReasonRepo: cancellationReasonRepo,
		r2:                     r2,
		logger:                 logger,
		enabled:                enabled,
		schedule:               schedule,
	}
}

// GetSchedule returns the cron expression for when this job should run.
func (j *DatabaseResetJob) GetSchedule() string { return j.schedule }

// GetName returns the job name for logging and monitoring purposes.
func (j *DatabaseResetJob) GetName() string { return "DatabaseReset" }

// IsEnabled indicates whether this job should be registered and run.
func (j *DatabaseResetJob) IsEnabled() bool { return j.enabled }

// Execute performs a complete database reset and re-seed operation.
func (j *DatabaseResetJob) Execute(ctx context.Context) error {
	j.logger.Warn("Cron | Starting Scheduled Database Reset (WIPE and SEED)...")

	if err := j.db.ResetDatabase(ctx); err != nil {
		j.logger.Errorf("Cron | Database reset FAILED during Wipe phase: %v", err)
		return err
	}
	j.logger.Info("Cron | Database wipe completed successfully")

	if err := seeder.RunSeeders(
		ctx,
		j.pool,
		j.userRepo,
		j.categoryRepo,
		j.paymentMethodRepo,
		j.cancellationReasonRepo,
		j.r2,
		j.logger,
	); err != nil {
		j.logger.Errorf("Cron | Database reset FAILED during Seeding phase: %v", err)
		return err
	}

	j.logger.Info("Cron | Scheduled Database Reset completed successfully")
	return nil
}
