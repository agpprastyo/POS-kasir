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
	JWT    jwtConfig
	Minio  minioConfig
}

type minioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
}

type jwtConfig struct {
	Secret   string
	Duration time.Duration
	Issuer   string
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
		JWT: jwtConfig{
			Secret:   getEnv("JWT_SECRET", "secret"),
			Duration: time.Duration(getInt("JWT_DURATION_HOURS", 24)) * time.Hour,
			Issuer:   getEnv("JWT_ISSUER", "agprastyo"),
		},
		Minio: minioConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin123"),
			UseSSL:    getBool("MINIO_USE_SSL", false),
			Bucket:    getEnv("MINIO_BUCKET", "pos-kasir"),
		},
	}
}
