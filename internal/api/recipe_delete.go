package api

import (
	"net/http"
	"recipe-api/internal/models"
	"strconv"

	"github.com/gorilla/mux"
)

// Delete recipe using id
func (app *App) deleteRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	// Convert to int
	id, err := strconv.Atoi(recipeID)
	if err != nil {
		http.Error(w, "invalid recipe ID", http.StatusBadRequest)
		return
	}

	var check models.Recipe
	result := app.Repo.DB.First(&check, id)
	if result.Error != nil {
		app.Logger.Println("Recipe not found")
		return
	}

	result = app.Repo.DB.Delete(&models.Recipe{}, id)
	if result.Error != nil {
		http.Error(w, "Failed to delete recipe by ID", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Table unaffacted", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	app.Logger.Printf("Recipe '%s' deleted successfully", recipeID)
}

// Delete recipe using name
func (app *App) deleteRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeName := vars["name"]

	var check models.Recipe
	result := app.Repo.DB.First(&check, check.Name)
	if result.Error != nil {
		app.Logger.Println("Recipe not found")
		return
	}

	result = app.Repo.DB.Where("name = ?", recipeName).Delete(&models.Recipe{})
	if result.Error != nil {
		http.Error(w, "Failed to delete recipe by name", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Table unaffacted", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	app.Logger.Printf("Recipe '%s' deleted successfully", recipeName)
}
