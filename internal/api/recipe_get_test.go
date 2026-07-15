package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"recipe-api/internal/models"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetByID(t *testing.T) {
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

	// Get created id
	var created models.Recipe
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/recipe/id/{id}", testApp.getRecipeByID).Methods("GET")

	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/recipe/id/%d", created.RecipeID), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	getRes := w.Result()
	defer getRes.Body.Close()

	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d No Content, got %d", http.StatusOK, getRes.StatusCode)
	}

	var returned models.Recipe

	err = json.NewDecoder(getRes.Body).Decode(&returned)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if returned.RecipeID != created.RecipeID {
		t.Fatalf("Expected id %d, got %d", created.RecipeID, returned.RecipeID)
	}

	if returned.Name != created.Name {
		t.Fatalf("Expected name %q, got %q", created.Name, returned.Name)
	}

	if returned.Difficulty != created.Difficulty {
		t.Fatalf("Expected Difficulty %d, got %d", created.Difficulty, returned.Difficulty)
	}

	if len(returned.Ingredients) != len(created.Ingredients) {
		t.Fatalf("expected %d ingredients, got %d", len(created.Ingredients), len(returned.Ingredients))
	}

	if len(returned.Instructions) != len(created.Instructions) {
		t.Fatalf("expected %d instructions, got %d", len(created.Instructions), len(returned.Instructions))
	}

	// if returned.Method != created.Method {
	// 	t.Fatalf("Expected Method %q, got %q", created.Method, returned.Method)
	// }
}

func TestGetByName(t *testing.T) {
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

	// Get created id
	var created models.Recipe
	if err := json.NewDecoder(createRes.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/recipe/name/{name}", testApp.getRecipeByName).Methods("GET")

	urlSafeName := url.PathEscape(created.Name)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/recipe/name/%s", urlSafeName), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	getRes := w.Result()
	defer getRes.Body.Close()

	if getRes.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d No Content, got %d", http.StatusOK, getRes.StatusCode)
	}

	var returned models.Recipe

	err = json.NewDecoder(getRes.Body).Decode(&returned)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if returned.RecipeID != created.RecipeID {
		t.Fatalf("Expected id %d, got %d", created.RecipeID, returned.RecipeID)
	}

	if returned.Name != created.Name {
		t.Fatalf("Expected name %q, got %q", created.Name, returned.Name)
	}

	if returned.Difficulty != created.Difficulty {
		t.Fatalf("Expected Difficulty %d, got %d", created.Difficulty, returned.Difficulty)
	}

	if len(returned.Ingredients) != len(created.Ingredients) {
		t.Fatalf("expected %d ingredients, got %d", len(created.Ingredients), len(returned.Ingredients))
	}

	if len(returned.Instructions) != len(created.Instructions) {
		t.Fatalf("expected %d instructions, got %d", len(created.Instructions), len(returned.Instructions))
	}

	// if returned.Method != created.Method {
	// 	t.Fatalf("Expected Method %q, got %q", created.Method, returned.Method)
	// }
}
