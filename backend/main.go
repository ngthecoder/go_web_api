package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Ingredient struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Calories    int    `json:"calories_per_100g"`
	Description string `json:"description"`
}
type Recipe struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	PrepTimeMinutes int    `json:"prep_time_minutes"`
	CookTimeMinutes int    `json:"cook_time_minutes"`
	Servings        int    `json:"servings"`
	Difficulty      string `json:"difficulty"`
	Instructions    string `json:"instructions"`
	Description     string `json:"description"`
}

type RecipeIngredient struct {
	RecipeID     int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Notes        string  `json:"notes"`
}

type RecipeWithIngredients struct {
	Recipe      Recipe                   `json:"recipe"`
	Ingredients []IngredientWithQuantity `json:"ingredients"`
}

type IngredientWithQuantity struct {
	IngredientID int     `json:"ingredient_id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Notes        string  `json:"notes"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./foods.db")
	if err != nil {
		log.Fatal(err)
	}

	createIngredientsTable := `
		CREATE TABLE IF NOT EXISTS ingredients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			calories_per_100g INTEGER NOT NULL,
			description TEXT
		);
	`

	createRecipesTable := `
		CREATE TABLE IF NOT EXISTS recipes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			prep_time_minutes INTEGER NOT NULL,
			cook_time_minutes INTEGER NOT NULL,
			servings INTEGER NOT NULL,
			difficulty TEXT NOT NULL,
			instructions TEXT NOT NULL,
			description TEXT
		);
	`

	createRecipesIngredientsTable := `
		CREATE TABLE IF NOT EXISTS recipes_ingredients (
			recipe_id INTEGER NOT NULL,
			ingredient_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit TEXT NOT NULL,
			notes TEXT,
			PRIMARY KEY (recipe_id, ingredient_id),
			FOREIGN KEY (recipe_id) REFERENCES recipes (id),
			FOREIGN KEY (ingredient_id) REFERENCES ingredients (id)
		)
	`

	_, err = db.Exec(createIngredientsTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createRecipesTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createRecipesIngredientsTable)
	if err != nil {
		log.Fatal(err)
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
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

func helloHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	name := query.Get("name")

	if name == "" {
		name = "匿名"
	}

	response := map[string]string{
		"message": fmt.Sprintf("ようこそ、%sさん！", name),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ingredientsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM ingredients")
	if err != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient
		rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Category, &ingredient.Calories, &ingredient.Description)
		ingredients = append(ingredients, ingredient)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ingredients": ingredients,
		"total":       len(ingredients),
	})
}

// /api/recipes/{id} レシピ詳細＋使用食材取得
func recipeDetailHandler(w http.ResponseWriter, r *http.Request) {
	// URLからIDを取得
	idStr := r.URL.Path[len("/api/recipes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe id", http.StatusBadRequest)
		return
	}

	// レシピ情報を取得
	var recipe Recipe
	err = db.QueryRow(`
		SELECT id, name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description 
		FROM recipes WHERE id = ?`, id).
		Scan(&recipe.ID, &recipe.Name, &recipe.Category, &recipe.PrepTimeMinutes, &recipe.CookTimeMinutes,
			&recipe.Servings, &recipe.Difficulty, &recipe.Instructions, &recipe.Description)
	if err != nil {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// 使用食材を取得
	rows, err := db.Query(`
		SELECT i.id, i.name, ri.quantity, ri.unit, ri.notes
		FROM recipes_ingredients ri
		JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = ?`, id)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ingredients []IngredientWithQuantity
	for rows.Next() {
		var ing IngredientWithQuantity
		if err := rows.Scan(&ing.IngredientID, &ing.Name, &ing.Quantity, &ing.Unit, &ing.Notes); err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}
		ingredients = append(ingredients, ing)
	}

	// JSONレスポンス
	w.Header().Set("Content-Type", "application/json")
	resp := RecipeWithIngredients{Recipe: recipe, Ingredients: ingredients}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	initDB()
	defer db.Close()

	fmt.Printf("ポート8000でAPIサーバーを起動\n")

	http.HandleFunc("/api/hello", enableCORS(helloHandler))
	http.HandleFunc("/api/ingredients", enableCORS(ingredientsHandler))
	http.HandleFunc("/api/recipes/", enableCORS(recipeDetailHandler))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
