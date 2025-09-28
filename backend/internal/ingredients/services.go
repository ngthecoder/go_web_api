package ingredients

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ngthecoder/go_web_api/internal/recipes"
)

type IngredientService struct {
	db *sql.DB
}

func NewIngredientService(db *sql.DB) *IngredientService {
	return &IngredientService{
		db: db,
	}
}

func (s *IngredientService) ingredientsCounter(w *http.ResponseWriter, search, category string) (error, int) {
	sqlCountQuery, args := s.buildIngredientCountQuery(search, category)

	var total int
	err := s.db.QueryRow(sqlCountQuery, args...).Scan(&total)
	if err != nil {
		log.Printf("Counting Query: %s", sqlCountQuery)
		http.Error(*w, "Data scanning error", http.StatusInternalServerError)
		return errors.New("Data scanning error"), 0
	}

	return nil, total
}

func (s *IngredientService) ingredientsRetriever(w *http.ResponseWriter, search, category, sort, order string, limit, offset int) (error, []Ingredient) {
	sqlQuery, args := s.buildIngredientQuery(search, category, sort, order, limit, offset)

	rows, err := s.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Query: %s", sqlQuery)
		http.Error(*w, "Database error", http.StatusInternalServerError)
		return errors.New("Database error"), nil
	}
	defer rows.Close()

	var ingredients []Ingredient
	for rows.Next() {
		var ingredient Ingredient
		err = rows.Scan(&ingredient.ID, &ingredient.Name, &ingredient.Category, &ingredient.Calories, &ingredient.Description)
		if err != nil {
			http.Error(*w, "Data scanning error", http.StatusInternalServerError)
			return errors.New("Data scanning error"), nil
		}
		ingredients = append(ingredients, ingredient)
	}

	return nil, ingredients
}

func (s *IngredientService) buildIngredientCountQuery(search, category string) (string, []interface{}) {
	query := "SELECT COUNT(*) FROM ingredients"
	conditions := []string{}
	args := []interface{}{}

	if search != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	return query, args
}

func (s *IngredientService) buildIngredientQuery(search, category, sort, order string, limit, offset int) (string, []interface{}) {
	query := "SELECT id, name, category, calories_per_100g, description FROM ingredients"
	conditions := []string{}
	args := []interface{}{}

	if search != "" {
		conditions = append(conditions, "(name LIKE ? OR description LIKE ?)")
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm, searchTerm)
	}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	orderByClause := ""
	switch sort {
	case "calories":
		orderByClause = "calories_per_100g"
	default:
		orderByClause = "name"
	}

	query += " ORDER BY " + orderByClause + " " + strings.ToUpper(order)
	query += " LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	return query, args
}

func (s *IngredientService) ingredientDetailsWithRecipesRetriever(w *http.ResponseWriter, ingredientID int) (error, Ingredient, []recipes.Recipe) {
	query := `
        SELECT
            i.id, i.name, i.category, i.calories_per_100g, i.description,
            ri.quantity, ri.unit, ri.notes,
            r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes,
            r.servings, r.difficulty, r.instructions, r.description
        FROM
            ingredients AS i
        LEFT JOIN
            recipe_ingredients AS ri ON i.id = ri.ingredient_id
        LEFT JOIN
            recipes AS r ON ri.recipe_id = r.id
        WHERE
            i.id = ?;
    `

	rows, err := s.db.Query(query, ingredientID)
	if err != nil {
		http.Error(*w, "Database query error", http.StatusInternalServerError)
		return errors.New("Database query error"), Ingredient{}, nil
	}
	defer rows.Close()

	var ingredient Ingredient
	var associatedRecipes []recipes.Recipe
	foundIngredient := false

	for rows.Next() {
		var (
			iID, calories                                     int
			iName, iCategory, iDesc                           string
			quantity                                          sql.NullFloat64
			unit, notes                                       sql.NullString
			rID, prepTime, cookTime, servings                 sql.NullInt64
			rName, rCategory, difficulty, instructions, rDesc sql.NullString
		)

		err := rows.Scan(
			&iID, &iName, &iCategory, &calories, &iDesc,
			&quantity, &unit, &notes,
			&rID, &rName, &rCategory, &prepTime, &cookTime,
			&servings, &difficulty, &instructions, &rDesc,
		)
		if err != nil {
			log.Println("Scan error:", err)
			continue
		}

		if !foundIngredient {
			ingredient = Ingredient{
				ID:          iID,
				Name:        iName,
				Category:    iCategory,
				Calories:    calories,
				Description: iDesc,
			}
			foundIngredient = true
		}

		if rID.Valid {
			associatedRecipes = append(associatedRecipes, recipes.Recipe{
				ID:              int(rID.Int64),
				Name:            rName.String,
				Category:        rCategory.String,
				PrepTimeMinutes: int(prepTime.Int64),
				CookTimeMinutes: int(cookTime.Int64),
				Servings:        int(servings.Int64),
				Difficulty:      difficulty.String,
				Instructions:    instructions.String,
				Description:     rDesc.String,
			})
		}
	}

	if !foundIngredient {
		http.Error(*w, "Ingredient not found", http.StatusNotFound)
		return errors.New("Ingredient not found"), Ingredient{}, nil
	}

	return nil, ingredient, associatedRecipes
}
