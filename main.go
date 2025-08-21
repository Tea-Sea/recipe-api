package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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
	RecipeID     int                `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string             `gorm:"unique" json:"name"`
	Difficulty   int                `json:"difficulty"`
	Description  *string            `json:"description,omitempty"` //optional
	Ingredients  []RecipeIngredient `gorm:"foreignKey:RecipeID" json:"ingredients,omitempty"`
	Instructions []Instruction      `gorm:"foreignKey:RecipeID" json:"instructions,omitempty"`
}

// Single Ingredient
type Ingredient struct {
	IngredientID int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Label        string `gorm:"type:varchar(32);not null" json:"label"`
}

type Unit struct {
	UnitID int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Label  string `gorm:"type:varchar(32);not null" json:"label"`
}

type RecipeIngredient struct {
	RecipeIngredientID int         `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipeID           int         `gorm:"not null;index" json:"recipe_id"`
	IngredientID       int         `gorm:"not null;index" json:"ingredient_id"`
	UnitID             *int        `gorm:"index" json:"unit_id,omitempty"`            //optional
	Amount             *float32    `gorm:"type:numeric(4,2)" json:"amount,omitempty"` //optional
	Ingredient         *Ingredient `gorm:"foreignKey:IngredientID;references:IngredientID" json:"ingredient"`
	Unit               *Unit       `gorm:"foreignKey:UnitID;references:UnitID" json:"unit,omitempty"`
}

type Instruction struct {
	InstructionID int     `gorm:"primaryKey;autoIncrement" json:"id"`
	RecipeID      int     `gorm:"not null;index" json:"recipe_id"`
	StepNumber    int     `gorm:"not null;check:step_number>0" json:"step_number"`
	StepText      string  `json:"step_text"`
	Duration      *int    `json:"duration,omitempty"` //optional
	Notes         *string `json:"notes,omitempty"`    //optional
}

// Applicaton struct to prevent use of globals
type App struct {
	db *gorm.DB
}

func init() {
	fmt.Println("App is starting…")
}

func main() {

	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

	app := &App{db: db}

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	app.db.AutoMigrate(&Recipe{})

	sqlDB, err := app.db.DB()
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

	router := mux.NewRouter()

	router.HandleFunc("/", displayLanding)

	router.HandleFunc("/recipe/add", app.addRecipe).Methods("POST")
	router.HandleFunc("/recipe/all", app.getAllRecipes).Methods("GET")

	router.HandleFunc("/recipe/id/{id}", app.getRecipeByID).Methods("GET")
	router.HandleFunc("/recipe/name/{name}", app.getRecipeByName).Methods("GET")

	router.HandleFunc("/recipe/id/{id}", app.updateRecipeByID).Methods("PUT")
	router.HandleFunc("/recipe/name/{name}", app.updateRecipeByName).Methods("PUT")

	router.HandleFunc("/recipe/id/{id}", app.deleteRecipeByID).Methods("DELETE")
	router.HandleFunc("/recipe/name/{name}", app.deleteRecipeByName).Methods("DELETE")

	router.HandleFunc("/recipe/random", app.selectRandomRecipe).Methods("GET")
	router.HandleFunc("/recipe/random/{difficulty}", app.filterRandomRecipe).Methods("GET")

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
func (app *App) getAllRecipes(w http.ResponseWriter, r *http.Request) {
	var recipes []Recipe

	result := app.db.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("ingredient_id ASC")
	}).
		Preload("Ingredients.Ingredient"). // load Ingredient details
		Preload("Ingredients.Unit").       // load Unit details
		Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).
		Find(&recipes)

	if result.Error != nil {
		http.Error(w, "Error fetching recipes.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}

// Add new recipe
func (app *App) addRecipe(w http.ResponseWriter, r *http.Request) {
	var data Recipe
	check := json.NewDecoder(r.Body).Decode(&data)
	if check != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	var recipe = Recipe{}
	result := app.db.Transaction(func(tx *gorm.DB) error {
		// Rebuild structs
		// Recipe
		recipe = Recipe{
			Name:       data.Name,
			Difficulty: data.Difficulty,
		}

		result := tx.Create(&recipe) // Check if exists
		if result.Error != nil {
			fmt.Println("Recipe error:", result.Error)
			// return result.Error
		}

		// Insert instructions
		for i := range data.Instructions {
			var instruction = Instruction{
				RecipeID:   recipe.RecipeID,
				StepNumber: data.Instructions[i].StepNumber,
				StepText:   data.Instructions[i].StepText,
				Duration:   data.Instructions[i].Duration,
				Notes:      data.Instructions[i].Notes,
			}

			result := tx.Create(&instruction)
			if result.Error != nil {
				fmt.Println("Instruction error:", result.Error)
				// return result.Error
			}
			recipe.Instructions = append(recipe.Instructions, instruction) // For return created object
		}

		//For every Recipe_Ingredient
		for i := range data.Ingredients {
			// Create new RI linker with recipeID and amount
			ri := RecipeIngredient{
				RecipeID: recipe.RecipeID,
				Amount:   data.Ingredients[i].Amount,
			}
			// Create the Ingredient
			if data.Ingredients[i].Ingredient != nil {
				var ingredient = Ingredient{
					Label: data.Ingredients[i].Ingredient.Label,
				}
				result := tx.FirstOrCreate(&ingredient, Ingredient{Label: ingredient.Label}) // Check if exists
				if result.Error != nil {
					fmt.Println("Ingredient error:", result.Error)
					// return result.Error
				}
				// Set IngredientID in linker
				ri.IngredientID = ingredient.IngredientID
			}
			// Create Unit
			if data.Ingredients[i].Unit != nil {
				var unit = Unit{
					Label: data.Ingredients[i].Unit.Label,
				}
				result := tx.FirstOrCreate(&unit, Unit{Label: unit.Label}) // Check if exists
				if result.Error != nil {
					fmt.Println("Unit error:", result.Error)
					// return result.Error
				}
				// Set UnitID in linker
				ri.UnitID = &unit.UnitID
			}
			// Create the linker
			result := tx.Create(&ri)
			if result.Error != nil {
				fmt.Println("RecipeIngredient error:", result.Error)
				// return result.Error
			}
			recipe.Ingredients = append(recipe.Ingredients, ri)
		}
		return nil
	})

	if result != nil {
		fmt.Println("Transaction Failed:", result)
		http.Error(w, result.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}

// Find a recipe using the ID
func (app *App) getRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	id, err := strconv.Atoi(recipeID)
	if err != nil {
		http.Error(w, "invalid recipe ID", http.StatusBadRequest)
		return
	}

	var recipe Recipe

	result := app.db.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Preload("Ingredients.Ingredient"). // load Ingredient details
						Preload("Ingredients.Unit"). // load Unit details
						Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).First(&recipe, id)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Recipe with id %s not found", recipeID), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Find a recipe using its name
func (app *App) getRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeName := vars["name"]
	var recipe Recipe

	result := app.db.Preload("Ingredients", func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}).Preload("Ingredients.Ingredient"). // load Ingredient details
						Preload("Ingredients.Unit"). // load Unit details
						Preload("Instructions", func(db *gorm.DB) *gorm.DB {
			return db.Order("step_number ASC")
		}).First(&recipe, "name = ?", recipeName)
	if result.Error != nil {
		http.Error(w, fmt.Sprintf("Recipe %s not found", recipeName), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) updateRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	result := app.db.Model(&recipe).Where("id = ?", recipeID).Updates(recipe)
	if result.Error != nil {
		http.Error(w, "Failed to edit recipe by id", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) updateRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["name"]
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	result := app.db.Model(&recipe).Where("name = ?", recipeID).Updates(recipe)
	if result.Error != nil {
		http.Error(w, "Failed to edit recipe by name", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(recipe)
}

// Delete recipe using id
func (app *App) deleteRecipeByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeID := vars["id"]
	result := app.db.Where("id = ?", recipeID).Delete(&Recipe{})
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
func (app *App) deleteRecipeByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	recipeName := vars["name"]
	result := app.db.Where("name = ?", recipeName).Delete(&Recipe{})
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

func (app *App) selectRandomRecipe(w http.ResponseWriter, r *http.Request) {
	//SELECT COUNT(*) FROM recipes
	var recipe Recipe
	result := app.db.Order("RANDOM()").First(&recipe)
	if result.Error != nil {
		http.Error(w, "No recipe found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) filterRandomRecipe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var recipe Recipe
	result := app.db.Where("difficulty <= ?", vars["difficulty"]).Order("RANDOM()").First(&recipe)
	if result.Error != nil {
		http.Error(w, "No recipe found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

func (app *App) numberOfRecipes() {
	var recipes []Recipe
	result := app.db.Find(&recipes)
	if result.Error != nil {
		log.Println("DB error:", result.Error)
	}
	log.Printf("Found %d recipes", len(recipes))
}
