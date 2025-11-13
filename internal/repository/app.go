package repository

import (
	"recipe-api/internal/models"

	"gorm.io/gorm"
)

// Applicaton struct to prevent use of globals
type App struct {
	DB *gorm.DB
}

// Constructor
func NewApp(db *gorm.DB) *App {
	return &App{DB: db}
}

func (app *App) AutoMigrate() error {
	return app.DB.AutoMigrate(&models.Recipe{})
}
