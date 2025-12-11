package config

import (
	"os"
	"strconv"
	"time"
)

func getBool(key string, fallback bool) bool {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := strconv.ParseBool(strValue)
	if err != nil {
		return fallback
	}
	return value
}

// Helper functions for environment variables
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getInt(key string, fallback int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := strconv.Atoi(strValue)
	if err != nil {
		return fallback
	}
	return value
}

// getDuration (kept for convenience; prefer explicit getInt + unit in config)
func getDuration(key string, fallback time.Duration) time.Duration {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := time.ParseDuration(strValue)
	if err != nil {
		return fallback
	}
	return value
}

// getInt64
func getInt64(key string, fallback int64) int64 {
	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	value, err := strconv.ParseInt(strValue, 10, 64)
	if err != nil {
		return fallback
	}
	return value
}

// getEnvEnum returns env value only if it's in validValues; otherwise returns fallback.
// It also ensures fallback is valid (if not, the first validValues entry will be used).
func getEnvEnum(key string, validValues []string, fallback string) string {
	// ensure fallback is valid
	validFallback := false
	for _, v := range validValues {
		if v == fallback {
			validFallback = true
			break
		}
	}
	if !validFallback && len(validValues) > 0 {
		fallback = validValues[0]
	}

	strValue := getEnv(key, "")
	if strValue == "" {
		return fallback
	}
	for _, value := range validValues {
		if strValue == value {
			return strValue
		}
	}
	return fallback
}
