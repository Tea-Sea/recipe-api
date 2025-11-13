package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"

	"recipe-api/internal/models"
)

// Get all recipes
func (app *App) getAllRecipes(w http.ResponseWriter, r *http.Request) {
	var recipes []models.Recipe

	result := app.Repo.DB.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("ingredient_id ASC")
	}).
		Preload("Ingredients.Ingredient"). // load Ingredient details
		Preload("Ingredients.Unit").       // load Unit details
		Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).
		Find(&recipes)

	if result.Error != nil {
		http.Error(w, "Error fetching recipes.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// Find a recipe using the ID
func (app *App) getRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	id, err := strconv.Atoi(recipeID)
	if err != nil {
		http.Error(w, "invalid recipe ID", http.StatusBadRequest)
		return
	}

	var recipe models.Recipe

	result := app.Repo.DB.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("ingredient_id ASC")
	}).Preload("Ingredients.Ingredient"). // load Ingredient details
						Preload("Ingredients.Unit"). // load Unit details
						Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).First(&recipe, id)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Recipe with id %s not found", recipeID), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Find a recipe using its name
func (app *App) getRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeName := vars["name"]
	var recipe models.Recipe

	result := app.Repo.DB.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("recipe_id ASC")
	}).Preload("Ingredients.Ingredient"). // load Ingredient details
						Preload("Ingredients.Unit"). // load Unit details
						Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).First(&recipe, "name = ?", recipeName)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Recipe %s not found", recipeName), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) getNumberOfRecipes(w http.ResponseWriter, r *http.Request) {
	var recipes []models.Recipe
	result := app.Repo.DB.Find(&recipes)
	if result.Error != nil {
		app.Logger.Println("DB error:", result.Error)
	}
	app.Logger.Printf("Found %d recipes", len(recipes))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(len(recipes))
}
