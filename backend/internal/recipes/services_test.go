package recipes

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestRecipeServiceBasics(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewRecipesService(db)

	if service == nil {
		t.Fatal("NewRecipesService returned nil")
	}

	recipes, err := service.matchedRecipesRetriever("partial", []int{1}, 10, "")

	if err != nil {
		t.Fatalf("matchedRecipesRetriever() error: %v", err)
	}

	t.Logf("Query returned %d recipes", len(recipes))
}

func TestRecipeDetailsWithIngredientsRetriever(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	service := NewRecipesService(db)

	recipe, ingredients, err := service.recipeDetailsWithIngredientsRetriever(1, "")
	if err != nil {
		t.Fatalf("recipeDetailsWithIngredientsRetriever() error: %v", err)
	}

	if recipe.Name != "Tomato Rice" {
		t.Errorf("Recipe name = %v, want Tomato Rice", recipe.Name)
	}

	if len(ingredients) != 3 {
		t.Errorf("Got %d ingredients, want 3", len(ingredients))
	}

	for _, ing := range ingredients {
		if ing.Name == "" {
			t.Error("Ingredient has empty name")
		}
		if ing.Quantity <= 0 {
			t.Error("Ingredient has zero or negative quantity")
		}
	}

	t.Logf("Recipe details test passed!")
}

func TestDatabaseSetup(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM recipes").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count recipes: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected 2 recipes, got %d", count)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM ingredients").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count ingredients: %v", err)
	}

	if count != 4 {
		t.Errorf("Expected 4 ingredients, got %d", count)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM recipe_ingredients").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count recipe_ingredients: %v", err)
	}

	if count != 5 {
		t.Errorf("Expected 5 recipe_ingredients, got %d", count)
	}

	t.Logf("Database setup test passed!")
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE recipes (
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

		CREATE TABLE ingredients (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			calories_per_100g INTEGER NOT NULL,
			description TEXT
		);

		CREATE TABLE recipe_ingredients (
			recipe_id INTEGER NOT NULL,
			ingredient_id INTEGER NOT NULL,
			quantity REAL NOT NULL,
			unit TEXT NOT NULL,
			notes TEXT,
			PRIMARY KEY (recipe_id, ingredient_id)
		);

		CREATE TABLE users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE user_liked_recipes (
			user_id TEXT NOT NULL,
			recipe_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, recipe_id)
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create tables: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO ingredients (id, name, category, calories_per_100g, description) VALUES
		(1, 'Tomato', 'Vegetables', 18, 'Fresh tomatoes'),
		(2, 'Onion', 'Vegetables', 40, 'Yellow onions'),
		(3, 'Rice', 'Grains', 130, 'White rice'),
		(4, 'Chicken', 'Protein', 165, 'Chicken breast');

		INSERT INTO recipes (id, name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description) VALUES
		(1, 'Tomato Rice', 'lunch', 10, 20, 4, 'easy', 'Cook rice with tomatoes', 'Delicious'),
		(2, 'Chicken Rice', 'dinner', 15, 30, 4, 'medium', 'Cook rice with chicken', 'Tasty');

		INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, unit, notes) VALUES
		(1, 1, 2, 'pieces', 'diced'),
		(1, 2, 1, 'piece', 'chopped'),
		(1, 3, 2, 'cups', 'uncooked'),
		(2, 4, 1, 'piece', 'diced'),
		(2, 3, 2, 'cups', 'cooked');
	`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	return db
}
