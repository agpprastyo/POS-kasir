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

func getEnvEnum(key string, validValues []string, fallback string) string {
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
