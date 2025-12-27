package config

import "os"

type Config struct {
	AppEnv     string
	AppPort    string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	UploadPath string
	BaseURL    string
}

func LoadConfig() *Config {
	return &Config{
		AppEnv:     getEnv("APP_ENV", "development"),
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "bulky_db"),
		JWTSecret:  getEnv("JWT_SECRET", "your-secret-key"),
		UploadPath: getEnv("UPLOAD_PATH", "./uploads"),
		BaseURL:    getEnv("BASE_URL", "http://localhost:8080/uploads"),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
