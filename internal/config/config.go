package config

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	// Postgres
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	DBSSLMode  string

	// URA
	URAAccessKey string

	// Environment
	Env string
}

// Loads env variables
func Load() *Config {
	// reads .env file and sets them as env variables
	_ = godotenv.Load("../../.env")


	cfg := &Config{
		DBUser: getEnvOrFail("POSTGRES_USER"),
		DBPassword: getEnvOrFail("POSTGRES_PASSWORD"),
		DBName: getEnvOrFail("POSTGRES_DB"),
		DBHost: getEnvOrFail("POSTGRES_HOST"),
		DBPort: getEnvOrDefault("POSTGRES_PORT", "5433"),
		DBSSLMode: getEnvOrDefault("POSTGRES_SSLMODE", "disable"),
		URAAccessKey: getEnvOrFail("URA_ACCESS_KEY"),
		Env: getEnvOrDefault("ENV", "development"),
	}

	log.Printf("[config] Loaded configuration for envrionment: %s", cfg.Env)
	return cfg
}

// Returns PostgreSQL data source name (connection string)
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}