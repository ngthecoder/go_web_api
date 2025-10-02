package users

import "time"

type UserProfile struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type LikedRecipeRequest struct {
	RecipeID int `json:"recipe_id"`
}

type LikedRecipeResponse struct {
	RecipeIDs []int `json:"recipe_ids"`
}
type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}
type DeleteAccountRequest struct {
	Password string `json:"password"`
}
