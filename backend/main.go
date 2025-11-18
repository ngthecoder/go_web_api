package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/ngthecoder/go_web_api/internal/auth"
	"github.com/ngthecoder/go_web_api/internal/database"
	"github.com/ngthecoder/go_web_api/internal/ingredients"
	"github.com/ngthecoder/go_web_api/internal/recipes"
	"github.com/ngthecoder/go_web_api/internal/stats"
	"github.com/ngthecoder/go_web_api/internal/users"
)

func enableCORS(allowedOrigins []string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		for _, allowedOrigin := range allowedOrigins {
			if origin == allowedOrigin {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next(w, r)

		duration := time.Since(start)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, duration)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	err := database.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "unhealthy",
			"error":  "database connection failed",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func getDatabaseConfig() database.Config {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL != "" {
		log.Println("Using DATABASE_URL for database connection")
		return database.Config{URL: dbURL}
	}

	return database.Config{
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		DBName:   os.Getenv("POSTGRES_DB"),
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	dbConfig := getDatabaseConfig()
	port := os.Getenv("PORT")
	allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")

	if jwtSecret == "" {
		log.Fatal("Missing JWT_SECRET attribute")
	}
	log.Println("Successfully loaded JWT_SECRET")

	if port == "" {
		log.Fatal("Missing PORT attribute")
	}
	log.Println("Successfully loaded PORT")

	if len(allowedOrigins) == 0 {
		log.Fatal("Missing ALLOWED_ORIGINS attribute")
	}
	log.Println("Successfully loaded ALLOWED_ORIGINS")

	if err := database.InitDB(dbConfig); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := database.SeedData(); err != nil {
		log.Printf("Warning: failed to seed database: %v", err)
	}

	authService := auth.NewAuthService(database.DB, jwtSecret)
	authHandler := auth.NewAuthHandler(authService)

	userService := users.NewUserService(database.DB)
	userHandler := users.NewUserHandler(userService)

	recipesService := recipes.NewRecipesService(database.DB)
	recipesHandler := recipes.NewRecipesHandler(recipesService)

	ingredientsService := ingredients.NewIngredientsService(database.DB)
	ingredientsHandler := ingredients.NewIngredientsHandler(ingredientsService)

	statsService := stats.NewStatsService(database.DB)
	statsHandler := stats.NewStatsHandler(statsService)

	log.Println("Server running on port 8000")

	http.HandleFunc("/api/auth/register", loggingMiddleware(enableCORS(allowedOrigins, authHandler.RegisterHandler)))
	http.HandleFunc("/api/auth/login", loggingMiddleware(enableCORS(allowedOrigins, authHandler.LoginHandler)))

	http.HandleFunc("/api/user/profile", loggingMiddleware(enableCORS(allowedOrigins, authHandler.AuthMiddleware(userHandler.GetProfile))))
	http.HandleFunc("/api/user/liked-recipes", loggingMiddleware(enableCORS(allowedOrigins, authHandler.AuthMiddleware(userHandler.GetLikedRecipes))))
	http.HandleFunc("/api/user/liked-recipes/add", loggingMiddleware(enableCORS(allowedOrigins, authHandler.AuthMiddleware(userHandler.AddLikedRecipe))))
	http.HandleFunc("/api/user/liked-recipes/", loggingMiddleware(enableCORS(allowedOrigins, authHandler.AuthMiddleware(userHandler.RemoveLikedRecipe))))
	http.HandleFunc("/api/user/profile/update", loggingMiddleware(enableCORS(allowedOrigins, authHandler.AuthMiddleware(userHandler.UpdateProfile))))
	http.HandleFunc("/api/user/password", loggingMiddleware(enableCORS(allowedOrigins, authHandler.AuthMiddleware(userHandler.ChangePassword))))
	http.HandleFunc("/api/user/account", loggingMiddleware(enableCORS(allowedOrigins, authHandler.AuthMiddleware(userHandler.DeleteAccount))))

	http.HandleFunc("/api/recipes", loggingMiddleware(enableCORS(allowedOrigins, authHandler.OptionalAuthMiddleware(recipesHandler.AllRecipesHandler))))
	http.HandleFunc("/api/recipes/", loggingMiddleware(enableCORS(allowedOrigins, authHandler.OptionalAuthMiddleware(recipesHandler.RecipeDetailHandler))))
	http.HandleFunc("/api/recipes/find-by-ingredients", loggingMiddleware(enableCORS(allowedOrigins, authHandler.OptionalAuthMiddleware(recipesHandler.FindRecipesByIngredientsHandler))))
	http.HandleFunc("/api/recipes/shopping-list/", loggingMiddleware(enableCORS(allowedOrigins, recipesHandler.ShoppingListHandler)))

	http.HandleFunc("/api/ingredients", loggingMiddleware(enableCORS(allowedOrigins, ingredientsHandler.AllIngredientsHandler)))
	http.HandleFunc("/api/ingredients/", loggingMiddleware(enableCORS(allowedOrigins, ingredientsHandler.IngredientDetailsHandler)))

	http.HandleFunc("/api/categories", loggingMiddleware(enableCORS(allowedOrigins, statsHandler.CategoriesHandler)))
	http.HandleFunc("/api/stats", loggingMiddleware(enableCORS(allowedOrigins, statsHandler.StatsHandler)))

	http.HandleFunc("/api/health", loggingMiddleware(enableCORS(allowedOrigins, healthHandler)))

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
