package users

import (
	"database/sql"

	"github.com/ngthecoder/go_web_api/internal/auth"
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
	err := s.db.QueryRow("SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1", userID).
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
		WHERE ulr.user_id = $1;
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
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM recipes WHERE id = $1)", recipeID).Scan(&exists)
	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}
	if !exists {
		return errors.NewNotFoundError("Recipe not found")
	}

	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_liked_recipes WHERE user_id = $1 AND recipe_id = $2)", userID, recipeID).Scan(&exists)
	if err != nil {
		return errors.NewInternalServerError("Database error", err)
	}
	if exists {
		return errors.NewConflictError("Recipe already in liked list")
	}

	_, err = s.db.Exec("INSERT INTO user_liked_recipes (user_id, recipe_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP);", userID, recipeID)
	if err != nil {
		return errors.NewInternalServerError("Failed to add liked recipe", err)
	}

	return nil
}

func (s *UserService) removeLikedRecipe(userID string, recipeID int) error {
	removeRecipeQuery := `DELETE FROM user_liked_recipes WHERE user_id = $1 AND recipe_id = $2;`
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

func (s *UserService) updateUserProfile(userID string, username, email string) (UserProfile, error) {
	var exists bool
	err := s.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE (username = $1 OR email = $2) AND id != $3)",
		username, email, userID,
	).Scan(&exists)

	if err != nil {
		return UserProfile{}, errors.NewInternalServerError("Database error", err)
	}
	if exists {
		return UserProfile{}, errors.NewConflictError("Username or email already taken")
	}

	_, err = s.db.Exec(
		"UPDATE users SET username = $1, email = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3",
		username, email, userID,
	)
	if err != nil {
		return UserProfile{}, errors.NewInternalServerError("Failed to update profile", err)
	}

	return s.getUserProfile(userID)
}

func (s *UserService) changePassword(userID string, currentPassword, newPassword string) error {
	var storedEncodedPasswordHash string
	err := s.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&storedEncodedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFoundError("User not found")
		}
		return errors.NewInternalServerError("Database error", err)
	}

	verified, err := auth.VerifyPasswordHash(currentPassword, storedEncodedPasswordHash)
	if err != nil {
		return errors.NewInternalServerError("Failed to verify password", err)
	}
	if !verified {
		return errors.NewBadRequestError("Current password is incorrect")
	}

	newEncodedPasswordHash, err := auth.HashPassword(newPassword)
	if err != nil {
		return errors.NewInternalServerError("Failed to hash password", err)
	}

	_, err = s.db.Exec(
		"UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2",
		newEncodedPasswordHash, userID,
	)
	if err != nil {
		return errors.NewInternalServerError("Failed to update password", err)
	}

	return nil
}

func (s *UserService) deleteAccount(userID string, password string) error {
	var storedEncodedPasswordHash string
	err := s.db.QueryRow("SELECT password_hash FROM users WHERE id = $1", userID).Scan(&storedEncodedPasswordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.NewNotFoundError("User not found")
		}
		return errors.NewInternalServerError("Database error", err)
	}

	verified, err := auth.VerifyPasswordHash(password, storedEncodedPasswordHash)
	if err != nil {
		return errors.NewInternalServerError("Failed to verify password", err)
	}
	if !verified {
		return errors.NewBadRequestError("Password is incorrect")
	}

	_, err = s.db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return errors.NewInternalServerError("Failed to delete account", err)
	}

	return nil
}
