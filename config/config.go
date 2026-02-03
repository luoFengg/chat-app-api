package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Main Struct
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

// AppConfig
type AppConfig struct {
	Name string
	Port string
	Env  string
}

// DatabaseConfig for PostgreSQL
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// RedisConfig for Redis
type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

// JWTConfig for JWT
type JWTConfig struct {
	Secret string
	RefreshSecret string
	AccessExpiryMins int
	RefreshExpiryDays int
}

// LoadConfig to read .env dan return Config struct
func LoadConfig() *Config {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Parse JWT access token expiry (minutes)
	accessExpiry, err := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRY_MINUTES", "15"))
	if err != nil {
		accessExpiry = 15
	}

	// Parse JWT refresh token expiry (days)
	refreshExpiry, err := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRY_DAYS", "30"))
	if err != nil {
		refreshExpiry = 30
	}

	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "ChatApp API"),
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "chatapp"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "secret"),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", "refresh_secret"),
			AccessExpiryMins: accessExpiry,
			RefreshExpiryDays: refreshExpiry,
		},
	}
}

// getEnv to get environment variable with default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}