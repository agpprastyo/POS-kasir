package logger

import (
	"POS-kasir/config"
	"os"

	"github.com/sirupsen/logrus"
)

type ILogger = logrus.FieldLogger

type Fields = logrus.Fields

func New(cfg *config.AppConfig) ILogger {
	base := logrus.New()

	if cfg != nil && cfg.Logger.Output != nil {
		base.SetOutput(cfg.Logger.Output)
	} else {
		base.SetOutput(os.Stdout)
	}

	if cfg != nil && cfg.Logger.JSONFormat {
		base.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	} else {
		base.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
			DisableQuote:    true,
		})
	}

	// Level
	var lvl logrus.Level
	if cfg != nil {
		if parsed, err := logrus.ParseLevel(cfg.Logger.Level); err == nil {
			lvl = parsed
		} else {
			lvl = logrus.InfoLevel
		}
	} else {
		lvl = logrus.InfoLevel
	}
	base.SetLevel(lvl)

	if cfg != nil {
		return base.WithFields(logrus.Fields{
			"app": cfg.Server.AppName,
			"env": cfg.Server.Env,
		})
	}

	return base.WithFields(logrus.Fields{})
}
