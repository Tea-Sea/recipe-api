package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ToPtr[T any](v T) *T {
	return &v
}

func TestAddRecipe(t *testing.T) {
	// Open temporary connection in memory
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	db.AutoMigrate(&Recipe{}, &Ingredient{}, &Unit{}, &RecipeIngredient{}, &Instruction{})
	app := &App{db: db}

	// Prepare request body
	newUnit := Unit{
		Label: "Cup",
	}

	newIngredient := Ingredient{
		Label: "Salt",
	}

	newIngredient2 := Ingredient{
		Label: "Water",
	}

	newRecipe := Recipe{
		Name:       "test recipe",
		Difficulty: 5,
		Ingredients: []RecipeIngredient{
			{
				Amount:     ToPtr(float32(1)),
				Ingredient: &newIngredient,
				Unit:       &newUnit,
			},
			{
				Amount:     ToPtr(float32(4)),
				Ingredient: &newIngredient2,
				Unit:       &newUnit,
			},
		},
		Instructions: []Instruction{
			{
				StepNumber: 1,
				StepText:   "boil it idk",
			},
			{
				StepNumber: 2,
				StepText:   "dry it?",
			},
		},
	}

	body, err := json.Marshal(newRecipe)
	if err != nil {
		t.Fatalf("failed to marshal recipe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/recipe/add", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	app.addRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf("expected status %d OK, got %d", http.StatusCreated, res.StatusCode)
	}

	var created Recipe
	if err := json.NewDecoder(res.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if created.RecipeID == 0 {
		t.Errorf("expected recipe ID to be set")
	}

	if created.Name != newRecipe.Name {
		t.Errorf("expected name %q, got %q", newRecipe.Name, created.Name)
	}

	if len(created.Ingredients) != 2 {
		t.Fatalf("expected 1 ingredient, got %d", len(created.Ingredients))
	}
	if created.Ingredients[0].RecipeIngredientID == 0 {
		t.Fatalf("expected RecipeIngredientID to be set")
	}
	if created.Ingredients[0].UnitID == nil || *created.Ingredients[0].UnitID == 0 {
		t.Fatalf("expected UnitID to be set")
	}
	if created.Ingredients[0].IngredientID == 0 {
		t.Fatalf("expected IngredientID to be set")
	}

	// Check instructions
	if len(created.Instructions) != 2 {
		t.Fatalf("expected 1 instruction, got %d", len(created.Instructions))
	}
	if created.Instructions[0].InstructionID == 0 {
		t.Fatalf("expected InstructionID to be set")
	}
}

// func TestDeleteByID(t *testing.T) {

// 	newRecipe := Recipe{
// 		RecipeID:   100,
// 		Name:       "test id",
// 		Difficulty: 5,
// 	}
// 	body, _ := json.Marshal(newRecipe)

// 	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewReader(body))
// 	w := httptest.NewRecorder()

// 	app.addRecipe(w, req)

// 	router := mux.NewRouter()
// 	router.HandleFunc("/recipes/id/{id}", deleteRecipeByID).Methods("DELETE")

// 	req = httptest.NewRequest(http.MethodDelete, "/recipes/id/100", nil)
// 	w = httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
// 	}

// 	var recipe Recipe
// 	result := db.Where("id = ?", newRecipe.RecipeID).First(&recipe)
// 	if result.Error == nil {
// 		t.Fatalf("recipe ID: '%v' was not deleted", newRecipe.RecipeID)
// 	}
// 	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		t.Fatalf("unexpected DB error: %v", result.Error)
// 	}
// }

// func TestDeleteByName(t *testing.T) {

// 	recipeName := "test recipe"

// 	router := mux.NewRouter()
// 	router.HandleFunc("/recipes/name/{name}", deleteRecipeByName).Methods("DELETE")

// 	req := httptest.NewRequest(http.MethodDelete, "/recipes/name/test%20recipe", nil)
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
// 	}

// 	var recipe Recipe
// 	result := db.Where("name = ?", recipeName).First(&recipe)
// 	if result.Error == nil {
// 		t.Fatalf("recipe '%s' was not deleted", recipeName)
// 	}
// 	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 		t.Fatalf("unexpected DB error: %v", result.Error)
// 	}
// }

// func TestGetByID(t *testing.T) {

// 	var expected = Recipe{
// 		Name:       "Chocolate Cake",
// 		Difficulty: 2,
// 	}

// 	router := mux.NewRouter()
// 	router.HandleFunc("/recipes/id/{id}", getRecipeByID).Methods("GET")

// 	req := httptest.NewRequest(http.MethodGet, "/recipes/id/1", nil)
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
// 	}

// 	var data = Recipe{}

// 	err := json.NewDecoder(res.Body).Decode(&data)
// 	if err != nil {
// 		t.Fatalf("Failed to decode response: %v", err)
// 	}

// 	if data.Name != expected.Name {
// 		t.Fatalf("Expected name %q, got %q", expected.Name, data.Name)
// 	}

// 	if data.Difficulty != expected.Difficulty {
// 		t.Fatalf("Expected Difficulty %q, got %q", expected.Difficulty, data.Difficulty)
// 	}

// 	// if data.Method != expected.Method {
// 	// 	t.Fatalf("Expected Method %q, got %q", expected.Method, data.Method)
// 	// }
// }

// func TestGetByName(t *testing.T) {

// 	var expected = Recipe{
// 		Name:       "test recipe",
// 		Difficulty: 5,
// 	}

// 	router := mux.NewRouter()
// 	router.HandleFunc("/recipes/name/{name}", getRecipeByName).Methods("GET")

// 	req := httptest.NewRequest(http.MethodGet, "/recipes/name/test%20recipe", nil)
// 	w := httptest.NewRecorder()

// 	router.ServeHTTP(w, req)

// 	res := w.Result()
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
// 	}

// 	var data = Recipe{}

// 	err := json.NewDecoder(res.Body).Decode(&data)
// 	if err != nil {
// 		t.Fatalf("Failed to decode response: %v", err)
// 	}

// 	if data.Name != expected.Name {
// 		t.Fatalf("Expected name %q, got %q", expected.Name, data.Name)
// 	}

// 	if data.Difficulty != expected.Difficulty {
// 		t.Fatalf("Expected Difficulty %q, got %q", expected.Difficulty, data.Difficulty)
// 	}

// 	// if data.Method != expected.Method {
// 	// 	t.Fatalf("Expected Method %q, got %q", expected.Method, data.Method)
// 	// }
// }
