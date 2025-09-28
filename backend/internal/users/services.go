package users

import (
	"database/sql"

	"github.com/ngthecoder/go_web_api/internal/errors"
	"github.com/ngthecoder/go_web_api/internal/recipes"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (s *UserService) getUserProfile(userID string) (UserProfile, error) {
	var userProfile UserProfile
	err := s.db.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = ?", userID).
		Scan(&userProfile.ID, &userProfile.Username, &userProfile.Email, &userProfile.CreatedAt, &userProfile.UpdatedAt)

	if err == sql.ErrNoRows {
		return UserProfile{}, errors.NewNotFoundError("User not found")
	}
	if err != nil {
		return UserProfile{}, errors.NewInternalServerError("Database scanning error", err)
	}

	return userProfile, nil
}

func (s *UserService) getLikedRecipes(userID string) ([]recipes.Recipe, error) {
	query := `
		SELECT r.id, r.name, r.category, r.prep_time_minutes, r.cook_time_minutes, r.servings, r.difficulty, r.instructions, r.description
		FROM user_liked_recipes ulr
		JOIN recipes r on ulr.recipe_id = r.id
		WHERE ulr.user_id = ?;
	`
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, errors.NewInternalServerError("Database error", err)
	}
	defer rows.Close()

	var recipesList []recipes.Recipe
	for rows.Next() {
		var recipe recipes.Recipe
		err := rows.Scan(&recipe.ID, &recipe.Name, &recipe.Category,
			&recipe.PrepTimeMinutes, &recipe.CookTimeMinutes,
			&recipe.Servings, &recipe.Difficulty,
			&recipe.Instructions, &recipe.Description)
		if err != nil {
			return nil, errors.NewInternalServerError("Data scanning error", err)
		}
		recipesList = append(recipesList, recipe)
	}

	return recipesList, nil
}

func (s *UserService) addLikedRecipe(userID string, recipeID int) error {
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM recipes WHERE id = ?)", recipeID).Scan(&exists)
	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}
	if !exists {
		return errors.NewNotFoundError("Recipe not found")
	}

	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_liked_recipes WHERE user_id = ? AND recipe_id = ?)", userID, recipeID).Scan(&exists)
	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}
	if exists {
		return errors.NewConflictError("Recipe already in liked list")
	}

	_, err = s.db.Exec("INSERT INTO user_liked_recipes (user_id, recipe_id, created_at) VALUES (?, ?, datetime('now'));", userID, recipeID)
	if err != nil {
		return errors.NewInternalServerError("Failed to add liked recipe", err)
	}

	return nil
}

func (s *UserService) removeLikedRecipe(userID string, recipeID int) error {
	removeRecipeQuery := `DELETE FROM user_liked_recipes WHERE user_id = ? AND recipe_id = ?;`
	result, err := s.db.Exec(removeRecipeQuery, userID, recipeID)
	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.NewInternalServerError("Failed to check affected rows", err)
	}
	if rowsAffected == 0 {
		return errors.NewNotFoundError("Recipe not in liked list")
	}

	return nil
}
