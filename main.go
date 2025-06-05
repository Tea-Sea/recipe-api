package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Recipe struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Name       string `json:"name"`
	Difficulty int    `json:"difficulty"`
	Method     string `json:"method"`
}

var db *gorm.DB

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	db.AutoMigrate(&Recipe{})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get generic database object: %v", err)
	}

	if err != nil {
		log.Fatal("failed to get sql.DB from gorm.DB:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("database ping failed:", err)
	}

	fmt.Println("Database connection is alive.")

	// defer sqlDB.Close()
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", displayLanding)

	router.HandleFunc("/recipes", addRecipe).Methods("POST")
	router.HandleFunc("/recipes", getAllRecipes).Methods("GET")

	router.HandleFunc("/recipes/id/{id}", getRecipeByID).Methods("GET")
	router.HandleFunc("/recipes/name/{name}", getRecipeByName).Methods("GET")

	router.HandleFunc("/recipes/id/{id}", updateRecipeByID).Methods("PUT")
	router.HandleFunc("/recipes/name/{name}", updateRecipeByName).Methods("PUT")

	router.HandleFunc("/recipes/id/{id}", deleteRecipeByID).Methods("DELETE")
	router.HandleFunc("/recipes/name/{name}", deleteRecipeByName).Methods("DELETE")

	router.HandleFunc("/recipes/random", selectRandomRecipe).Methods("GET")
	router.HandleFunc("/recipes/random/{difficulty}", filterRandomRecipe).Methods("GET")

	http.Handle("/", router)

	fmt.Printf("Server listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", router))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 Not Found: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})

}

// Default
func displayLanding(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You have successfully connected using this API"))
}

// Get all recipes
func getAllRecipes(w http.ResponseWriter, r *http.Request) {
	var recipe []Recipe
	db.Find(&recipe)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Add new recipe
func addRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	result := db.Where("name = ?", recipe.Name).First(&recipe)
	if result.Error == nil {
		http.Error(w, "Recipe with that name already exists.", http.StatusInternalServerError)
		return
	}
	db.Create(&recipe)
	json.NewEncoder(w).Encode(recipe)
}

// Find a recipe using the ID
func getRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	var recipe Recipe
	db.First(&recipe, recipeID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Find a recipe using its name
func getRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeName := vars["name"]
	var recipe Recipe
	result := db.Where("name = ?", recipeName).First(&recipe)
	if result.Error != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func updateRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	result := db.Model(&recipe).Where("id = ?", recipeID).Updates(recipe)
	if result.Error != nil {
		http.Error(w, "Failed to edit recipe by id", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

func updateRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["name"]
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	result := db.Model(&recipe).Where("name = ?", recipeID).Updates(recipe)
	if result.Error != nil {
		http.Error(w, "Failed to edit recipe by name", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

// Delete recipe using id
func deleteRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	result := db.Where("id = ?", recipeID).Delete(&Recipe{})
	if result.Error != nil {
		http.Error(w, "Failed to delete recipe by ID", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Recipe '%s' deleted successfully", recipeID)
}

// Delete recipe using name
func deleteRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeName := vars["name"]
	result := db.Where("name = ?", recipeName).Delete(&Recipe{})
	if result.Error != nil {
		http.Error(w, "Failed to delete recipe", http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	fmt.Fprintf(w, "Recipe '%s' deleted successfully", recipeName)
}

func selectRandomRecipe(w http.ResponseWriter, r *http.Request) {
	//SELECT COUNT(*) FROM recipes
	var recipe Recipe
	result := db.Order("RANDOM()").First(&recipe)
	if result.Error != nil {
		http.Error(w, "No recipe found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func filterRandomRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var recipe Recipe
	result := db.Where("difficulty <= ?", vars["difficulty"]).Order("RANDOM()").First(&recipe)
	if result.Error != nil {
		http.Error(w, "No recipe found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func numberOfRecipes() {
	var recipes []Recipe
	result := db.Find(&recipes)
	if result.Error != nil {
		log.Println("DB error:", result.Error)
	}
	log.Printf("Found %d recipes", len(recipes))
}
