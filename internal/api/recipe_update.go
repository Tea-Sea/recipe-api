package api

import (
	"encoding/json"
	"net/http"
	"recipe-api/internal/models"

	"github.com/gorilla/mux"
)

func (app *App) updateRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	result := app.Repo.DB.Model(&recipe).Where("id = ?", recipeID).Updates(recipe)
	if result.Error != nil {
		http.Error(w, "Failed to edit recipe by id", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) updateRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["name"]
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	result := app.Repo.DB.Model(&recipe).Where("name = ?", recipeID).Updates(recipe)
	if result.Error != nil {
		http.Error(w, "Failed to edit recipe by name", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipe)
}
