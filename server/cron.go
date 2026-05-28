package server

// SetupCron initializes the cron scheduler and registers all configured jobs.
//
// This function is the orchestration point for all scheduled tasks. It:
// 1. Creates a new CronScheduler instance
// 2. Gets all jobs from the JobsRegistry
// 3. Registers each job to the scheduler
// 4. Starts the scheduler
//
// Jobs are defined using the CronJob interface and located in their respective
// modules (e.g., internal/shift/job.go) or in server/jobs_*.go files for infrastructure jobs.
// The JobsRegistry acts as a central factory for creating all jobs.
func SetupCron(app *App, container *AppContainer) {
	scheduler := NewCronScheduler(app)

	// Get all configured jobs from registry
	registry := NewJobsRegistry(app, container)
	jobs := registry.RegisterAll()

	if len(jobs) == 0 {
		app.Logger.Warn("Cron | No jobs registered")
		return
	}

	// Register all jobs to the scheduler
	registered := scheduler.RegisterJobs(jobs)
	app.Logger.Infof("Cron | Registered %d out of %d jobs", registered, len(jobs))

	scheduler.Start()
}
