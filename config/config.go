package config

import (
	"io"
	"os"
	"time"
)

type AppConfig struct {
	Server serverConfig
	DB     dbConfig
	Logger loggerConfig
}

type serverConfig struct {
	AppName string
	Env     string
	Port    string
}

type loggerConfig struct {
	Level      string
	JSONFormat bool
	Output     io.Writer
}

type dbConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	DBName      string
	SSLMode     string
	MaxOpenConn int
	MaxIdleConn int
	MaxLifetime time.Duration
}

func Load() *AppConfig {
	return &AppConfig{
		DB: dbConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5432"),
			User:        getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", "postgres"),
			DBName:      getEnv("DB_NAME", "agprastyowsl"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			MaxOpenConn: getInt("DB_MAX_OPEN_CONNECTIONS", 10),
			MaxIdleConn: getInt("DB_MAX_IDLE_CONNECTIONS", 2),
			MaxLifetime: time.Duration(getInt("DB_MAX_LIFETIME", 10)),
		},
		Logger: loggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			JSONFormat: getBool("LOG_JSON_FORMAT", true),
			Output:     os.Stdout,
		},
		Server: serverConfig{
			AppName: getEnv("APP_NAME", "agprastyowsl"),
			Env:     getEnv("APP_ENV", "production"),
			Port:    getEnv("APP_PORT", "8080"),
		},
	}
}
