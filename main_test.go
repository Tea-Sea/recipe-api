package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddRecipe(t *testing.T) {

	// Prepare request body
	newRecipe := Recipe{
		Name:       "Chocolate Cake",
		Difficulty: 3,
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
