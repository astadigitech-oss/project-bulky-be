package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppEnv                        string
	AppPort                       string
	DBHost                        string
	DBPort                        string
	DBUser                        string
	DBPassword                    string
	DBName                        string
	JWTSecret                     string
	JWTAccessDuration             time.Duration
	BcryptCost                    int
	UploadPath                    string
	BaseURL                       string
	InternalAPIKey                string
	WMSAPIKey                     string
	DelivereeWebhookAuthorization string
	ForwarderWebhookAuthorization string
	WMSBaseURL                    string
	WMSClientID                   string
	WMSClientSecret               string
	StorefrontBaseURL             string
	BackupPath                    string
	BackupRetentionDays           int
	BackupScheduleHour            int
	BackupScheduleMinute          int
}

func LoadConfig() *Config {
	// Parse JWT access duration (24 hours for single token)
	accessDuration := parseDuration(getEnv("JWT_ACCESS_EXPIRY", "24h"), 24*time.Hour)

	// Get bcrypt cost based on environment
	bcryptCost := getBcryptCost(getEnv("APP_ENV", "development"))

	uploadPath := getEnv("UPLOAD_PATH", "./uploads")
	backupPath := getEnv("BACKUP_PATH", "")
	if backupPath == "" {
		backupPath = uploadPath + "/backups"
	}

	return &Config{
		AppEnv:                        getEnv("APP_ENV", "development"),
		AppPort:                       getEnv("APP_PORT", "8080"),
		DBHost:                        getEnv("DB_HOST", "localhost"),
		DBPort:                        getEnv("DB_PORT", "5432"),
		DBUser:                        getEnv("DB_USER", "postgres"),
		DBPassword:                    getEnv("DB_PASSWORD", ""),
		DBName:                        getEnv("DB_NAME", "bulky_db"),
		JWTSecret:                     getEnv("JWT_SECRET", "your-secret-key-minimum-32-characters-long"),
		JWTAccessDuration:             accessDuration,
		BcryptCost:                    bcryptCost,
		UploadPath:                    uploadPath,
		BaseURL:                       getEnv("BASE_URL", "http://localhost:8080"),
		InternalAPIKey:                getEnv("INTERNAL_API_KEY", ""),
		WMSAPIKey:                     getEnv("WMS_API_KEY", ""),
		DelivereeWebhookAuthorization: getEnv("DELIVEREE_WEBHOOK_AUTHORIZATION", ""),
		ForwarderWebhookAuthorization: getEnv("FORWARDER_WEBHOOK_AUTHORIZATION", ""),
		WMSBaseURL:                    getEnv("WMS_BASE_URL", ""),
		WMSClientID:                   getEnv("WMS_CLIENT_ID", ""),
		WMSClientSecret:               getEnv("WMS_CLIENT_SECRET", ""),
		StorefrontBaseURL:             getEnv("STOREFRONT_BASE_URL", ""),
		BackupPath:                    backupPath,
		BackupRetentionDays:           getEnvInt("BACKUP_RETENTION_DAYS", 7),
		BackupScheduleHour:            getEnvInt("BACKUP_SCHEDULE_HOUR_UTC", 17), // 17:00 UTC = 00:00 WIB
		BackupScheduleMinute:          getEnvInt("BACKUP_SCHEDULE_MINUTE_UTC", 0),
	}
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}
	return defaultValue
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

// getBcryptCost returns the bcrypt cost based on environment
func getBcryptCost(env string) int {
	// Allow override via environment variable
	if costStr := os.Getenv("BCRYPT_COST"); costStr != "" {
		if cost, err := strconv.Atoi(costStr); err == nil && cost >= 4 && cost <= 31 {
			return cost
		}
	}

	// Default based on environment
	switch env {
	case "production":
		return 12 // ~400ms
	case "staging":
		return 11 // ~200ms
	default: // development, testing
		return 10 // ~100ms
	}
}
