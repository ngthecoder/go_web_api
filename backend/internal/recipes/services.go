package recipes

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/ngthecoder/go_web_api/internal/errors"
)

type RecipesService struct {
	db *sql.DB
}

func NewRecipesService(db *sql.DB) *RecipesService {
	return &RecipesService{
		db: db,
	}
}

func (s *RecipesService) recipesCounter(search, category, difficulty string, maxTime int) (int, error) {
	sqlCountQuery, args := s.buildRecipeCountQuery(search, category, difficulty, maxTime)

	var total int
	err := s.db.QueryRow(sqlCountQuery, args...).Scan(&total)
	if err != nil {
		return 0, errors.NewInternalServerError("Data scanning error", err)
	}

	return total, nil
}

func (s *RecipesService) recipesRetriever(search, category, difficulty, sort, order string, maxTime, limit, offset int, userID string) ([]Recipe, error) {
	sqlQuery, args := s.buildRecipeQuery(search, category, difficulty, sort, order, maxTime, limit, offset, userID)

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var recipe Recipe

		err = rows.Scan(
			&recipe.ID,
			&recipe.Name,
			&recipe.Category,
			&recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes,
			&recipe.Servings,
			&recipe.Difficulty,
			&recipe.Instructions,
			&recipe.Description,
			&recipe.IsLiked,
		)
		if err != nil {
			return nil, errors.NewInternalServerError("Data scanning error", err)
		}

		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

func (s *RecipesService) recipeDetailsWithIngredientsRetriever(id int, userID string) (Recipe, []IngredientWithQuantity, error) {
	var recipe Recipe

	query := `
		SELECT 
			r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, 
			r.servings, r.difficulty, r.instructions, r.description,
			CASE WHEN ulr.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked
		FROM recipes r
		LEFT JOIN user_liked_recipes ulr ON r.id = ulr.recipe_id AND ulr.user_id = $1
		WHERE r.id = $2
	`

	err := s.db.QueryRow(query, userID, id).
		Scan(&recipe.ID, &recipe.Name, &recipe.Category, &recipe.PrepTimeMinutes, &recipe.CookTimeMinutes,
			&recipe.Servings, &recipe.Difficulty, &recipe.Instructions, &recipe.Description,
			&recipe.IsLiked)

	if err == sql.ErrNoRows {
		return Recipe{}, nil, errors.NewNotFoundError("Recipe not found")
	} else if err != nil {
		return Recipe{}, nil, errors.NewInternalServerError("Database error", err)
	}

	rows, err := s.db.Query(`
		SELECT i.id, i.name, ri.quantity, ri.unit, ri.notes
		FROM recipe_ingredients ri
		JOIN ingredients i ON ri.ingredient_id = i.id
		WHERE ri.recipe_id = $1`, id)
	if err != nil {
		return Recipe{}, nil, errors.NewInternalServerError("Database error", err)
	}
	defer rows.Close()

	var ingredients []IngredientWithQuantity
	for rows.Next() {
		var ing IngredientWithQuantity
		err := rows.Scan(&ing.IngredientID, &ing.Name, &ing.Quantity, &ing.Unit, &ing.Notes)
		if err != nil {
			return Recipe{}, nil, errors.NewInternalServerError("Data scanning error", err)
		}
		ingredients = append(ingredients, ing)
	}
	if err = rows.Err(); err != nil {
		return Recipe{}, nil, errors.NewInternalServerError("Data scanning error", err)
	}

	return recipe, ingredients, nil
}

func (s *RecipesService) matchedRecipesRetriever(matchType string, ingredientIDs []int, limit int, userID string) ([]MatchedRecipe, error) {
	sqlQuery := ""
	args := []interface{}{}

	placeholderNum := 2

	placeholders := make([]string, 0, len(ingredientIDs))
	for _, ingredientID := range ingredientIDs {
		placeholders = append(placeholders, "$"+strconv.Itoa(placeholderNum))
		args = append(args, ingredientID)
		placeholderNum++
	}

	args = append(args, userID)
	args = append(args, limit)

	if matchType == "partial" {
		sqlQuery = fmt.Sprintf(
			`
			SELECT 
				r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, 
				r.servings, r.difficulty, r.instructions, r.description,
				COUNT(ri.ingredient_id) as match_ingredients_count,
				(SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id) as total_ingredients_count,
				CASE WHEN ulr.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked
			FROM recipes r
			JOIN recipe_ingredients ri on r.id = ri.recipe_id
			LEFT JOIN user_liked_recipes ulr ON r.id = ulr.recipe_id AND ulr.user_id = $1
			WHERE ri.ingredient_id in (%s)
			GROUP BY r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description, ulr.user_id
			ORDER BY match_ingredients_count DESC, total_ingredients_count ASC
			LIMIT $%d
		`, strings.Join(placeholders, ","), placeholderNum)
	}

	if matchType == "exact" {
		sqlQuery = fmt.Sprintf(
			`
			SELECT 
				r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, 
				r.servings, r.difficulty, r.instructions, r.description,
				COUNT(ri.ingredient_id) as match_ingredients_count,
				(SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id) as total_ingredients_count,
				CASE WHEN ulr.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked
			FROM recipes r
			JOIN recipe_ingredients ri on r.id = ri.recipe_id
			LEFT JOIN user_liked_recipes ulr ON r.id = ulr.recipe_id AND ulr.user_id = $1
			WHERE ri.ingredient_id in (%s)
			GROUP BY r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description, ulr.user_id
			HAVING COUNT(ri.ingredient_id) = (SELECT COUNT(*) FROM recipe_ingredients WHERE recipe_id = r.id)
			ORDER BY match_ingredients_count DESC, total_ingredients_count ASC
			LIMIT $%d
		`, strings.Join(placeholders, ","), placeholderNum)
	}

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
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
			&matchedRecipe.IsLiked,
		)

		matchedRecipe.MatchScore = float32(matchedRecipe.MatchedIngredientsCount) / float32(matchedRecipe.TotalIngredientsCount)

		if err != nil {
			return nil, errors.NewInternalServerError("Database scanning error", err)
		}
		matchedRecipes = append(matchedRecipes, matchedRecipe)
	}

	return matchedRecipes, nil
}

func (s *RecipesService) shoppingListRetriever(recipeID int, haveIngredientIDs map[int]struct{}) ([]IngredientWithQuantity, error) {
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
			ri.recipe_id = $1;
	`

	rows, err := s.db.Query(query, recipeID)
	if err != nil {
		return nil, errors.NewInternalServerError("Database query error", err)
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
			continue
		}

		if _, ok := haveIngredientIDs[ingredient.IngredientID]; !ok {
			shoppingList = append(shoppingList, ingredient)
		}
	}

	if !recipeFound {
		return nil, errors.NewNotFoundError("Recipe not found")
	}

	return shoppingList, nil
}

func (s *RecipesService) buildRecipeCountQuery(search, category, difficulty string, maxTime int) (string, []interface{}) {
	query := "SELECT COUNT(*) FROM recipes"
	conditions := []string{}
	args := []interface{}{}

	placeholderNum := 1

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(name LIKE $%d OR instructions LIKE $%d OR description LIKE $%d)", placeholderNum, placeholderNum+1, placeholderNum+2))
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
		placeholderNum += 3
	}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", placeholderNum))
		args = append(args, category)
		placeholderNum++
	}

	if difficulty != "" {
		conditions = append(conditions, fmt.Sprintf("difficulty = $%d", placeholderNum))
		args = append(args, difficulty)
		placeholderNum++
	}

	if maxTime > 0 {
		conditions = append(conditions, fmt.Sprintf("(prep_time_minutes + cook_time_minutes) <= $%d", placeholderNum))
		args = append(args, maxTime)
		placeholderNum++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}

func (s *RecipesService) buildRecipeQuery(search, category, difficulty, sort, order string, maxTime, limit, offset int, userID string) (string, []interface{}) {
	placeholderNum := 1

	query := fmt.Sprintf("SELECT r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description, CASE WHEN ulr.user_id IS NOT NULL THEN 1 ELSE 0 END as is_liked FROM recipes r LEFT JOIN user_liked_recipes ulr ON r.id = ulr.recipe_id AND ulr.user_id = $%d", placeholderNum)
	placeholderNum++

	conditions := []string{}
	args := []interface{}{userID}

	if search != "" {
		conditions = append(conditions, fmt.Sprintf("(r.name LIKE $%d OR r.instructions LIKE $%d OR r.description LIKE $%d)", placeholderNum, placeholderNum+1, placeholderNum+2))
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm, searchTerm)
		placeholderNum += 3
	}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("r.category = $%d", placeholderNum))
		args = append(args, category)
		placeholderNum++
	}

	if difficulty != "" {
		conditions = append(conditions, fmt.Sprintf("r.difficulty = $%d", placeholderNum))
		args = append(args, difficulty)
		placeholderNum++
	}

	if maxTime > 0 {
		conditions = append(conditions, fmt.Sprintf("(r.prep_time_minutes + r.cook_time_minutes) <= $%d", placeholderNum))
		args = append(args, maxTime)
		placeholderNum++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	orderByClause := ""
	switch sort {
	case "prep_time":
		orderByClause = "r.prep_time_minutes"
	case "cook_time":
		orderByClause = "r.cook_time_minutes"
	case "total_time":
		orderByClause = "(r.prep_time_minutes + r.cook_time_minutes)"
	case "servings":
		orderByClause = "r.servings"
	case "difficulty":
		orderByClause = "r.difficulty"
	default:
		orderByClause = "r.name"
	}

	query += " ORDER BY " + orderByClause + " " + strings.ToUpper(order)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", placeholderNum, placeholderNum+1)
	args = append(args, limit, offset)

	return query, args
}
