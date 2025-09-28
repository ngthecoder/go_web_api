package ingredients

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ngthecoder/go_web_api/internal/errors"
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

	total, err := h.ingredientService.ingredientsCounter(search, category)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	totalPages := (total + limit - 1) / limit
	hasNext := page < totalPages

	ingredients, err := h.ingredientService.ingredientsRetriever(search, category, sort, order, limit, offset)
	if err != nil {
		errors.WriteHTTPError(w, err)
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
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid URL format. Use /api/ingredients/{id}"))
		return
	}

	ingredientID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid ingredient ID"))
		return
	}

	ingredient, associatedRecipes, err := h.ingredientService.ingredientDetailsWithRecipesRetriever(ingredientID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	response := map[string]interface{}{
		"ingredient": ingredient,
		"recipes":    associatedRecipes,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
