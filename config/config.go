package config

import (
	"fmt"
	"io"
	"os"
	"time"
)

type AppConfig struct {
	Server           ServerConfig
	DB               DbConfig
	Logger           LoggerConfig
	JWT              JwtConfig
	CloudflareR2     CloudflareR2Config
	Midtrans         MidtransConfig
	Redis            RedisConfig
	Metrics          MetricsConfig
	Alert            AlertConfig
	AutoMigrate      bool
	MigrationsPath   string
	EnableDbWipe     bool
	WipeCronSchedule string
}

type MetricsConfig struct {
	Enabled bool
	Path    string
}

type AlertConfig struct {
	SlackEnabled    bool
	SlackWebhookURL string
	SlackChannel    string

	DiscordEnabled    bool
	DiscordWebhookURL string

	TelegramEnabled  bool
	TelegramBotToken string
	TelegramChatID   string

	EmailEnabled  bool
	EmailSMTPHost string
	EmailSMTPPort int
	EmailFrom     string
	EmailTo       string
	EmailUsername string
	EmailPassword string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type CloudflareR2Config struct {
	AccountID    string
	AccessKey    string
	SecretKey    string
	Bucket       string
	PublicDomain string
	ExpirySec    int64
	Endpoint     string
	UseSSL       bool
}

type MidtransConfig struct {
	ServerKey string `mapstructure:"server_key"`
	IsProd    bool   `mapstructure:"is_prod"`
}

type JwtConfig struct {
	Secret               string
	Duration             time.Duration
	RefreshTokenDuration time.Duration
	Issuer               string
}

type ServerConfig struct {
	AppName                string
	Env                    string
	Port                   string
	CookieDomain           string
	WebFrontendCrossOrigin bool
	CorsAllowOrigins       string
}

type LoggerConfig struct {
	Level      string
	JSONFormat bool
	SentryDSN  string
	Output     io.Writer
}

type DbConfig struct {
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

func (d DbConfig) BuildDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode)
}

func Load() *AppConfig {
	maxLifetimeMinutes := getInt("DB_MAX_LIFETIME_MINUTES", 10)
	return &AppConfig{
		Midtrans: MidtransConfig{
			ServerKey: getEnv("MIDTRANS_SERVER_KEY", "SB-Mid-server-1234567890"),
			IsProd:    getBool("MIDTRANS_IS_PROD", false),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getInt("REDIS_DB", 0),
		},
		DB: DbConfig{
			Host:        getEnv("DB_HOST", "localhost"),
			Port:        getEnv("DB_PORT", "5432"),
			User:        getEnv("DB_USER", "postgres"),
			Password:    getEnv("DB_PASSWORD", "postgres"),
			DBName:      getEnv("DB_NAME", "agprastyowsl"),
			SSLMode:     getEnv("DB_SSLMODE", "disable"),
			MaxOpenConn: getInt("DB_MAX_OPEN_CONNECTIONS", 10),
			MaxIdleConn: getInt("DB_MAX_IDLE_CONNECTIONS", 2),
			MaxLifetime: time.Duration(maxLifetimeMinutes) * time.Minute,
		},
		Logger: LoggerConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			JSONFormat: getBool("LOG_JSON_FORMAT", true),
			SentryDSN:  getEnv("SENTRY_DSN", ""),
			Output:     os.Stdout,
		},
		Server: ServerConfig{
			AppName:                getEnv("APP_NAME", "hmm"),
			Env:                    getEnvEnum("APP_ENV", []string{"production", "development"}, "production"),
			Port:                   getEnv("APP_PORT", "8080"),
			CookieDomain:           getEnv("COOKIE_DOMAIN", ""),
			WebFrontendCrossOrigin: getBool("WEB_FRONTEND_CROSS_ORIGIN", false),
			CorsAllowOrigins:       getEnv("CORS_ALLOW_ORIGINS", ""),
		},
		JWT: JwtConfig{
			Secret:               getEnv("JWT_SECRET", "secret"),
			Duration:             time.Duration(getInt("JWT_DURATION_HOURS", 24)) * time.Hour,
			RefreshTokenDuration: time.Duration(getInt("JWT_REFRESH_DURATION_DAYS", 7)) * 24 * time.Hour,
			Issuer:               getEnv("JWT_ISSUER", "poskasir"),
		},
		CloudflareR2: CloudflareR2Config{
			AccountID:    getEnv("R2_ACCOUNT_ID", ""),
			AccessKey:    getEnv("R2_ACCESS_KEY", ""),
			SecretKey:    getEnv("R2_SECRET_KEY", ""),
			Bucket:       getEnv("R2_BUCKET", "pos-kasir"),
			PublicDomain: getEnv("R2_PUBLIC_DOMAIN", ""),
			ExpirySec:    getInt64("R2_EXPIRY_SECONDS", 3600),
			Endpoint:     getEnv("R2_ENDPOINT", ""),
			UseSSL:       getBool("R2_USE_SSL", true),
		},
		Metrics: MetricsConfig{
			Enabled: getBool("METRICS_ENABLED", true),
			Path:    getEnv("METRICS_PATH", "/metrics"),
		},
		Alert: AlertConfig{
			SlackEnabled:    getBool("ALERT_SLACK_ENABLED", false),
			SlackWebhookURL: getEnv("ALERT_SLACK_WEBHOOK_URL", ""),
			SlackChannel:    getEnv("ALERT_SLACK_CHANNEL", "#pos-alerts"),

			DiscordEnabled:    getBool("ALERT_DISCORD_ENABLED", false),
			DiscordWebhookURL: getEnv("ALERT_DISCORD_WEBHOOK_URL", ""),

			TelegramEnabled:  getBool("ALERT_TELEGRAM_ENABLED", false),
			TelegramBotToken: getEnv("ALERT_TELEGRAM_BOT_TOKEN", ""),
			TelegramChatID:   getEnv("ALERT_TELEGRAM_CHAT_ID", ""),

			EmailEnabled:  getBool("ALERT_EMAIL_ENABLED", false),
			EmailSMTPHost: getEnv("ALERT_EMAIL_SMTP_HOST", "smtp.gmail.com"),
			EmailSMTPPort: getInt("ALERT_EMAIL_SMTP_PORT", 587),
			EmailFrom:     getEnv("ALERT_EMAIL_FROM", ""),
			EmailTo:       getEnv("ALERT_EMAIL_TO", ""),
			EmailUsername: getEnv("ALERT_EMAIL_USERNAME", ""),
			EmailPassword: getEnv("ALERT_EMAIL_PASSWORD", ""),
		},
		AutoMigrate:      getBool("AUTO_MIGRATE", false),
		MigrationsPath:   getEnv("MIGRATIONS_PATH", "file://./sqlc/migrations"),
		EnableDbWipe:     getBool("ENABLE_DB_WIPE", false),
		WipeCronSchedule: getEnv("WIPE_CRON_SCHEDULE", "0 1 * * *"),
	}
}
