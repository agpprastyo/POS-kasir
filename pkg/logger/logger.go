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

// In pkg/logger/logger.go
type ILogger interface {
	Errorf(format string, args ...interface{})
	Error(args ...interface{})
	Infof(format string, args ...interface{})
	Info(args ...interface{})

	Warnf(format string, args ...interface{})
	Warn(args ...interface{})
	Println(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Fatalln(args ...interface{})
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.Logger.Error(args...)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}
func (l *Logger) Info(args ...interface{}) {
	l.Logger.Info(args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.Logger.Warn(args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.Logger.Println(args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.Logger.Fatalln(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Logger.Fatal(args...)
}

// New creates a new configured logger
func New(cfg *config.AppConfig) ILogger {
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
