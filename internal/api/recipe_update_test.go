package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"recipe-api/internal/models"
	"testing"

	"github.com/gorilla/mux"
)

func TestUpdateByID(t *testing.T) {
	defer clearDatabase(testApp)
	// Create the test recipe
	recipe := createTestRecipe(t, testApp)

	// Prepare request body
	body, err := json.Marshal(recipe)
	if err != nil {
		t.Fatalf("failed to marshal recipe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/recipe/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	testApp.addRecipe(w, req)

	createRes := w.Result()
	defer createRes.Body.Close()
	// Was created successfully
	if createRes.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d OK, got %d", http.StatusCreated, createRes.StatusCode)
	}

	// Update recipe
	recipe.Name = "Recipe Tester"
	recipe.Difficulty = 5
	recipe.Ingredients = []models.RecipeIngredient{
		{
			Amount:     ToPtr(float32(10)),
			Ingredient: createTestIngredient("Sugar"),
			Unit:       createTestUnit("Gram"),
		},
		{
			Amount:     ToPtr(float32(4)),
			Ingredient: createTestIngredient("Water"),
			Unit:       createTestUnit("L"),
		},
		{
			Amount:     ToPtr(float32(1)),
			Ingredient: createTestIngredient("Salt"),
			Unit:       createTestUnit("gram"),
		},
	}

	// Get created id
	var recipeInDB models.Recipe
	if err := json.NewDecoder(createRes.Body).Decode(&recipeInDB); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/recipe/id/{id}", testApp.updateRecipeByID).Methods("PUT")

	updateBody, err := json.Marshal(recipe)
	if err != nil {
		t.Fatal(err)
	}

	req = httptest.NewRequest(http.MethodPut, fmt.Sprintf("/recipe/id/%d", recipeInDB.RecipeID), bytes.NewReader(updateBody))
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	updateRes := w.Result()
	defer updateRes.Body.Close()

	if updateRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d OK, got %d", http.StatusOK, updateRes.StatusCode)
	}

	var updated models.Recipe
	result := testApp.Repo.DB.
		Preload("Ingredients.Ingredient").
		Preload("Ingredients.Unit").
		Preload("Instructions").
		First(&updated, recipeInDB.RecipeID)

	if updated.RecipeID != recipeInDB.RecipeID {
		t.Fatalf("Expected id %d, got %d", recipeInDB.RecipeID, updated.RecipeID)
	}

	if updated.Name != recipe.Name {
		t.Fatalf("Expected name %s, got %s. Recipe was not updated", recipe.Name, updated.Name)
	}
	if updated.Difficulty != recipe.Difficulty {
		t.Fatalf("Expected updated difficulty of %d, got %d,", recipe.Difficulty, updated.Difficulty)
	}

	if len(updated.Ingredients) == len(recipeInDB.Ingredients) {
		t.Fatalf("Expected %d ingredients, got %d. Recipe was not updated", len(recipe.Ingredients), len(updated.Ingredients))
	}

	if updated.Ingredients[0].Ingredient.Label != recipe.Ingredients[0].Ingredient.Label {
		t.Fatalf("Expected first ingredient %s, got %s. Ingredient was not updated", recipe.Ingredients[0].Ingredient.Label, updated.Ingredients[0].Ingredient.Label)
	}

	if result.Error != nil {
		t.Fatalf("Failed to load updated recipe: '%v'", result.Error)
	}
}

// func TestUpdateByName(t *testing.T) {
// 	defer clearDatabase(testApp)
// 	// Create the test recipe
// 	recipe := createTestRecipe(t, testApp)

// 	// Prepare request body
// 	body, err := json.Marshal(recipe)
// 	if err != nil {
// 		t.Fatalf("failed to marshal recipe: %v", err)
// 	}

// 	req := httptest.NewRequest(http.MethodPost, "/recipe/add", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	testApp.addRecipe(w, req)

// 	createRes := w.Result()
// 	defer createRes.Body.Close()
// 	// Was created successfully
// 	if createRes.StatusCode != http.StatusCreated {
// 		t.Fatalf("expected status %d OK, got %d", http.StatusCreated, createRes.StatusCode)
// 	}

// 	// Get created id
// 	var created models.Recipe
// 	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
// 		t.Fatalf("failed to decode response: %v", err)
// 	}

// 	router := mux.NewRouter()
// 	router.HandleFunc("/recipe/name/{name}", testApp.deleteRecipeByName).Methods("DELETE")

// 	urlSafeName := url.PathEscape(created.Name)
// 	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/recipe/name/%s", urlSafeName), nil)
// 	w = httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	deleteRes := w.Result()
// 	defer deleteRes.Body.Close()

// 	if deleteRes.StatusCode != http.StatusNoContent {
// 		t.Fatalf("expected status %d No Content, got %d", http.StatusNoContent, deleteRes.StatusCode)
// 	}

// 	var removed models.Recipe
// 	result := testApp.Repo.DB.Where("name = ?", created.Name).First(&removed)
// 	if result.Error == nil {
// 		t.Fatalf("recipe ID: '%v' was not deleted", created.Name)
// 	}
// 	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		t.Fatalf("unexpected DB error: %v", result.Error)
// 	}
// }
