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
	router.HandleFunc("/recipes", getRecipes).Methods("GET")
	router.HandleFunc("/recipes", addRecipe).Methods("POST")
	// router.HandleFunc("/recipes/{id}", getRecipe).Methods("GET")
	// router.HandleFunc("/recipes/{id}", updateRecipe).Methods("PUT")
	// router.HandleFunc("/recipes/{id}", deleteRecipe).Methods("DELETE")

	http.Handle("/", router)

	fmt.Printf("Server listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Default
func displayLanding(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("You have successfully connected using this API"))
}

// Get all recipes
func getRecipes(w http.ResponseWriter, r *http.Request) {

	var recipe []Recipe
	db.Find(&recipe)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipe)
}

// Add new recipe
func addRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe Recipe
	json.NewDecoder(r.Body).Decode(&recipe)
	db.Create(&recipe)
	json.NewEncoder(w).Encode(recipe)
}
