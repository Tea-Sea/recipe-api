package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
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
	LoadEnvFromRoot()

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

func LoadEnvFromRoot() {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	root := filepath.Join(dir, "../../")

	envPath := filepath.Join(root, ".env")

	err := godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: .env file not found at: %v. Will refer to default values instead. \n", envPath)
	}
}
