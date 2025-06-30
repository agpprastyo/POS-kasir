package logger

import (
	"POS-kasir/config"
	"github.com/sirupsen/logrus"
	"os"
)

// Logger extends logrus.Logger
type Logger struct {
	*logrus.Logger
}

// New creates a new configured logger
func New(cfg *config.AppConfig) *Logger {
	logger := &Logger{Logger: logrus.New()}

	// Set output
	if cfg.Logger.Output != nil {
		logger.SetOutput(cfg.Logger.Output)
	} else {
		logger.SetOutput(os.Stdout)
	}

	// Set formatter
	if cfg.Logger.JSONFormat {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  "2006-01-02T15:04:05.000Z07:00",
			DisableColors:    false,
			DisableTimestamp: false,
		})
	}

	// Set level
	level, err := logrus.ParseLevel(cfg.Logger.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	return logger
}

// Fields type to define log fields
type Fields logrus.Fields
