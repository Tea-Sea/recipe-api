package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	FrontendURL string
	Debug       bool
}

func Load() *Config {
	var err = godotenv.Load("../../.env")
	if err != nil {
		log.Println("Error loading .env file. Refering to default values.")
	}

	debug, _ := strconv.ParseBool(os.Getenv("DEBUG"))

	cfg := &Config{
		Port:        loadEnv("PORT", "8080"),
		DatabaseURL: loadEnv("DATABASE_URL", ""),
		FrontendURL: loadEnv("FRONTEND_URL", ""),
		Debug:       debug,
	}
	return cfg
}

// Load variable from Env file, or default values if they don't exist
func loadEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
