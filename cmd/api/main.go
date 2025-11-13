package main

import (
	"fmt"
	"log"
	"net/http"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"recipe-api/internal/api"
	"recipe-api/internal/config"
	"recipe-api/internal/logger"
	"recipe-api/internal/repository"
)

func init() {
	fmt.Println("App is starting…")
}

func main() {
	cfg := config.Load()
	log.Println("Starting server on port", cfg.Port)

	appLogger := logger.NewAppLogger("[recipes]", cfg.Debug)
	appLogger.Println("Starting application...")

	gormLogger := logger.NewGormLogger()

	dbConn, err := gorm.Open(
		postgres.Open(cfg.DatabaseURL),
		&gorm.Config{
			Logger: gormLogger,
		},
	)
	if err != nil {
		appLogger.Fatal("Failed to connect to database: ", err)
	}

	repoApp := repository.NewApp(dbConn)

	err = repoApp.AutoMigrate()
	if err != nil {
		appLogger.Fatal("failed to run database migrations:", err)
	}

	apiApp := &api.App{
		Repo:   repoApp,
		Logger: appLogger,
	}

	// Get the underlying *sql.DB for connection pooling configuration
	sqlDB, err := repoApp.DB.DB()
	if err != nil {
		appLogger.Fatalf("failed to get generic database object: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		appLogger.Fatal("database ping failed:", err)
	}

	appLogger.Println("Database connection is alive.")

	router := api.NewRouter(apiApp, cfg.FrontendURL, appLogger)

	appLogger.Fatal(http.ListenAndServe(":"+cfg.Port, router))

	appLogger.Printf("Server listening on port: %v \n", cfg.Port)

	appLogger.Println("App initialized successfully.")
}
