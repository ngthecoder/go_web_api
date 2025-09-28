package ingredients

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type IngredientHandler struct {
	ingredientService *IngredientService
}

func NewIngredientHandler(s *IngredientService) *IngredientHandler {
	return &IngredientHandler{
		ingredientService: s,
	}
}

func (h *IngredientHandler) AllIngredientsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	search := strings.TrimSpace(query.Get("search"))
	category := strings.TrimSpace(query.Get("category"))
	sort := query.Get("sort")
	order := query.Get("order")
	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	if sort != "name" && sort != "calories" {
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

	err, total := h.ingredientService.ingredientsCounter(&w, search, category)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	totalPages := (total + limit - 1) / limit
	hasNext := page < totalPages

	err, ingredients := h.ingredientService.ingredientsRetriever(&w, search, category, sort, order, limit, offset)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ingredients": ingredients,
		"total":       total,
		"page":        page,
		"page_size":   limit,
		"total_pages": totalPages,
		"has_next":    hasNext,
	})
}

func (h *IngredientHandler) IngredientDetailsHandler(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 4 || pathParts[3] == "" {
		http.Error(w, "Invalid URL format. Use /api/ingredients/{id}", http.StatusBadRequest)
		return
	}

	ingredientID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		http.Error(w, "Invalid ingredient ID", http.StatusBadRequest)
		return
	}

	err, ingredient, associatedRecipes := h.ingredientService.ingredientDetailsWithRecipesRetriever(&w, ingredientID)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	response := map[string]interface{}{
		"ingredient": ingredient,
		"recipes":    associatedRecipes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
