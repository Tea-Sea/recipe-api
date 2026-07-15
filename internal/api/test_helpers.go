package api

import (
	"fmt"
	"recipe-api/internal/models"
	"recipe-api/internal/repository"
	"testing"
	"time"

	"gorm.io/gorm"
)

func ToPtr[T any](v T) *T {
	return &v
}

func createTestUnit(label string) *models.Unit {
	return &models.Unit{
		Label: label,
	}
}

func createTestIngredient(label string) *models.Ingredient {
	return &models.Ingredient{
		Label: label,
	}
}

func createTestRecipe(t *testing.T, app *App, uniqueID ...bool) models.Recipe {
	t.Helper()

	idRequired := false
	if len(uniqueID) > 0 {
		idRequired = uniqueID[0]
	}

	name := "Test Recipe"
	if idRequired {
		name = fmt.Sprintf(name+" %d", time.Now().UnixNano())
	}

	recipe := models.Recipe{
		Name:       name,
		Difficulty: 5,
		Ingredients: []models.RecipeIngredient{
			{
				Amount:     ToPtr(float32(1)),
				Ingredient: createTestIngredient("Salt"),
				Unit:       createTestUnit("Cup"),
			},
			{
				Amount:     ToPtr(float32(4)),
				Ingredient: createTestIngredient("Water"),
				Unit:       createTestUnit("Cup"),
			},
		},
		Instructions: []models.Instruction{
			{
				StepNumber: 1,
				StepText:   "boil it",
			},
			{
				StepNumber: 2,
				StepText:   "dry it",
			},
		},
	}
	return recipe
}

func createRepository(db *gorm.DB) *repository.App {
	return &repository.App{
		DB: db,
	}
}

func clearDatabase(app *App) {
	app.Repo.DB.Exec("DELETE FROM recipe_ingredients")
	app.Repo.DB.Exec("DELETE FROM recipes")
}
