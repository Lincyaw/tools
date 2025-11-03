package config

import (
	"log"
	"os"
	"strconv"
)

// Config application configuration
type Config struct {
	Environment string
	Port        string
	BaseURL     string
	Database    DatabaseConfig
	Redis       RedisConfig
}

// DatabaseConfig database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// RedisConfig Redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// Load load configuration
func Load() *Config {
	cfg := &Config{
		Environment: getEnv("APP_ENV", "development"),
		Port:        getEnv("APP_PORT", "8080"),
		BaseURL:     getEnv("BASE_URL", "http://localhost:8080"),
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "tools"),
			Password: getEnv("DB_PASSWORD", "tools123"),
			DBName:   getEnv("DB_NAME", "tools"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
	}

	log.Printf("Configuration loaded: env=%s, port=%s", cfg.Environment, cfg.Port)
	return cfg
}

// getEnv get environment variable, return default if not exists
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt get environment variable and convert to integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: invalid integer value for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return value
}
