package api

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"

	"recipe-api/internal/models"
)

// Add new recipe
func (app *App) addRecipe(w http.ResponseWriter, r *http.Request) {
	var data models.Recipe
	check := json.NewDecoder(r.Body).Decode(&data)
	if check != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Rebuild structs
	var recipe = models.Recipe{}
	result := app.Repo.DB.Transaction(func(tx *gorm.DB) error {
		// Recipe
		recipe = models.Recipe{
			Name:       data.Name,
			Difficulty: data.Difficulty,
		}

		result := tx.Create(&recipe) // Check if exists
		if result.Error != nil {
			app.Logger.Println("Recipe error:", result.Error)
			// return result.Error
		}

		// Insert instructions
		for i := range data.Instructions {
			var instruction = models.Instruction{
				RecipeID:   recipe.RecipeID,
				StepNumber: data.Instructions[i].StepNumber,
				StepText:   data.Instructions[i].StepText,
				Duration:   data.Instructions[i].Duration,
				Notes:      data.Instructions[i].Notes,
			}

			result := tx.Create(&instruction)
			if result.Error != nil {
				app.Logger.Println("Instruction error:", result.Error)
				// return result.Error
			}
			recipe.Instructions = append(recipe.Instructions, instruction) // For return created object
		}

		//For every Recipe_Ingredient
		for i := range data.Ingredients {
			// Create new RI linker with recipeID and amount
			ri := models.RecipeIngredient{
				RecipeID: recipe.RecipeID,
				Amount:   data.Ingredients[i].Amount,
			}
			// Create the Ingredient
			if data.Ingredients[i].Ingredient != nil {
				var ingredient = models.Ingredient{
					Label: data.Ingredients[i].Ingredient.Label,
				}
				result := tx.FirstOrCreate(&ingredient, models.Ingredient{Label: ingredient.Label}) // Check if exists
				if result.Error != nil {
					app.Logger.Println("Ingredient error:", result.Error)
					// return result.Error
				}
				// Set IngredientID in linker
				ri.IngredientID = ingredient.IngredientID
			}
			// Create Unit
			if data.Ingredients[i].Unit != nil {
				var unit = models.Unit{
					Label: data.Ingredients[i].Unit.Label,
				}
				result := tx.FirstOrCreate(&unit, models.Unit{Label: unit.Label}) // Check if exists
				if result.Error != nil {
					app.Logger.Println("Unit error:", result.Error)
					// return result.Error
				}
				// Set UnitID in linker
				ri.UnitID = &unit.UnitID
			}
			// Create the linker
			result := tx.Create(&ri)
			if result.Error != nil {
				app.Logger.Println("RecipeIngredient error:", result.Error)
				// return result.Error
			}
			recipe.Ingredients = append(recipe.Ingredients, ri)
		}
		return nil
	})

	if result != nil {
		app.Logger.Println("Transaction Failed:", result)
		http.Error(w, result.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}
