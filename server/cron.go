package server

import (
	"context"

	"github.com/robfig/cron/v3"
)

func SetupCron(app *App, container *AppContainer) {
	c := cron.New()

	// Auto-close shift at 03:00 every day
	_, err := c.AddFunc("0 3 * * *", func() {
		app.Logger.Info("Cron | Starting auto-close shifts job...")
		err := container.ShiftService.AutoCloseShifts(context.Background())
		if err != nil {
			app.Logger.Errorf("Cron | Auto-close shifts job failed: %v", err)
		} else {
			app.Logger.Info("Cron | Auto-close shifts job completed successfully")
		}
	})

	if err != nil {
		app.Logger.Errorf("Failed to setup cron: %v", err)
		return
	}

	c.Start()
	app.Logger.Info("Cron | Scheduler started")
}
