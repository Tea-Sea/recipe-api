package api

import (
	"os"
	"recipe-api/internal/logger"
	"recipe-api/internal/models"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var testApp *App

func TestMain(m *testing.M) {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		os.Exit(1)
	}

	err = db.AutoMigrate(
		&models.Recipe{},
		&models.Ingredient{},
		&models.Unit{},
		&models.RecipeIngredient{},
		&models.Instruction{},
	)

	repo := createRepository(db)

	testApp = &App{
		Repo:   repo,
		Logger: logger.NewAppLogger("TEST", true),
	}
	exitCode := m.Run()
	os.Exit(exitCode)
}
