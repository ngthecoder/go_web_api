package stats

type CategoryCount struct {
	Category string `json:"category"`
	Count    int    `json:"count"`
}

type CategoryCountsResponse struct {
	IngredientCategories map[string]int `json:"ingredient_categories"`
	RecipeCategories     map[string]int `json:"recipe_categories"`
}

type Stats struct {
	TotalIngredients       int            `json:"total_ingredients"`
	TotalRecipes           int            `json:"total_recipes"`
	AvgPrepTime            float32        `json:"avg_prep_time"`
	AvgCookTime            float32        `json:"avg_cook_time"`
	DifficultyDistribution map[string]int `json:"difficulty_distribution"`
}
