package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type Recipe struct {
	id         int    `json:"id"`
	name       string `json:"name"`
	difficulty int    `json:"difficulty"`
	method     string `json:"method"`
}

var db *gorm.DB

func init() {
	// connStr := "host=localhost port=5432 user=postgres password=testpassword dbname=recipes sslmode=disable"
	// db, err := pgx.Connect(context.Background(), connStr)
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	db.AutoMigrate(&Recipe{})

	sqlDB, err := db.DB()

	defer sqlDB.Close()

	if err != nil {
		log.Fatal("failed to get sql.DB from gorm.DB:", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatal("database ping failed:", err)
	}

	fmt.Println("Database connection is alive.")
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/recipes", getRecipes).Methods("GET")
	// router.HandleFunc("/recipes", addRecipe).Methods("POST")
	// router.HandleFunc("/recipes/{id}", getRecipe).Methods("GET")
	// router.HandleFunc("/recipes/{id}", updateRecipe).Methods("PUT")
	// router.HandleFunc("/recipes/{id}", deleteRecipe).Methods("DELETE")

	http.Handle("/", router)

	fmt.Printf("Server listening on port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Get all users
func getRecipes(w http.ResponseWriter, r *http.Request) {
	var recipe []Recipe
	db.Find(&recipe)
	json.NewEncoder(w).Encode(recipe)
}
