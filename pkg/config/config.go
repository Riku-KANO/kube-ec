package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	ServerPort  int
	GRPCPort    int
	DatabaseURL string
	JWTSecret   string
	Environment string
}

func Load() (*Config, error) {
	cfg := &Config{
		ServerPort:  getEnvAsInt("SERVER_PORT", 8080),
		GRPCPort:    getEnvAsInt("GRPC_PORT", 50051),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
		Environment: getEnv("ENVIRONMENT", "development"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}
	if c.JWTSecret == "" && c.Environment == "production" {
		return fmt.Errorf("JWT_SECRET is required in production")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
