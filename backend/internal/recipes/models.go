package recipes

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

type IngredientWithQuantity struct {
	IngredientID int     `json:"ingredient_id"`
	Name         string  `json:"name"`
	Quantity     float64 `json:"quantity"`
	Unit         string  `json:"unit"`
	Notes        string  `json:"notes"`
}

type RecipeWithIngredients struct {
	Recipe      Recipe                   `json:"recipe"`
	Ingredients []IngredientWithQuantity `json:"ingredients"`
}

type MatchedRecipe struct {
	ID                      int     `json:"id"`
	Name                    string  `json:"name"`
	Category                string  `json:"category"`
	PrepTimeMinutes         int     `json:"prep_time_minutes"`
	CookTimeMinutes         int     `json:"cook_time_minutes"`
	Servings                int     `json:"servings"`
	Difficulty              string  `json:"difficulty"`
	Instructions            string  `json:"instructions"`
	Description             string  `json:"description"`
	MatchedIngredientsCount int     `json:"matched_ingredients_count"`
	TotalIngredientsCount   int     `json:"total_ingredients_count"`
	MatchScore              float32 `json:"match_score"`
}
