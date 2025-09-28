package users

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ngthecoder/go_web_api/internal/errors"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(s *UserService) *UserHandler {
	return &UserHandler{
		userService: s,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	userProfile, err := h.userService.getUserProfile(userID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	json.NewEncoder(w).Encode(userProfile)
}

func (h *UserHandler) GetLikedRecipes(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	recipesList, err := h.userService.getLikedRecipes(userID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	json.NewEncoder(w).Encode(recipesList)
}

func (h *UserHandler) AddLikedRecipe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 4 || pathParts[3] == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid URL format. Use /api/user/liked-recipes/{id}"))
		return
	}
	recipeID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid recipe ID format"))
		return
	}

	err = h.userService.addLikedRecipe(userID, recipeID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe added to liked list"})
}

func (h *UserHandler) RemoveLikedRecipe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 4 || pathParts[3] == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid URL format. Use /api/user/liked-recipes/{id}"))
		return
	}
	recipeID, err := strconv.Atoi(pathParts[3])
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid recipe ID format"))
		return
	}

	err = h.userService.removeLikedRecipe(userID, recipeID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe removed from liked list"})
}
