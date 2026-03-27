package server

import (
	"POS-kasir/pkg/database/seeder"
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
		app.Logger.Errorf("Failed to setup shift auto-close cron: %v", err)
	}

	// Daily Database Reset (for portfolio demo consistency)
	if app.Config.EnableDbWipe {
		_, err = c.AddFunc(app.Config.WipeCronSchedule, func() {
			app.Logger.Warn("Cron | Starting Scheduled Database Reset (WIPE and SEED)...")
			ctx := context.Background()

			// 1. Wipe
			if err := app.DB.ResetDatabase(ctx); err != nil {
				app.Logger.Errorf("Cron | Database reset FAILED during Wipe: %v", err)
				return
			}

			// 2. Re-seed
			if err := seeder.RunSeeders(
				ctx,
				app.DB.GetPool(),
				container.UserRepo,
				container.CategoryRepo,
				container.PaymentMethodRepo,
				container.CancellationReasonRepo,
				app.R2,
				app.Logger,
			); err != nil {
				app.Logger.Errorf("Cron | Database reset FAILED during Seeding: %v", err)
				return
			}

			app.Logger.Info("Cron | Scheduled Database Reset completed successfully")
		})

		if err != nil {
			app.Logger.Errorf("Failed to setup database reset cron: %v", err)
		} else {
			app.Logger.Infof("Cron | Database reset scheduled at: %s", app.Config.WipeCronSchedule)
		}
	}

	c.Start()
	app.Logger.Info("Cron | Scheduler started")
}
