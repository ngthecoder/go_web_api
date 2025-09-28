package ingredients

type Ingredient struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Calories    int    `json:"calories_per_100g"`
	Description string `json:"description"`
}
