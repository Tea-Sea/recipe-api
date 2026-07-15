package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"recipe-api/internal/models"
	"testing"
)

func TestAddRecipe(t *testing.T) {
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

	res := w.Result()
	defer res.Body.Close()
	// Was created successfully
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d OK, got %d", http.StatusCreated, res.StatusCode)
	}
	// Created recipe matches one posted
	var created models.Recipe
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.RecipeID == 0 {
		t.Errorf("expected recipe ID to be set")
	}

	if created.Name != recipe.Name {
		t.Errorf("expected name %q, got %q", recipe.Name, created.Name)
	}

	if len(created.Ingredients) != 2 {
		t.Fatalf("expected 2 ingredients, got %d", len(created.Ingredients))
	}
	if created.Ingredients[0].IngredientID == 0 {
		t.Fatalf("expected IngredientID to be set")
	}
	if created.Ingredients[0].RecipeIngredientID == 0 {
		t.Fatalf("expected RecipeIngredientID to be set")
	}
	if created.Ingredients[0].UnitID == nil || *created.Ingredients[0].UnitID == 0 {
		t.Fatalf("expected UnitID to be set")
	}

	// Check instructions
	if len(created.Instructions) != 2 {
		t.Fatalf("expected 2 instructions, got %d", len(created.Instructions))
	}
	if created.Instructions[0].InstructionID == 0 {
		t.Fatalf("expected InstructionID to be set")
	}

	clearDatabase(testApp)
}
