package shift

import (
	"POS-kasir/pkg/logger"
	"context"
)

// ShiftAutoCloseJob implements the CronJob interface for automatically closing open shifts.
//
// This job is typically scheduled to run daily (e.g., 03:00 AM) to ensure no shifts
// are left open indefinitely. It gracefully handles errors and logs them without failing
// the entire cron run.
type ShiftAutoCloseJob struct {
	service Service
	logger  logger.ILogger
}

// NewShiftAutoCloseJob creates a new instance of ShiftAutoCloseJob.
//
// Parameters:
//   - service: The shift service that implements the AutoCloseShifts logic
//   - logger: Logger for job execution logs
func NewShiftAutoCloseJob(service Service, logger logger.ILogger) *ShiftAutoCloseJob {
	return &ShiftAutoCloseJob{
		service: service,
		logger:  logger,
	}
}

// GetSchedule returns the cron expression for when this job should run.
// "0 3 * * *" means 03:00 AM every day.
func (j *ShiftAutoCloseJob) GetSchedule() string {
	return "0 3 * * *"
}

// GetName returns the job name for logging and monitoring purposes.
func (j *ShiftAutoCloseJob) GetName() string {
	return "ShiftAutoClose"
}

// IsEnabled indicates whether this job should be registered and run.
// Currently always enabled, but can be made conditional if needed.
func (j *ShiftAutoCloseJob) IsEnabled() bool {
	return true
}

// Execute runs the auto-close shifts logic by calling the shift service.
//
// This method is called by the cron scheduler at the scheduled time. Any errors
// encountered are returned to the scheduler for logging.
func (j *ShiftAutoCloseJob) Execute(ctx context.Context) error {
	return j.service.AutoCloseShifts(ctx)
}
