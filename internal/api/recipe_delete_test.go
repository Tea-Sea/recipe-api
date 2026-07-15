package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"recipe-api/internal/models"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func TestDeleteByID(t *testing.T) {
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
	router.HandleFunc("/recipe/id/{id}", testApp.deleteRecipeByID).Methods("DELETE")

	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/recipe/id/%d", created.RecipeID), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	deleteRes := w.Result()
	defer deleteRes.Body.Close()

	if deleteRes.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d No Content, got %d", http.StatusNoContent, deleteRes.StatusCode)
	}

	var removed models.Recipe
	result := testApp.Repo.DB.Where("recipe_id = ?", created.RecipeID).First(&removed)
	if result.Error == nil {
		t.Fatalf("recipe ID: '%v' was not deleted", recipe.RecipeID)
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("unexpected DB error: %v", result.Error)
	}
}

func TestDeleteByName(t *testing.T) {
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
	router.HandleFunc("/recipe/name/{name}", testApp.deleteRecipeByName).Methods("DELETE")

	urlSafeName := url.PathEscape(created.Name)
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/recipe/name/%s", urlSafeName), nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	deleteRes := w.Result()
	defer deleteRes.Body.Close()

	if deleteRes.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d No Content, got %d", http.StatusNoContent, deleteRes.StatusCode)
	}

	var removed models.Recipe
	result := testApp.Repo.DB.Where("name = ?", created.Name).First(&removed)
	if result.Error == nil {
		t.Fatalf("recipe ID: '%v' was not deleted", created.Name)
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("unexpected DB error: %v", result.Error)
	}
}
