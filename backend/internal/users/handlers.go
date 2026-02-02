package users

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ngthecoder/go_web_api/internal/auth"
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
	userID := auth.GetUserIDFromRequest(r)

	userProfile, err := h.userService.getUserProfile(userID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userProfile)
}

func (h *UserHandler) GetLikedRecipes(w http.ResponseWriter, r *http.Request) {
	userID := auth.GetUserIDFromRequest(r)

	recipesList, err := h.userService.getLikedRecipes(userID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipesList)
}

func (h *UserHandler) AddLikedRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	userID := auth.GetUserIDFromRequest(r)

	var request LikedRecipeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid JSON"))
		return
	}

	err = h.userService.addLikedRecipe(userID, request.RecipeID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe added to liked list"})
}

func (h *UserHandler) RemoveLikedRecipe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	userID := auth.GetUserIDFromRequest(r)

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 5 || pathParts[4] == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid URL format"))
		return
	}

	recipeID, err := strconv.Atoi(pathParts[4])
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid recipe ID"))
		return
	}

	err = h.userService.removeLikedRecipe(userID, recipeID)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Recipe removed from liked list"})
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	userID := auth.GetUserIDFromRequest(r)

	var request UpdateProfileRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid JSON"))
		return
	}

	if request.Username == "" || request.Email == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Username and email are required"))
		return
	}

	updatedProfile, err := h.userService.updateUserProfile(userID, request.Username, request.Email)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedProfile)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	userID := auth.GetUserIDFromRequest(r)

	var request ChangePasswordRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid JSON"))
		return
	}

	if request.CurrentPassword == "" || request.NewPassword == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Current password and new password are required"))
		return
	}

	if len(request.NewPassword) < 6 {
		errors.WriteHTTPError(w, errors.NewBadRequestError("New password must be at least 6 characters"))
		return
	}

	err = h.userService.changePassword(userID, request.CurrentPassword, request.NewPassword)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Password updated successfully"})
}

func (h *UserHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	userID := auth.GetUserIDFromRequest(r)

	var request DeleteAccountRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid JSON"))
		return
	}

	if request.Password == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Password is required"))
		return
	}

	err = h.userService.deleteAccount(userID, request.Password)
	if err != nil {
		errors.WriteHTTPError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted successfully"})
}
