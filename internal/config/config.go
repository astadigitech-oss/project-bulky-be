package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv             string
	AppPort            string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	JWTAccessDuration  time.Duration
	JWTRefreshDuration time.Duration
	UploadPath         string
	BaseURL            string
	PlayStoreURL       string
	AppStoreURL        string
}

func LoadConfig() *Config {
	// Parse JWT durations
	accessDuration := parseDuration(getEnv("JWT_ACCESS_EXPIRY", "1h"), time.Hour)
	refreshDuration := parseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"), 168*time.Hour) // 7 days

	return &Config{
		AppEnv:             getEnv("APP_ENV", "development"),
		AppPort:            getEnv("APP_PORT", "8080"),
		DBHost:             getEnv("DB_HOST", "localhost"),
		DBPort:             getEnv("DB_PORT", "5432"),
		DBUser:             getEnv("DB_USER", "postgres"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "bulky_db"),
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key-minimum-32-characters-long"),
		JWTAccessDuration:  accessDuration,
		JWTRefreshDuration: refreshDuration,
		UploadPath:         getEnv("UPLOAD_PATH", "./uploads"),
		BaseURL:            getEnv("BASE_URL", "http://localhost:8080/uploads"),
		PlayStoreURL:       getEnv("PLAY_STORE_URL", "https://play.google.com/store/apps/details?id=com.bulky"),
		AppStoreURL:        getEnv("APP_STORE_URL", "https://apps.apple.com/app/bulky/id123456789"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseDuration(durationStr string, defaultDuration time.Duration) time.Duration {
	// Try to parse as duration string (e.g., "1h", "30m", "168h")
	if duration, err := time.ParseDuration(durationStr); err == nil {
		return duration
	}

	// Try to parse as seconds (for backward compatibility)
	if seconds, err := strconv.Atoi(durationStr); err == nil {
		return time.Duration(seconds) * time.Second
	}

	return defaultDuration
}
