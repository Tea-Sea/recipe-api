package api

import (
	"encoding/json"
	"net/http"
	"recipe-api/internal/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (app *App) selectRandomRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe

	result := app.Repo.DB.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("ingredient_id ASC")
	}).Preload("Ingredients.Ingredient"). // load Ingredient details
						Preload("Ingredients.Unit"). // load Unit details
						Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).Order("RANDOM()").First(&recipe)
	if result.Error != nil {
		http.Error(w, "Random recipe not retrieved", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) filterRandomRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var recipe models.Recipe
	result := app.Repo.DB.Where("difficulty <= ?", vars["difficulty"]).Order("RANDOM()").First(&recipe)
	if result.Error != nil {
		http.Error(w, "No recipe found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}
