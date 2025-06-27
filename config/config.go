package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all configuration values
type Config struct {
	Token    string `env:"TOKEN"`
	Timezone string `env:"TIMEZONE"`
}

// Load loads configuration from environment variables and .env file
func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}

	config := &Config{
		Token:    getEnv("TOKEN", ""),
		Timezone: getEnv("TIMEZONE", "Europe/Madrid"),
	}

	return config
}

// getEnv gets an environment variable with a fallback value
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
