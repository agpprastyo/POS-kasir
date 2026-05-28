package server

import (
	"POS-kasir/internal/shift"
	"POS-kasir/pkg/database/maintenance"
)

// JobsRegistry provides a centralized factory for creating and registering all cron jobs.

type JobsRegistry struct {
	app       *App
	container *AppContainer
}

// NewJobsRegistry creates a new registry instance.
func NewJobsRegistry(app *App, container *AppContainer) *JobsRegistry {
	return &JobsRegistry{
		app:       app,
		container: container,
	}
}

// RegisterAll returns a list of all cron jobs that should be registered.
//
// This method is called once during application startup to build the complete
// list of jobs to be scheduled. Each job is created with its required dependencies.
func (r *JobsRegistry) RegisterAll() []CronJob {
	var jobs []CronJob

	// ==== SHIFT MODULE JOBS ====
	jobs = append(jobs, shift.NewShiftAutoCloseJob(
		r.container.ShiftService,
		r.app.Logger,
	))

	// ==== DATABASE JOBS ====
	// Only add database reset job if it's enabled in config
	if r.app.Config.EnableDbWipe {
		jobs = append(jobs, maintenance.NewDatabaseResetJob(
			r.app.DB,
			r.app.DB.GetPool(),
			r.container.UserRepo,
			r.container.CategoryRepo,
			r.container.PaymentMethodRepo,
			r.container.CancellationReasonRepo,
			r.app.R2,
			r.app.Logger,
			true, // enabled
			r.app.Config.WipeCronSchedule,
		))
	}

	// ==== ADD NEW JOBS BELOW THIS LINE ====
	// Example:
	// jobs = append(jobs, report.NewDailyReportJob(
	//     r.container.ReportService,
	//     r.app.Logger,
	// ))
	//
	// jobs = append(jobs, promotions.NewCleanupExpiredJob(
	//     r.container.PromotionService,
	//     r.app.Logger,
	// ))

	return jobs
}
