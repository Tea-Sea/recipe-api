package api

import (
	"encoding/json"
	"net/http"
	"recipe-api/internal/models"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func (app *App) updateRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	var recipe models.Recipe
	json.NewDecoder(r.Body).Decode(&recipe)

	id, err := strconv.Atoi(recipeID)
	if err != nil {
		http.Error(w, "invalid recipe ID", http.StatusBadRequest)
		return
	}
	result := app.Repo.DB.Transaction(func(tx *gorm.DB) error {

		var check models.Recipe
		result := tx.First(&check, id)
		if result.Error != nil {
			app.Logger.Println("Recipe not found")
			return result.Error
		}

		// Update Recipe object
		result = tx.Model(&models.Recipe{}).Where("recipe_id = ?", recipeID).Updates(recipe)
		if result.Error != nil {
			app.Logger.Println(w, "Failed to edit recipe by id", http.StatusInternalServerError)
			return result.Error
		}

		// Delete child linker
		result = tx.Where("recipe_id = ?", id).Delete(&models.RecipeIngredient{})
		if result.Error != nil {
			app.Logger.Println(w, "Failed to edit child linker object", http.StatusInternalServerError)
			return result.Error
		}

		// Delete Instructions object
		result = tx.Where("recipe_id = ?", id).Delete(&models.Instruction{})
		if result.Error != nil {
			app.Logger.Println(w, "Failed to delete child instructions", http.StatusInternalServerError)
			return result.Error
		}

		// Rebuild Ingredient objects to ensure existance (including units)
		for i := range recipe.Ingredients {
			recipe.Ingredients[i].RecipeID = id

			result = tx.Where("label = ?", recipe.Ingredients[i].Unit.Label).FirstOrCreate(recipe.Ingredients[i].Unit)
			if result.Error != nil {
				app.Logger.Println(w, "Failed to rebuild Unit objects")
				return result.Error
			}

			result = tx.Where("label = ?", recipe.Ingredients[i].Ingredient.Label).FirstOrCreate(recipe.Ingredients[i].Ingredient)
			if result.Error != nil {
				app.Logger.Println(w, "Failed to rebuild ingredient objects")
				return result.Error
			}

			// Rebuild Linker object
			linker := models.RecipeIngredient{
				RecipeID:     id,
				IngredientID: recipe.Ingredients[i].Ingredient.IngredientID,
				UnitID:       &recipe.Ingredients[i].Unit.UnitID,
				Amount:       recipe.Ingredients[i].Amount,
			}
			result = tx.Create(&linker)
			if result.Error != nil {
				app.Logger.Println(w, "Recipe Ingredient linking failed for %s", recipe.Ingredients[i].Ingredient.Label)
				return result.Error
			}
		}

		// Rebuild Instruction objects
		for i := range recipe.Instructions {
			recipe.Instructions[i].RecipeID = id
			result = tx.Create(&recipe.Instructions[i])
			if result.Error != nil {
				app.Logger.Println(w, "Failed to rebuild ingredient objects")
			}
		}
		return nil
	})

	if result != nil {
		app.Logger.Println("Transaction Failed:", result)
		http.Error(w, result.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) updateRecipeByName(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	// recipeID := vars["name"]
	// var recipe models.Recipe
	// json.NewDecoder(r.Body).Decode(&recipe)
	// result := app.Repo.DB.Model(&recipe).Where("name = ?", recipeID).Updates(recipe)
	//
	//	if result.Error != nil {
	//		http.Error(w, "Failed to edit recipe by name", http.StatusInternalServerError)
	//		return
	//	}
	//
	// json.NewEncoder(w).Encode(recipe)
}
