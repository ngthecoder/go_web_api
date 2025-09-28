package recipes

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type RecipeService struct {
	db *sql.DB
}

func NewRecipeService(db *sql.DB) *RecipeService {
	return &RecipeService{
		db: db,
	}
}

func (s *RecipeService) recipesCounter(w *http.ResponseWriter, search, category, difficulty string, maxTime int) (int, error) {
	sqlCountQuery, args := s.buildRecipeCountQuery(search, category, difficulty, maxTime)

	var total int
	err := s.db.QueryRow(sqlCountQuery, args...).Scan(&total)
	if err != nil {
		http.Error(*w, "Data scanning error", http.StatusInternalServerError)
		return 0, err
	}

	return total, nil
}

func (s *RecipeService) recipesRetriever(w *http.ResponseWriter, search, category, difficulty, sort, order string, maxTime, limit, offset int) ([]Recipe, error) {
	sqlQuery, args := s.buildRecipeQuery(search, category, difficulty, sort, order, maxTime, limit, offset)

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		http.Error(*w, "Database error", http.StatusInternalServerError)
		return nil, errors.New("Database error")
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe
		err = rows.Scan(&recipe.ID, &recipe.Name, &recipe.Category, &recipe.PrepTimeMinutes, &recipe.CookTimeMinutes, &recipe.Servings, &recipe.Difficulty, &recipe.Instructions, &recipe.Description)
		if err != nil {
			http.Error(*w, "Data scanning error", http.StatusInternalServerError)
			return nil, errors.New("Data scanning error")
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (s *RecipeService) buildRecipeCountQuery(search, category, difficulty string, maxTime int) (string, []interface{}) {
	query := "SELECT COUNT(*) FROM recipes"
	conditions := []string{}
	args := []interface{}{}

	if search != "" {
		conditions = append(conditions, "(name LIKE ? OR instructions LIKE ? OR description LIKE ?)")
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if difficulty != "" {
		conditions = append(conditions, "difficulty = ?")
		args = append(args, difficulty)
	}

	if maxTime > 0 {
		conditions = append(conditions, "(prep_time_minutes + cook_time_minutes) <= ?")
		args = append(args, maxTime)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}

func (s *RecipeService) buildRecipeQuery(search, category, difficulty, sort, order string, maxTime, limit, offset int) (string, []interface{}) {
	query := "SELECT id, name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description FROM recipes"
	conditions := []string{}
	args := []interface{}{}

	if search != "" {
		conditions = append(conditions, "(name LIKE ? OR instructions LIKE ? OR description LIKE ?)")
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
	}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if difficulty != "" {
		conditions = append(conditions, "difficulty = ?")
		args = append(args, difficulty)
	}

	if maxTime > 0 {
		conditions = append(conditions, "(prep_time_minutes + cook_time_minutes) <= ?")
		args = append(args, maxTime)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	orderByClause := ""
	switch sort {
	case "prep_time":
		orderByClause = "prep_time_minutes"
	case "cook_time":
		orderByClause = "cook_time_minutes"
	case "total_time":
		orderByClause = "(prep_time_minutes + cook_time_minutes)"
	case "servings":
		orderByClause = "servings"
	case "difficulty":
		orderByClause = "difficulty"
	default:
		orderByClause = "name"
	}

	query += " ORDER BY " + orderByClause + " " + strings.ToUpper(order)
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	return query, args
}

func (s *RecipeService) recipeDetailsWithIngredientsRetriever(w *http.ResponseWriter, id int) (Recipe, []IngredientWithQuantity, error) {
	var recipe Recipe
	err := s.db.QueryRow(`
		SELECT id, name, category, prep_time_minutes, cook_time_minutes, servings, difficulty, instructions, description 
		FROM recipes WHERE id = ?`, id).
		Scan(&recipe.ID, &recipe.Name, &recipe.Category, &recipe.PrepTimeMinutes, &recipe.CookTimeMinutes,
			&recipe.Servings, &recipe.Difficulty, &recipe.Instructions, &recipe.Description)
	if err == sql.ErrNoRows {
		http.Error(*w, "Recipe not found", http.StatusNotFound)
		return Recipe{}, nil, errors.New("Recipe not found")
	} else if err != nil {
		http.Error(*w, "Database error", http.StatusInternalServerError)
		return Recipe{}, nil, errors.New("Database error")
	}

	rows, err := s.db.Query(`
		SELECT i.id, i.name, ri.quantity, ri.unit, ri.notes
		FROM recipe_ingredients ri
		JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = ?`, id)
	if err == sql.ErrNoRows {
		http.Error(*w, "Recipe ingredients not found", http.StatusNotFound)
		return Recipe{}, nil, errors.New("Recipe ingredients not found")
	} else if err != nil {
		http.Error(*w, "Database error", http.StatusInternalServerError)
		return Recipe{}, nil, errors.New("Database error")
	}
	defer rows.Close()

	var ingredients []IngredientWithQuantity
	for rows.Next() {
		var ing IngredientWithQuantity
		err := rows.Scan(&ing.IngredientID, &ing.Name, &ing.Quantity, &ing.Unit, &ing.Notes)
		if err != nil {
			http.Error(*w, "Data scanning error", http.StatusInternalServerError)
			return Recipe{}, nil, errors.New("Data scanning error")
		}
		ingredients = append(ingredients, ing)
	}
	if err = rows.Err(); err != nil {
		http.Error(*w, "Data scanning error", http.StatusInternalServerError)
		return Recipe{}, nil, errors.New("Data scanning error")
	}

	return recipe, ingredients, nil
}

func (s *RecipeService) matchedRecipesRetriever(w *http.ResponseWriter, matchType string, ingredientIDs []int, limit int) ([]MatchedRecipe, error) {
	sqlQuery := ""
	args := []interface{}{}

	placeholders := make([]string, 0, len(ingredientIDs))
	for _, ingredientID := range ingredientIDs {
		placeholders = append(placeholders, "?")
		args = append(args, ingredientID)
	}
	args = append(args, limit)

	if matchType == "partial" {
		sqlQuery = fmt.Sprintf(
			`
			SELECT r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description,
			COUNT(ri.ingredient_id) as match_ingredients_count,
			(SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id) as total_ingredients_count
			FROM recipes r
			JOIN recipe_ingredients ri on r.id = ri.recipe_id
			WHERE ri.ingredient_id in (%s)
			GROUP BY r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description
			ORDER BY match_ingredients_count DESC, total_ingredients_count ASC
			LIMIT ?
		`, strings.Join(placeholders, ","))
	}

	if matchType == "exact" {
		sqlQuery = fmt.Sprintf(
			`
			SELECT r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description,
			COUNT(ri.ingredient_id) as match_ingredients_count,
			(SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id) as total_ingredients_count
			FROM recipes r
			JOIN recipe_ingredients ri on r.id = ri.recipe_id
			WHERE ri.ingredient_id in (%s)
			GROUP BY r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description
			HAVING COUNT(ri.ingredient_id) = (SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id)
			ORDER BY match_ingredients_count DESC, total_ingredients_count ASC
			LIMIT ?
		`, strings.Join(placeholders, ","))
	}

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("SQL Error: %v", err)
		http.Error(*w, "Database error", http.StatusInternalServerError)
		return nil, errors.New("Database error")
	}
	defer rows.Close()

	matchedRecipes := []MatchedRecipe{}
	for rows.Next() {
		var matchedRecipe MatchedRecipe
		err = rows.Scan(
			&matchedRecipe.ID,
			&matchedRecipe.Name,
			&matchedRecipe.Category,
			&matchedRecipe.PrepTimeMinutes,
			&matchedRecipe.CookTimeMinutes,
			&matchedRecipe.Servings,
			&matchedRecipe.Difficulty,
			&matchedRecipe.Instructions,
			&matchedRecipe.Description,
			&matchedRecipe.MatchedIngredientsCount,
			&matchedRecipe.TotalIngredientsCount,
		)
		matchedRecipe.MatchScore = float32(matchedRecipe.MatchedIngredientsCount) / float32(matchedRecipe.TotalIngredientsCount)
		if err != nil {
			http.Error(*w, "Database scanning error", http.StatusInternalServerError)
			return nil, errors.New("Database scanning error")
		}
		matchedRecipes = append(matchedRecipes, matchedRecipe)
	}

	return matchedRecipes, nil
}

func (s *RecipeService) shoppingListRetriever(w *http.ResponseWriter, recipeID int, haveIngredientIDs map[int]struct{}) ([]IngredientWithQuantity, error) {
	query := `
		SELECT
			ri.ingredient_id,
			i.name,
			ri.quantity,
			ri.unit,
			ri.notes
		FROM
			recipe_ingredients AS ri
		JOIN
			ingredients AS i ON ri.ingredient_id = i.id
		WHERE
			ri.recipe_id = ?;
	`

	rows, err := s.db.Query(query, recipeID)
	if err != nil {
		http.Error(*w, "Database query error", http.StatusInternalServerError)
		return nil, errors.New("Database query error")
	}
	defer rows.Close()

	var shoppingList []IngredientWithQuantity
	var recipeFound bool
	for rows.Next() {
		recipeFound = true
		var ingredient IngredientWithQuantity
		err := rows.Scan(
			&ingredient.IngredientID,
			&ingredient.Name,
			&ingredient.Quantity,
			&ingredient.Unit,
			&ingredient.Notes,
		)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		if _, ok := haveIngredientIDs[ingredient.IngredientID]; !ok {
			shoppingList = append(shoppingList, ingredient)
		}
	}

	if !recipeFound {
		http.Error(*w, "Recipe not found", http.StatusNotFound)
		return nil, errors.New("Recipe not found")
	}

	return shoppingList, nil
}
