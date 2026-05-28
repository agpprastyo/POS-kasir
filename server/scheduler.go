// Package server provides scheduler and cron job management.
package server

import (
	"context"

	"github.com/robfig/cron/v3"
)

// CronJob defines the interface that all cron jobs must implement.
// This ensures consistency and makes it easy to add new jobs.
type CronJob interface {
	// GetSchedule returns the cron expression (e.g., "0 3 * * *" for 3 AM daily)
	GetSchedule() string

	// GetName returns a human-readable job name for logging purposes
	GetName() string

	// IsEnabled should return true if the job should be registered
	IsEnabled() bool

	// Execute runs the job logic with the given context
	Execute(ctx context.Context) error
}

// CronScheduler wraps the robfig/cron library and provides a clean, dependency-injected API
// for registering and managing cron jobs.
type CronScheduler struct {
	c   *cron.Cron
	app *App
}

// NewCronScheduler creates a new scheduler instance.
func NewCronScheduler(app *App) *CronScheduler {
	return &CronScheduler{
		c:   cron.New(),
		app: app,
	}
}

// RegisterJob registers a single cron job to the scheduler.
//
// If the job's IsEnabled() returns false, it is skipped with a log message.
// If registration fails, the error is logged and returned.
func (cs *CronScheduler) RegisterJob(job CronJob) error {
	if !job.IsEnabled() {
		cs.app.Logger.Infof("Cron | Job '%s' is disabled, skipping registration", job.GetName())
		return nil
	}

	_, err := cs.c.AddFunc(job.GetSchedule(), func() {
		cs.app.Logger.Infof("Cron | Starting job: %s", job.GetName())
		if err := job.Execute(context.Background()); err != nil {
			cs.app.Logger.Errorf("Cron | Job '%s' failed: %v", job.GetName(), err)
		} else {
			cs.app.Logger.Infof("Cron | Job '%s' completed successfully", job.GetName())
		}
	})

	if err != nil {
		cs.app.Logger.Errorf("Cron | Failed to register job '%s': %v", job.GetName(), err)
		return err
	}

	cs.app.Logger.Debugf("Cron | Job '%s' registered with schedule: %s", job.GetName(), job.GetSchedule())
	return nil
}

// RegisterJobs registers multiple jobs at once.
// Returns the count of successfully registered jobs.
func (cs *CronScheduler) RegisterJobs(jobs []CronJob) int {
	count := 0
	for _, job := range jobs {
		if err := cs.RegisterJob(job); err != nil {
			continue
		}
		count++
	}
	return count
}

// Start begins the cron scheduler.
func (cs *CronScheduler) Start() {
	cs.c.Start()
	cs.app.Logger.Info("Cron | Scheduler started")
}

// Stop gracefully stops the scheduler and waits for running jobs to complete.
func (cs *CronScheduler) Stop() {
	<-cs.c.Stop().Done()
	cs.app.Logger.Info("Cron | Scheduler stopped")
}
