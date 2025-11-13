package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"recipe-api/internal/middleware"
)

func NewRouter(app *App, frontendURL string, appLogger *log.Logger) http.Handler {

	router := mux.NewRouter()

	// Default Landing Page
	router.HandleFunc("/", displayLanding)

	// Recipe endpoints
	// Add new Recipe
	router.HandleFunc("/recipe/add", app.addRecipe).Methods("POST")

	// Get all Recipes
	router.HandleFunc("/recipe/all", app.getAllRecipes).Methods("GET")

	// Get Recipes by value
	router.HandleFunc("/recipe/id/{id}", app.getRecipeByID).Methods("GET")
	router.HandleFunc("/recipe/name/{name}", app.getRecipeByName).Methods("GET")

	// Get Number of Recipes
	router.HandleFunc("/recipe/name/{name}", app.getNumberOfRecipes).Methods("GET")

	// Update Recipes
	router.HandleFunc("/recipe/id/{id}", app.updateRecipeByID).Methods("PUT")
	router.HandleFunc("/recipe/name/{name}", app.updateRecipeByName).Methods("PUT")

	// Delete Recipes
	router.HandleFunc("/recipe/id/{id}", app.deleteRecipeByID).Methods("DELETE")
	router.HandleFunc("/recipe/name/{name}", app.deleteRecipeByName).Methods("DELETE")

	//Filtered Recipes
	router.HandleFunc("/recipe/random", app.selectRandomRecipe).Methods("GET")
	router.HandleFunc("/recipe/random/{difficulty}", app.filterRandomRecipe).Methods("GET")

	http.Handle("/", router)

	handler := middleware.CorsMiddleware(router, frontendURL)

	appLogger.Fatal(http.ListenAndServe(":8080", handler))

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("404 Not Found: %s %s", r.Method, r.URL.Path)
		http.NotFound(w, r)
	})
	return router
}
