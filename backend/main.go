package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ngthecoder/go_web_api/internal/auth"
	"github.com/ngthecoder/go_web_api/internal/ingredients"
	"github.com/ngthecoder/go_web_api/internal/recipes"
	"github.com/ngthecoder/go_web_api/internal/users"
)

var db *sql.DB

type RecipeIngredient struct {
	RecipeID     int     `json:"recipe_id"`
	IngredientID int     `json:"ingredient_id"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Notes        string  `json:"notes"`
}

type CategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

type CategoryCountsResponse struct {
	IngredientCategories []CategoryCount `json:"ingredient_categories"`
	RecipeCategories     []CategoryCount `json:"recipe_categories"`
}

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./foods.db")
	if err != nil {
		log.Fatal(err)
	}

	db.Exec("DROP TABLE IF EXISTS recipe_ingredients")
	db.Exec("DROP TABLE IF EXISTS recipes")
	db.Exec("DROP TABLE IF EXISTS ingredients")
	db.Exec("DROP TABLE IF EXISTS users")
	db.Exec("DROP TABLE IF EXISTS user_liked_recipes")

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
		CREATE TABLE IF NOT EXISTS recipe_ingredients (
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

	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
	`

	createUserLikedRecipesTable := `
		CREATE TABLE IF NOT EXISTS user_liked_recipes (
			user_id TEXT NOT NULL,
			recipe_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, recipe_id),
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY (recipe_id) REFERENCES recipes(id) ON DELETE CASCADE
		);
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

	_, err = db.Exec(createUsersTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createUserLikedRecipesTable)
	if err != nil {
		log.Fatal(err)
	}

	createIndexes()
}

func createIndexes() {
	indexQueries := []string{
		"CREATE INDEX IF NOT EXISTS idx_ingredients_category ON ingredients(category)",
		"CREATE INDEX IF NOT EXISTS idx_ingredients_name ON ingredients(name)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_category ON recipes(category)",
		"CREATE INDEX IF NOT EXISTS idx_recipes_difficulty ON recipes(difficulty)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_recipe_id ON recipe_ingredients(recipe_id)",
		"CREATE INDEX IF NOT EXISTS idx_recipe_ingredients_ingredient_id ON recipe_ingredients(ingredient_id)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);",
		"CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);",
	}

	for _, query := range indexQueries {
		_, err := db.Exec(query)
		if err != nil {
			log.Printf("Error creating index: %v", err)
		}
	}
}

func populateTestData() {
	db.Exec("DELETE FROM recipe_ingredients")
	db.Exec("DELETE FROM recipes")
	db.Exec("DELETE FROM ingredients")

	// Comprehensive English ingredients data
	ingredientsData := []struct {
		name        string
		category    string
		calories    int
		description string
	}{
		// Vegetables
		{"Tomato", "Vegetables", 18, "Fresh red tomatoes"},
		{"Onion", "Vegetables", 40, "Sweet yellow onions"},
		{"Garlic", "Vegetables", 149, "Fresh garlic cloves"},
		{"Bell Pepper", "Vegetables", 31, "Colorful bell peppers"},
		{"Carrot", "Vegetables", 41, "Orange carrots"},
		{"Celery", "Vegetables", 16, "Crisp celery stalks"},
		{"Potato", "Vegetables", 77, "Russet potatoes"},
		{"Sweet Potato", "Vegetables", 86, "Orange sweet potatoes"},
		{"Broccoli", "Vegetables", 34, "Fresh broccoli florets"},
		{"Spinach", "Vegetables", 23, "Fresh spinach leaves"},
		{"Lettuce", "Vegetables", 15, "Crisp romaine lettuce"},
		{"Cucumber", "Vegetables", 16, "Fresh cucumbers"},
		{"Zucchini", "Vegetables", 17, "Green zucchini"},
		{"Mushrooms", "Vegetables", 22, "Button mushrooms"},
		{"Asparagus", "Vegetables", 20, "Fresh asparagus spears"},
		{"Green Beans", "Vegetables", 31, "Fresh green beans"},
		{"Corn", "Vegetables", 86, "Sweet corn kernels"},
		{"Peas", "Vegetables", 81, "Green peas"},
		{"Cabbage", "Vegetables", 25, "Fresh green cabbage"},
		{"Cauliflower", "Vegetables", 25, "White cauliflower"},
		{"Brussels Sprouts", "Vegetables", 43, "Fresh Brussels sprouts"},
		{"Kale", "Vegetables", 35, "Nutritious kale leaves"},
		{"Eggplant", "Vegetables", 25, "Purple eggplant"},
		{"Red Onion", "Vegetables", 40, "Purple red onions"},
		{"Leek", "Vegetables", 61, "Fresh leeks"},

		// Proteins
		{"Chicken Breast", "Protein", 165, "Boneless chicken breast"},
		{"Chicken Thighs", "Protein", 209, "Juicy chicken thighs"},
		{"Ground Beef", "Protein", 254, "Lean ground beef"},
		{"Beef Steak", "Protein", 271, "Premium beef steak"},
		{"Pork Chops", "Protein", 231, "Tender pork chops"},
		{"Ground Pork", "Protein", 297, "Fresh ground pork"},
		{"Salmon", "Protein", 208, "Atlantic salmon fillet"},
		{"Tuna", "Protein", 144, "Fresh tuna steak"},
		{"Shrimp", "Protein", 99, "Large shrimp"},
		{"Cod", "Protein", 82, "White cod fillet"},
		{"Eggs", "Protein", 155, "Large chicken eggs"},
		{"Tofu", "Protein", 76, "Firm tofu"},
		{"Black Beans", "Protein", 132, "Canned black beans"},
		{"Chickpeas", "Protein", 164, "Canned chickpeas"},
		{"Lentils", "Protein", 116, "Red lentils"},
		{"Turkey", "Protein", 189, "Ground turkey"},
		{"Bacon", "Protein", 541, "Crispy bacon strips"},
		{"Ham", "Protein", 145, "Sliced ham"},
		{"Sausage", "Protein", 301, "Italian sausage"},

		// Grains & Starches
		{"Rice", "Grains", 130, "Long grain white rice"},
		{"Brown Rice", "Grains", 123, "Whole grain brown rice"},
		{"Pasta", "Grains", 131, "Dried pasta"},
		{"Bread", "Grains", 265, "Whole wheat bread"},
		{"Quinoa", "Grains", 120, "Cooked quinoa"},
		{"Oats", "Grains", 68, "Rolled oats"},
		{"Flour", "Grains", 364, "All-purpose flour"},
		{"Couscous", "Grains", 112, "Pearl couscous"},
		{"Barley", "Grains", 123, "Pearl barley"},
		{"Noodles", "Grains", 138, "Egg noodles"},

		// Dairy
		{"Milk", "Dairy", 61, "Whole milk"},
		{"Cheese", "Dairy", 113, "Cheddar cheese"},
		{"Mozzarella", "Dairy", 280, "Fresh mozzarella"},
		{"Parmesan", "Dairy", 110, "Grated parmesan"},
		{"Greek Yogurt", "Dairy", 97, "Plain Greek yogurt"},
		{"Butter", "Dairy", 717, "Unsalted butter"},
		{"Heavy Cream", "Dairy", 340, "Heavy whipping cream"},
		{"Cream Cheese", "Dairy", 342, "Philadelphia cream cheese"},
		{"Sour Cream", "Dairy", 193, "Regular sour cream"},

		// Seasonings & Condiments
		{"Salt", "Seasonings", 0, "Table salt"},
		{"Black Pepper", "Seasonings", 251, "Ground black pepper"},
		{"Olive Oil", "Seasonings", 884, "Extra virgin olive oil"},
		{"Vegetable Oil", "Seasonings", 884, "Neutral cooking oil"},
		{"Soy Sauce", "Seasonings", 8, "Low sodium soy sauce"},
		{"Vinegar", "Seasonings", 18, "White vinegar"},
		{"Lemon Juice", "Seasonings", 22, "Fresh lemon juice"},
		{"Lime Juice", "Seasonings", 25, "Fresh lime juice"},
		{"Honey", "Seasonings", 304, "Pure honey"},
		{"Maple Syrup", "Seasonings", 260, "Pure maple syrup"},
		{"Ketchup", "Seasonings", 112, "Tomato ketchup"},
		{"Mustard", "Seasonings", 66, "Dijon mustard"},
		{"Mayonnaise", "Seasonings", 680, "Regular mayonnaise"},
		{"Hot Sauce", "Seasonings", 12, "Tabasco sauce"},

		// Herbs & Spices
		{"Basil", "Herbs & Spices", 22, "Fresh basil leaves"},
		{"Oregano", "Herbs & Spices", 265, "Dried oregano"},
		{"Thyme", "Herbs & Spices", 101, "Fresh thyme"},
		{"Rosemary", "Herbs & Spices", 131, "Fresh rosemary"},
		{"Parsley", "Herbs & Spices", 36, "Fresh parsley"},
		{"Cilantro", "Herbs & Spices", 23, "Fresh cilantro"},
		{"Paprika", "Herbs & Spices", 282, "Sweet paprika"},
		{"Cumin", "Herbs & Spices", 375, "Ground cumin"},
		{"Chili Powder", "Herbs & Spices", 282, "Chili powder blend"},
		{"Garlic Powder", "Herbs & Spices", 331, "Granulated garlic"},
		{"Onion Powder", "Herbs & Spices", 341, "Granulated onion"},
		{"Red Pepper Flakes", "Herbs & Spices", 318, "Crushed red pepper"},
		{"Bay Leaves", "Herbs & Spices", 313, "Dried bay leaves"},
		{"Cinnamon", "Herbs & Spices", 247, "Ground cinnamon"},
		{"Ginger", "Herbs & Spices", 80, "Fresh ginger root"},

		// Fruits
		{"Lemon", "Fruits", 29, "Fresh lemons"},
		{"Lime", "Fruits", 30, "Fresh limes"},
		{"Apple", "Fruits", 52, "Red apples"},
		{"Banana", "Fruits", 89, "Ripe bananas"},
		{"Orange", "Fruits", 47, "Navel oranges"},
		{"Strawberry", "Fruits", 32, "Fresh strawberries"},
		{"Blueberry", "Fruits", 57, "Fresh blueberries"},
		{"Avocado", "Fruits", 160, "Ripe avocados"},
		{"Pineapple", "Fruits", 50, "Fresh pineapple"},
		{"Mango", "Fruits", 60, "Ripe mango"},

		// Pantry Staples
		{"Chicken Stock", "Pantry", 12, "Low sodium chicken stock"},
		{"Vegetable Stock", "Pantry", 12, "Vegetable broth"},
		{"Canned Tomatoes", "Pantry", 18, "Crushed tomatoes"},
		{"Tomato Paste", "Pantry", 82, "Double concentrated tomato paste"},
		{"Coconut Milk", "Pantry", 230, "Canned coconut milk"},
		{"Baking Powder", "Pantry", 53, "Double acting baking powder"},
		{"Baking Soda", "Pantry", 0, "Sodium bicarbonate"},
		{"Vanilla Extract", "Pantry", 288, "Pure vanilla extract"},
		{"Sugar", "Pantry", 387, "Granulated white sugar"},
		{"Brown Sugar", "Pantry", 380, "Packed brown sugar"},
		{"Cornstarch", "Pantry", 381, "Corn starch for thickening"},
		{"Breadcrumbs", "Pantry", 395, "Plain breadcrumbs"},
	}

	for _, ing := range ingredientsData {
		_, err := db.Exec("INSERT INTO ingredients (name, category, calories_per_100g, description) VALUES (?, ?, ?, ?)",
			ing.name, ing.category, ing.calories, ing.description)
		if err != nil {
			log.Printf("Error adding %s to ingredients table: %v", ing.name, err)
		}
	}

	// Comprehensive English recipes data
	recipesData := []struct {
		name         string
		category     string
		prepTime     int
		cookTime     int
		servings     int
		difficulty   string
		instructions string
		description  string
	}{
		// Breakfast
		{
			"Scrambled Eggs", "Breakfast", 5, 5, 2, "easy",
			"1. Crack eggs into a bowl and whisk\n2. Heat butter in non-stick pan\n3. Pour in eggs and stir gently\n4. Season with salt and pepper\n5. Serve immediately",
			"Fluffy and creamy scrambled eggs",
		},
		{
			"Pancakes", "Breakfast", 10, 15, 4, "medium",
			"1. Mix flour, sugar, baking powder, and salt\n2. Combine milk, eggs, and melted butter\n3. Fold wet ingredients into dry\n4. Cook on griddle until bubbles form\n5. Flip and cook until golden",
			"Classic fluffy pancakes",
		},
		{
			"Oatmeal", "Breakfast", 2, 5, 1, "easy",
			"1. Bring milk to a simmer\n2. Add oats and cook stirring occasionally\n3. Add honey and cinnamon\n4. Top with fresh fruit",
			"Hearty breakfast oatmeal",
		},
		{
			"French Toast", "Breakfast", 10, 10, 4, "medium",
			"1. Whisk eggs, milk, and vanilla\n2. Dip bread slices in mixture\n3. Cook in buttered pan until golden\n4. Serve with maple syrup",
			"Classic French toast",
		},
		{
			"Breakfast Burrito", "Breakfast", 15, 10, 2, "medium",
			"1. Scramble eggs with cheese\n2. Cook bacon until crispy\n3. Sauté potatoes until golden\n4. Assemble in tortilla and roll\n5. Serve with salsa",
			"Hearty breakfast burrito",
		},

		// Lunch
		{
			"Caesar Salad", "Lunch", 15, 0, 4, "easy",
			"1. Chop romaine lettuce\n2. Make dressing with lemon juice, garlic, and parmesan\n3. Toss lettuce with dressing\n4. Top with croutons and more parmesan",
			"Classic Caesar salad",
		},
		{
			"Grilled Chicken Sandwich", "Lunch", 10, 15, 2, "medium",
			"1. Season chicken breast with salt and pepper\n2. Grill until cooked through\n3. Toast bread\n4. Assemble with lettuce, tomato, and mayo",
			"Juicy grilled chicken sandwich",
		},
		{
			"Vegetable Soup", "Lunch", 15, 30, 6, "easy",
			"1. Sauté onions, carrots, and celery\n2. Add vegetable stock and bring to boil\n3. Add remaining vegetables\n4. Simmer until tender\n5. Season with herbs",
			"Hearty vegetable soup",
		},
		{
			"Pasta Salad", "Lunch", 20, 10, 8, "easy",
			"1. Cook pasta according to package directions\n2. Cool completely\n3. Mix with vegetables and dressing\n4. Chill before serving",
			"Fresh pasta salad",
		},
		{
			"Chicken Quesadilla", "Lunch", 15, 10, 2, "medium",
			"1. Cook chicken with spices\n2. Place chicken and cheese on tortilla\n3. Top with another tortilla\n4. Cook until golden and cheese melts\n5. Cut into wedges",
			"Cheesy chicken quesadilla",
		},

		// Dinner
		{
			"Spaghetti Carbonara", "Dinner", 15, 15, 4, "medium",
			"1. Cook pasta until al dente\n2. Fry bacon until crispy\n3. Mix eggs with parmesan\n4. Toss hot pasta with egg mixture\n5. Add bacon and black pepper",
			"Classic Italian carbonara",
		},
		{
			"Chicken Parmesan", "Dinner", 20, 25, 4, "medium",
			"1. Pound chicken thin and season\n2. Bread with flour, egg, and breadcrumbs\n3. Pan fry until golden\n4. Top with sauce and cheese\n5. Bake until cheese melts",
			"Crispy chicken parmesan",
		},
		{
			"Beef Stir Fry", "Dinner", 15, 10, 4, "medium",
			"1. Slice beef thinly against grain\n2. Heat oil in wok\n3. Stir fry beef until browned\n4. Add vegetables and sauce\n5. Serve over rice",
			"Quick and flavorful stir fry",
		},
		{
			"Salmon with Vegetables", "Dinner", 10, 20, 4, "easy",
			"1. Season salmon with salt and pepper\n2. Roast salmon in oven\n3. Steam or roast vegetables\n4. Serve with lemon wedges",
			"Healthy salmon dinner",
		},
		{
			"Chicken Curry", "Dinner", 20, 35, 6, "medium",
			"1. Brown chicken pieces in oil\n2. Sauté onions until soft\n3. Add spices and cook until fragrant\n4. Add coconut milk and simmer\n5. Return chicken to pot and cook until tender",
			"Aromatic chicken curry",
		},
		{
			"Beef Tacos", "Dinner", 15, 15, 4, "easy",
			"1. Brown ground beef with onions\n2. Add spices and cook\n3. Warm tortillas\n4. Fill with beef and toppings\n5. Serve with lime wedges",
			"Flavorful beef tacos",
		},
		{
			"Vegetarian Chili", "Dinner", 20, 40, 8, "medium",
			"1. Sauté onions, peppers, and garlic\n2. Add spices and tomato paste\n3. Add beans, tomatoes, and stock\n4. Simmer until thick\n5. Garnish with cheese and cilantro",
			"Hearty vegetarian chili",
		},
		{
			"Pork Chops with Apples", "Dinner", 15, 25, 4, "medium",
			"1. Season pork chops with salt and pepper\n2. Sear in hot pan until browned\n3. Remove chops and sauté apples\n4. Return chops to pan and finish in oven",
			"Tender pork chops with apples",
		},
		{
			"Shrimp Scampi", "Dinner", 10, 10, 4, "medium",
			"1. Cook pasta until al dente\n2. Sauté garlic in butter and oil\n3. Add shrimp and cook until pink\n4. Add lemon juice and pasta\n5. Toss with parsley",
			"Garlicky shrimp scampi",
		},
		{
			"Chicken Stir Fry", "Dinner", 15, 12, 4, "medium",
			"1. Cut chicken into bite-sized pieces\n2. Heat oil in wok or large pan\n3. Stir fry chicken until cooked\n4. Add vegetables and sauce\n5. Serve over rice or noodles",
			"Quick chicken stir fry",
		},

		// Snacks & Sides
		{
			"Roasted Vegetables", "Side", 15, 30, 6, "easy",
			"1. Cut vegetables into even pieces\n2. Toss with olive oil and seasonings\n3. Roast in hot oven until tender\n4. Serve hot",
			"Colorful roasted vegetables",
		},
		{
			"Garlic Bread", "Side", 10, 12, 8, "easy",
			"1. Mix butter with minced garlic and herbs\n2. Spread on sliced bread\n3. Wrap in foil and bake\n4. Serve warm",
			"Crispy garlic bread",
		},
		{
			"Mashed Potatoes", "Side", 15, 25, 6, "easy",
			"1. Peel and chop potatoes\n2. Boil until tender\n3. Drain and mash with butter and milk\n4. Season with salt and pepper",
			"Creamy mashed potatoes",
		},
		{
			"Coleslaw", "Side", 15, 0, 8, "easy",
			"1. Shred cabbage and carrots\n2. Make dressing with mayo and vinegar\n3. Toss vegetables with dressing\n4. Chill before serving",
			"Crunchy coleslaw",
		},
		{
			"Rice Pilaf", "Side", 10, 25, 6, "medium",
			"1. Sauté rice in oil until lightly toasted\n2. Add stock and bring to boil\n3. Reduce heat and simmer covered\n4. Fluff with fork before serving",
			"Fluffy rice pilaf",
		},
	}

	for _, rec := range recipesData {
		_, err := db.Exec(`INSERT INTO recipes (name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description) 
						  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			rec.name, rec.category, rec.prepTime, rec.cookTime, rec.servings, rec.difficulty, rec.instructions, rec.description)
		if err != nil {
			log.Printf("Error adding %s to recipes table: %v", rec.name, err)
		}
	}

	// Comprehensive recipe-ingredient relationships
	recipeIngredientsData := []struct {
		recipeID     int
		ingredientID int
		quantity     float64
		unit         string
		notes        string
	}{
		// Scrambled Eggs (recipe 1)
		{1, 36, 4, "large", "room temperature"},
		{1, 55, 2, "tbsp", "for cooking"},
		{1, 58, 1, "tsp", "to taste"},
		{1, 59, 0.5, "tsp", "freshly ground"},

		// Pancakes (recipe 2)
		{2, 47, 2, "cups", "all-purpose"},
		{2, 87, 2, "tbsp", "granulated"},
		{2, 84, 2, "tsp", "double-acting"},
		{2, 58, 1, "tsp", ""},
		{2, 50, 1.5, "cups", "whole milk"},
		{2, 36, 2, "large", "beaten"},
		{2, 55, 0.25, "cup", "melted and cooled"},

		// Oatmeal (recipe 3)
		{3, 46, 0.5, "cup", "rolled oats"},
		{3, 50, 1, "cup", "whole milk"},
		{3, 66, 2, "tbsp", "to taste"},
		{3, 79, 0.5, "tsp", "ground"},
		{3, 74, 1, "medium", "sliced"},

		// French Toast (recipe 4)
		{4, 44, 8, "slices", "thick cut"},
		{4, 36, 4, "large", "beaten"},
		{4, 50, 0.5, "cup", "whole milk"},
		{4, 86, 1, "tsp", "pure"},
		{4, 55, 2, "tbsp", "for cooking"},
		{4, 67, 0.25, "cup", "for serving"},

		// Breakfast Burrito (recipe 5)
		{5, 36, 6, "large", "scrambled"},
		{5, 52, 1, "cup", "shredded cheddar"},
		{5, 41, 4, "strips", "cooked crispy"},
		{7, 2, 2, "medium", "diced and roasted"},
		// Note: Assuming tortillas would be ingredient ID that needs to be added

		// Caesar Salad (recipe 6)
		{6, 11, 2, "heads", "chopped"},
		{6, 54, 0.5, "cup", "grated"},
		{6, 64, 0.25, "cup", "fresh squeezed"},
		{6, 3, 3, "cloves", "minced"},
		{6, 60, 0.25, "cup", "extra virgin"},

		// Grilled Chicken Sandwich (recipe 7)
		{7, 26, 2, "pieces", "6 oz each"},
		{7, 44, 4, "slices", "toasted"},
		{7, 11, 4, "leaves", "fresh"},
		{7, 1, 2, "slices", "thick cut"},
		{7, 70, 2, "tbsp", ""},

		// Vegetable Soup (recipe 8)
		{8, 2, 1, "large", "diced"},
		{8, 5, 3, "large", "sliced"},
		{8, 6, 3, "stalks", "chopped"},
		{8, 90, 8, "cups", "low sodium"},
		{8, 7, 2, "medium", "cubed"},
		{8, 16, 2, "cups", "chopped"},
		{8, 58, 1, "tsp", "to taste"},
		{8, 59, 0.5, "tsp", "to taste"},

		// Pasta Salad (recipe 9)
		{9, 43, 1, "lb", "cooked and cooled"},
		{9, 1, 2, "medium", "diced"},
		{9, 4, 1, "large", "diced"},
		{9, 12, 1, "large", "sliced"},
		{9, 60, 0.33, "cup", "extra virgin"},
		{9, 63, 0.25, "cup", "red wine"},

		// Chicken Quesadilla (recipe 10)
		{10, 26, 2, "pieces", "cooked and diced"},
		{10, 52, 1, "cup", "shredded"},
		{10, 2, 0.5, "medium", "diced"},
		{10, 4, 0.5, "medium", "diced"},
		// Note: Tortillas needed

		// Spaghetti Carbonara (recipe 11)
		{11, 43, 1, "lb", "spaghetti"},
		{11, 41, 6, "strips", "chopped"},
		{11, 36, 4, "large", "beaten"},
		{11, 54, 1, "cup", "grated"},
		{11, 59, 1, "tsp", "freshly cracked"},

		// Chicken Parmesan (recipe 12)
		{12, 26, 4, "pieces", "6 oz each"},
		{12, 47, 1, "cup", "for dredging"},
		{12, 36, 2, "large", "beaten"},
		{12, 89, 2, "cups", "Italian seasoned"},
		{12, 53, 1, "cup", "shredded mozzarella"},
		{12, 91, 2, "cups", "marinara sauce"},

		// Beef Stir Fry (recipe 13)
		{13, 29, 1, "lb", "thinly sliced"},
		{13, 4, 2, "large", "julienned"},
		{13, 5, 2, "large", "sliced"},
		{13, 9, 2, "cups", "florets"},
		{13, 62, 2, "tbsp", "for cooking"},
		{13, 61, 3, "tbsp", "low sodium"},
		{13, 42, 2, "cups", "cooked"},

		// Salmon with Vegetables (recipe 14)
		{14, 32, 4, "fillets", "6 oz each"},
		{14, 15, 2, "cups", "trimmed"},
		{14, 5, 2, "large", "sliced"},
		{14, 9, 2, "cups", "florets"},
		{14, 60, 2, "tbsp", "extra virgin"},
		{14, 73, 1, "large", "cut into wedges"},

		// Chicken Curry (recipe 15)
		{15, 27, 2, "lbs", "cut into pieces"},
		{15, 2, 2, "large", "sliced"},
		{15, 3, 4, "cloves", "minced"},
		{15, 80, 1, "tbsp", "fresh grated"},
		{15, 75, 2, "tsp", "ground"},
		{15, 76, 1, "tsp", "ground"},
		{15, 82, 1, "can", "full fat"},

		// Beef Tacos (recipe 16)
		{16, 28, 1, "lb", "ground"},
		{16, 2, 1, "medium", "diced"},
		{16, 76, 2, "tsp", "ground"},
		{16, 77, 1, "tsp", "chili powder"},
		{16, 73, 1, "large", "cut into wedges"},
		// Note: Tortillas, cheese, lettuce for toppings

		// Vegetarian Chili (recipe 17)
		{17, 2, 1, "large", "diced"},
		{17, 4, 2, "large", "diced"},
		{17, 3, 4, "cloves", "minced"},
		{17, 37, 2, "cans", "drained and rinsed"},
		{17, 38, 2, "cans", "drained and rinsed"},
		{17, 92, 1, "can", "crushed"},
		{17, 90, 4, "cups", "vegetable"},

		// Pork Chops with Apples (recipe 18)
		{18, 31, 4, "pieces", "bone-in"},
		{18, 74, 2, "large", "sliced"},
		{18, 55, 2, "tbsp", "for cooking"},
		{18, 79, 1, "tsp", "ground"},
		{18, 58, 1, "tsp", "to taste"},

		// Shrimp Scampi (recipe 19)
		{19, 34, 1, "lb", "peeled and deveined"},
		{19, 43, 1, "lb", "linguine"},
		{19, 3, 6, "cloves", "minced"},
		{19, 55, 4, "tbsp", "unsalted"},
		{19, 60, 2, "tbsp", "extra virgin"},
		{19, 64, 0.25, "cup", "fresh"},
		{19, 72, 0.25, "cup", "chopped fresh"},

		// Chicken Stir Fry (recipe 20)
		{20, 26, 1, "lb", "cut into strips"},
		{20, 4, 2, "large", "julienned"},
		{20, 5, 2, "large", "sliced"},
		{20, 9, 2, "cups", "florets"},
		{20, 62, 2, "tbsp", "for cooking"},
		{20, 61, 3, "tbsp", "low sodium"},
		{20, 50, 2, "cups", "cooked jasmine"},

		// Roasted Vegetables (recipe 21)
		{21, 5, 3, "large", "chunked"},
		{21, 13, 2, "large", "sliced"},
		{21, 4, 2, "large", "chunked"},
		{21, 60, 3, "tbsp", "extra virgin"},
		{21, 58, 1, "tsp", "to taste"},
		{21, 59, 0.5, "tsp", "to taste"},

		// Garlic Bread (recipe 22)
		{22, 44, 1, "loaf", "sliced"},
		{22, 55, 0.5, "cup", "softened"},
		{22, 3, 4, "cloves", "minced"},
		{22, 72, 2, "tbsp", "chopped fresh"},

		// Mashed Potatoes (recipe 23)
		{23, 7, 3, "lbs", "peeled and cubed"},
		{23, 55, 0.5, "cup", "unsalted"},
		{23, 50, 0.5, "cup", "warm whole milk"},
		{23, 58, 1, "tsp", "to taste"},
		{23, 59, 0.5, "tsp", "white pepper"},

		// Coleslaw (recipe 24)
		{24, 19, 1, "head", "shredded"},
		{24, 5, 2, "large", "grated"},
		{24, 70, 0.75, "cup", "mayonnaise"},
		{24, 63, 2, "tbsp", "apple cider"},
		{24, 87, 2, "tbsp", "granulated"},

		// Rice Pilaf (recipe 25)
		{25, 42, 1.5, "cups", "long grain"},
		{25, 89, 3, "cups", "chicken stock"},
		{25, 2, 1, "medium", "finely diced"},
		{25, 60, 2, "tbsp", "extra virgin"},
		{25, 58, 1, "tsp", "to taste"},
	}

	for _, ri := range recipeIngredientsData {
		_, err := db.Exec(`INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, unit, notes) 
						  VALUES (?, ?, ?, ?, ?)`,
			ri.recipeID, ri.ingredientID, ri.quantity, ri.unit, ri.notes)
		if err != nil {
			log.Printf("Error adding recipe_ingredients relationship: %v", err)
		}
	}

	fmt.Println("Enhanced English test data populated successfully!")
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

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next(w, r)

		duration := time.Since(start)
		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, duration)
	}
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	// 食材のカテゴリ件数
	ingRows, err := db.Query(`
		SELECT category, COUNT(*)
		FROM ingredients
		GROUP BY category
		ORDER BY category
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer ingRows.Close()

	var ingCounts []CategoryCount
	for ingRows.Next() {
		var c CategoryCount
		if err := ingRows.Scan(&c.Category, &c.Count); err != nil {
			http.Error(w, "Data scanning error", http.StatusInternalServerError)
			return
		}
		ingCounts = append(ingCounts, c)
	}
	if err := ingRows.Err(); err != nil {
		http.Error(w, "Data scanning error", http.StatusInternalServerError)
		return
	}

	// レシピのカテゴリ件数
	recRows, err := db.Query(`
		SELECT category, COUNT(*)
		FROM recipes
		GROUP BY category
		ORDER BY category
	`)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer recRows.Close()

	var recCounts []CategoryCount
	for recRows.Next() {
		var c CategoryCount
		if err := recRows.Scan(&c.Category, &c.Count); err != nil {
			http.Error(w, "Data scanning error", http.StatusInternalServerError)
			return
		}
		recCounts = append(recCounts, c)
	}
	if err := recRows.Err(); err != nil {
		http.Error(w, "Data scanning error", http.StatusInternalServerError)
		return
	}

	// JSONレスポンス（構造体ベース）
	w.Header().Set("Content-Type", "application/json")
	resp := CategoryCountsResponse{IngredientCategories: ingCounts, RecipeCategories: recCounts}
	json.NewEncoder(w).Encode(resp)
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	type Stats struct {
		TotalIngredients       int            `json:"total_ingredients"`
		TotalRecipes           int            `json:"total_recipes"`
		AvgPrepTime            float32        `json:"avg_prep_time"`
		AvgCookTime            float32        `json:"avg_cook_time"`
		DifficultyDistribution map[string]int `json:"difficulty_distribution"`
	}

	var stats Stats
	stats.DifficultyDistribution = make(map[string]int)

	err := db.QueryRow(`SELECT COUNT(*) FROM ingredients`).Scan(&stats.TotalIngredients)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`SELECT COUNT(*) FROM recipes`).Scan(&stats.TotalRecipes)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`SELECT AVG(prep_time_minutes) FROM recipes`).Scan(&stats.AvgPrepTime)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	err = db.QueryRow(`SELECT AVG(cook_time_minutes) FROM recipes`).Scan(&stats.AvgCookTime)
	if err != nil {
		http.Error(w, "Database scanning error", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		SELECT difficulty, COUNT(*)
		FROM recipes
		GROUP BY difficulty
	`)
	if err != nil {
		log.Print(err)
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var difficulty string
		var count int
		err = rows.Scan(&difficulty, &count)
		if err != nil {
			log.Print(err)
			http.Error(w, "Database scanning error", http.StatusInternalServerError)
			return
		}
		stats.DifficultyDistribution[difficulty] = count
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func main() {
	initDB()
	populateTestData()
	defer db.Close()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("Missing JWT_SECRET attribute in .env file")
	}
	log.Println("Successfully loaded JWT_SECRET from .env file")

	authService := auth.NewAuthService(db, jwtSecret)
	authHandler := auth.NewAuthHandler(authService)

	userService := users.NewUserService(db)
	userHandler := users.NewUserHandler(userService)

	recipesService := recipes.NewRecipesService(db)
	recipesHandler := recipes.NewRecipesHandler(recipesService)

	ingredientsService := ingredients.NewIngredientsService(db)
	ingredientsHandler := ingredients.NewIngredientsHandler(ingredientsService)

	log.Println("Server running on port 8000")

	http.HandleFunc("/api/auth/register", loggingMiddleware(enableCORS(authHandler.RegisterHandler)))
	http.HandleFunc("/api/auth/login", loggingMiddleware(enableCORS(authHandler.LoginHandler)))

	http.HandleFunc("/api/user/profile", loggingMiddleware(enableCORS(authHandler.AuthMiddleware(userHandler.GetProfile))))
	http.HandleFunc("/api/user/liked-recipes", loggingMiddleware(enableCORS(authHandler.AuthMiddleware(userHandler.GetLikedRecipes))))
	http.HandleFunc("/api/user/liked-recipes/add", loggingMiddleware(enableCORS(authHandler.AuthMiddleware(userHandler.AddLikedRecipe))))
	http.HandleFunc("/api/user/liked-recipes/", loggingMiddleware(enableCORS(authHandler.AuthMiddleware(userHandler.RemoveLikedRecipe))))

	http.HandleFunc("/api/recipes", loggingMiddleware(enableCORS(authHandler.OptionalAuthMiddleware(recipesHandler.AllRecipesHandler))))
	http.HandleFunc("/api/recipes/", loggingMiddleware(enableCORS(authHandler.OptionalAuthMiddleware(recipesHandler.RecipeDetailHandler))))
	http.HandleFunc("/api/recipes/find-by-ingredients", loggingMiddleware(enableCORS(authHandler.OptionalAuthMiddleware(recipesHandler.FindRecipesByIngredientsHandler))))
	http.HandleFunc("/api/recipes/shopping-list/", loggingMiddleware(enableCORS(recipesHandler.ShoppingListHandler)))

	http.HandleFunc("/api/ingredients", loggingMiddleware(enableCORS(ingredientsHandler.AllIngredientsHandler)))
	http.HandleFunc("/api/ingredients/", loggingMiddleware(enableCORS(ingredientsHandler.IngredientDetailsHandler)))

	http.HandleFunc("/api/categories", loggingMiddleware(enableCORS(categoriesHandler)))
	http.HandleFunc("/api/stats", loggingMiddleware(enableCORS(statsHandler)))

	log.Fatal(http.ListenAndServe(":8000", nil))
}
