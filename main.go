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

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Edit CORs settings in future
		frontend, ok := os.LookupEnv("FRONTEND_URL")
		if !ok {
			fmt.Println("FRONTEND_URL is not set")
		}
		w.Header().Set("Access-Control-Allow-Origin", frontend)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type Recipe struct {
	RecipeID     int                `gorm:"primaryKey" json:"id"`
	Name         string             `json:"name"`
	Difficulty   int                `json:"difficulty"`
	Description  *string            `json:"decription,omitempty"` //optional
	Ingredients  []RecipeIngredient `gorm:"foreignKey:RecipeID" json:"ingredients,omitempty"`
	Instructions []Instruction      `gorm:"foreignKey:RecipeID" json:"instructions,omitempty"`
}

// Single Ingredient
type Ingredient struct {
	IngredientID int    `gorm:"primaryKey" json:"id"`
	Label        string `gorm:"type:varchar(32);not null" json:"label"`
	Sort         int    `gorm:"default:0;check:sort>=0" json:"sort"`
}

type Unit struct {
	UnitID int    `gorm:"primaryKey" json:"id"`
	Label  string `gorm:"type:varchar(32);not null" json:"label"`
	Sort   int    `gorm:"default:0;check:sort>=0" json:"sort"`
}

type RecipeIngredient struct {
	RecipeIngredientID int         `gorm:"primaryKey" json:"id"`
	RecipeID           int         `gorm:"not null;index" json:"recipe_id"`
	IngredientID       int         `gorm:"not null;index" json:"ingredient_id"`
	UnitID             *int        `gorm:"index" json:"unit_id,omitempty"`            //optional
	Amount             *float32    `gorm:"type:numeric(4,2)" json:"amount,omitempty"` //optional
	Sort               int         `gorm:"default:0;check:sort>=0" json:"sort"`
	Ingredient         *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient"`
	Unit               *Unit       `gorm:"foreignKey:UnitID;references:UnitID" json:"unit,omitempty"`
}

type Instruction struct {
	InstructionID int     `gorm:"primaryKey" json:"id"`
	RecipeID      int     `gorm:"not null;index" json:"recipe_id"`
	StepNumber    int     `gorm:"not null;check:step_number>0" json:"step_number"`
	StepText      string  `json:"text"`
	Duration      *int    `json:"duration_minutes,omitempty"` //optional
	Notes         *string `json:"notes,omitempty"`            //optional
	Sort          int     `gorm:"default:0;check:sort>=0" json:"sort"`
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

	router.HandleFunc("/recipe/add", addRecipe).Methods("POST")
	router.HandleFunc("/recipe/all", getAllRecipes).Methods("GET")

	router.HandleFunc("/recipe/id/{id}", getRecipeByID).Methods("GET")
	router.HandleFunc("/recipe/name/{name}", getRecipeByName).Methods("GET")

	router.HandleFunc("/recipe/id/{id}", updateRecipeByID).Methods("PUT")
	router.HandleFunc("/recipe/name/{name}", updateRecipeByName).Methods("PUT")

	router.HandleFunc("/recipe/id/{id}", deleteRecipeByID).Methods("DELETE")
	router.HandleFunc("/recipe/name/{name}", deleteRecipeByName).Methods("DELETE")

	router.HandleFunc("/recipe/random", selectRandomRecipe).Methods("GET")
	router.HandleFunc("/recipe/random/{difficulty}", filterRandomRecipe).Methods("GET")

	http.Handle("/", router)

	handler := corsMiddleware(router)

	fmt.Printf("Server listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", handler))

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

	err := db.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("sort ASC")
	}).Preload("Ingredients.Ingredient"). // load Ingredient details
						Preload("Ingredients.Unit"). // load Unit details
						Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).First(&recipe, recipeID).Error
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", recipe)
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
