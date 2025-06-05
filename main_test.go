package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func TestAddRecipe(t *testing.T) {

	// Prepare request body
	newRecipe := Recipe{
		Name:       "test recipe",
		Difficulty: 5,
		Method:     "This is where the method goes"}
	body, _ := json.Marshal(newRecipe)

	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	addRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	var created Recipe
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.ID == 0 {
		t.Errorf("expected recipe ID to be set")
	}

	if created.Name != newRecipe.Name {
		t.Errorf("expected title %q, got %q", newRecipe.Name, created.Name)
	}

	if created.Difficulty != newRecipe.Difficulty {
		t.Errorf("expected difficulty %v, got %v", newRecipe.Difficulty, created.Difficulty)
	}

}

func TestDeleteByID(t *testing.T) {

	newRecipe := Recipe{
		ID:         100,
		Name:       "test id",
		Difficulty: 5,
		Method:     "TESTING ID DELETION"}
	body, _ := json.Marshal(newRecipe)

	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewReader(body))
	w := httptest.NewRecorder()

	addRecipe(w, req)

	router := mux.NewRouter()
	router.HandleFunc("/recipes/id/{id}", deleteRecipeByID).Methods("DELETE")

	req = httptest.NewRequest(http.MethodDelete, "/recipes/id/100", nil)
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	var recipe Recipe
	result := db.Where("id = ?", newRecipe.ID).First(&recipe)
	if result.Error == nil {
		t.Fatalf("recipe ID: '%v' was not deleted", newRecipe.ID)
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("unexpected DB error: %v", result.Error)
	}
}

func TestDeleteByName(t *testing.T) {

	recipeName := "test recipe"

	router := mux.NewRouter()
	router.HandleFunc("/recipes/name/{name}", deleteRecipeByName).Methods("DELETE")

	req := httptest.NewRequest(http.MethodDelete, "/recipes/name/test%20recipe", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	var recipe Recipe
	result := db.Where("name = ?", recipeName).First(&recipe)
	if result.Error == nil {
		t.Fatalf("recipe '%s' was not deleted", recipeName)
	}
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		t.Fatalf("unexpected DB error: %v", result.Error)
	}
}

func TestGetByID(t *testing.T) {

	var expected = Recipe{
		Name:       "Pancakes",
		Difficulty: 2,
		Method:     "Mix ingredients and fry."}

	router := mux.NewRouter()
	router.HandleFunc("/recipes/id/{id}", getRecipeByID).Methods("GET")

	req := httptest.NewRequest(http.MethodGet, "/recipes/id/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	var data = Recipe{}

	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if data.Name != expected.Name {
		t.Fatalf("Expected name %q, got %q", expected.Name, data.Name)
	}

	if data.Difficulty != expected.Difficulty {
		t.Fatalf("Expected Difficulty %q, got %q", expected.Difficulty, data.Difficulty)
	}

	if data.Method != expected.Method {
		t.Fatalf("Expected Method %q, got %q", expected.Method, data.Method)
	}
}

func TestGetByName(t *testing.T) {

	var expected = Recipe{
		Name:       "test recipe",
		Difficulty: 5,
		Method:     "This is where the method goes"}

	router := mux.NewRouter()
	router.HandleFunc("/recipes/name/{name}", getRecipeByName).Methods("GET")

	req := httptest.NewRequest(http.MethodGet, "/recipes/name/test%20recipe", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	var data = Recipe{}

	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if data.Name != expected.Name {
		t.Fatalf("Expected name %q, got %q", expected.Name, data.Name)
	}

	if data.Difficulty != expected.Difficulty {
		t.Fatalf("Expected Difficulty %q, got %q", expected.Difficulty, data.Difficulty)
	}

	if data.Method != expected.Method {
		t.Fatalf("Expected Method %q, got %q", expected.Method, data.Method)
	}
}
