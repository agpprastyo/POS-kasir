package logger

import (
	"POS-kasir/config"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/getsentry/sentry-go"
)

// ILogger defines the logging interface used throughout the app.
// We maintain the logrus-style methods for backward compatibility during the slog migration.
type ILogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Print(args ...interface{})
	Printf(format string, args ...interface{})
	Println(args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	WithFields(fields Fields) ILogger
}

type Fields map[string]interface{}

type slogLogger struct {
	inner *slog.Logger
}

func New(cfg *config.AppConfig) ILogger {
	var level slog.Level
	if cfg != nil {
		switch strings.ToLower(cfg.Logger.Level) {
		case "debug":
			level = slog.LevelDebug
		case "info":
			level = slog.LevelInfo
		case "warn", "warning":
			level = slog.LevelWarn
		case "error":
			level = slog.LevelError
		default:
			level = slog.LevelInfo
		}
	} else {
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var output io.Writer = os.Stdout
	if cfg != nil && cfg.Logger.Output != nil {
		output = cfg.Logger.Output
	}

	var handler slog.Handler
	if cfg != nil && cfg.Logger.JSONFormat {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}

	logger := slog.New(handler)

	if cfg != nil && cfg.Logger.SentryDSN != "" {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              cfg.Logger.SentryDSN,
			Environment:      cfg.Server.Env,
			Release:          cfg.Server.AppName,
			Debug:            cfg.Server.Env == "development",
			TracesSampleRate: 1.0,
		})
		if err != nil {
			logger.Error("Failed to initialize Sentry", slog.Any("error", err))
		} else {
			logger.Info("Sentry initialized successfully")
		}
	}

	return &slogLogger{inner: logger}
}

func (l *slogLogger) Debug(args ...interface{}) {
	l.inner.Debug(fmt.Sprint(args...))
}

func (l *slogLogger) Debugf(format string, args ...interface{}) {
	l.inner.Debug(fmt.Sprintf(format, args...))
}

func (l *slogLogger) Info(args ...interface{}) {
	l.inner.Info(fmt.Sprint(args...))
}

func (l *slogLogger) Infof(format string, args ...interface{}) {
	l.inner.Info(fmt.Sprintf(format, args...))
}

func (l *slogLogger) Warn(args ...interface{}) {
	l.inner.Warn(fmt.Sprint(args...))
}

func (l *slogLogger) Warnf(format string, args ...interface{}) {
	l.inner.Warn(fmt.Sprintf(format, args...))
}

func (l *slogLogger) Error(args ...interface{}) {
	l.inner.Error(fmt.Sprint(args...))
}

func (l *slogLogger) Errorf(format string, args ...interface{}) {
	l.inner.Error(fmt.Sprintf(format, args...))
}

func (l *slogLogger) Print(args ...interface{}) {
	l.inner.Info(fmt.Sprint(args...))
}

func (l *slogLogger) Printf(format string, args ...interface{}) {
	l.inner.Info(fmt.Sprintf(format, args...))
}

func (l *slogLogger) Println(args ...interface{}) {
	l.inner.Info(fmt.Sprint(args...))
}

func (l *slogLogger) Fatal(args ...interface{}) {
	l.inner.Error(fmt.Sprint(args...))
	os.Exit(1)
}

func (l *slogLogger) Fatalf(format string, args ...interface{}) {
	l.inner.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *slogLogger) Panic(args ...interface{}) {
	msg := fmt.Sprint(args...)
	l.inner.Error(msg)
	panic(msg)
}

func (l *slogLogger) Panicf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	l.inner.Error(msg)
	panic(msg)
}

func (l *slogLogger) WithFields(fields Fields) ILogger {
	var attrs []interface{}
	for k, v := range fields {
		attrs = append(attrs, slog.Any(k, v))
	}
	return &slogLogger{inner: l.inner.With(attrs...)}
}
