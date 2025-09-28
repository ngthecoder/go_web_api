package recipes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RecipeHandler struct {
	recipeService *RecipeService
}

func NewRecipesHandler(s *RecipeService) *RecipeHandler {
	return &RecipeHandler{
		recipeService: s,
	}
}

func (h *RecipeHandler) AllRecipesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))
	category := strings.TrimSpace(query.Get("category"))
	difficulty := strings.TrimSpace(query.Get("difficulty"))
	maxTimeStr := query.Get("max_time")
	sort := query.Get("sort")
	order := query.Get("order")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	var maxTime int
	if maxTimeStr != "" {
		if m, err := strconv.Atoi(maxTimeStr); err == nil && m > 0 {
			maxTime = m
		}
	}

	validSorts := map[string]bool{
		"name": true, "prep_time": true, "cook_time": true,
		"total_time": true, "servings": true, "difficulty": true,
	}
	if sort == "" || !validSorts[sort] {
		sort = "name"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	offset := (page - 1) * limit

	err, total := h.recipeService.recipesCounter(&w, search, category, difficulty, maxTime)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	totalPages := (total + limit - 1) / limit
	hasNext := page < totalPages

	err, recipes := h.recipeService.recipesRetriever(&w, search, category, difficulty, sort, order, maxTime, limit, offset)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recipes":     recipes,
		"total":       total,
		"page":        page,
		"page_size":   limit,
		"total_pages": totalPages,
		"has_next":    hasNext,
	})
}

func (h *RecipeHandler) RecipeDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/recipes/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe id", http.StatusBadRequest)
		return
	}

	err, recipe, ingredients := h.recipeService.recipeDetailsWithIngredientsRetriever(&w, id)

	w.Header().Set("Content-Type", "application/json")
	resp := RecipeWithIngredients{Recipe: recipe, Ingredients: ingredients}
	json.NewEncoder(w).Encode(resp)
}

func (h *RecipeHandler) FindRecipesByIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ingredientsParams := query.Get("ingredients")
	matchType := query.Get("match_type")
	limitParams := query.Get("limit")

	if ingredientsParams == "" {
		http.Error(w, "Missing required parameters: ingredients", http.StatusBadRequest)
		return
	}

	if matchType == "" || matchType == "exact" {
		matchType = "partial"
	}

	limit := 10
	if limitParams != "" {
		if l, err := strconv.Atoi(limitParams); err == nil && l > 0 {
			limit = l
		}
	}

	ingredientIDStrings := strings.Split(ingredientsParams, ",")
	ingredientIDs := make([]int, 0, len(ingredientIDStrings))

	for _, idStr := range ingredientIDStrings {
		if id, err := strconv.Atoi(strings.TrimSpace(idStr)); err == nil {
			ingredientIDs = append(ingredientIDs, id)
		}
	}

	if len(ingredientIDs) == 0 {
		http.Error(w, "Invalid ingredient IDs", http.StatusBadRequest)
		return
	}

	err, matchedRecipes := h.recipeService.matchedRecipesRetriever(&w, matchType, ingredientIDs, limit)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(matchedRecipes)
}

func (h *RecipeHandler) ShoppingListHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] == "" {
		http.Error(w, "Invalid URL format. Use /api/recipes/shopping-list/{id}", http.StatusBadRequest)
		return
	}

	recipeID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	haveIngredientsStr := r.URL.Query().Get("have_ingredients")
	haveIngredientIDs := make(map[int]struct{})
	if haveIngredientsStr != "" {
		ids := strings.Split(haveIngredientsStr, ",")
		for _, idStr := range ids {
			id, err := strconv.Atoi(idStr)
			if err == nil {
				haveIngredientIDs[id] = struct{}{}
			}
		}
	}

	err, shoppingList := h.recipeService.shoppingListRetriever(&w, recipeID, haveIngredientIDs)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"recipe_id":     recipeID,
		"shopping_list": shoppingList,
	})
}
