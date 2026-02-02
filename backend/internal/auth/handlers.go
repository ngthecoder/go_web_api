package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ngthecoder/go_web_api/internal/errors"
)

type AuthHandler struct {
	service *AuthService
}

func NewAuthHandler(service *AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

func (h *AuthHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			errors.WriteHTTPError(w, errors.NewUnauthorizedError("Authorization header missing"))
			return
		}

		claims, err := h.service.validateJWT(authHeader)
		if err != nil {
			errors.WriteHTTPError(w, errors.NewUnauthorizedError("Invalid or expired token"))
			return
		}

		ctx := SetUserIDInContext(r.Context(), claims.UserID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func (h *AuthHandler) OptionalAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next(w, r)
			return
		}

		claims, err := h.service.validateJWT(authHeader)
		if err != nil {
			next(w, r)
			return
		}

		ctx := SetUserIDInContext(r.Context(), claims.UserID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	var request RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid JSON"))
		return
	}

	if request.Username == "" || request.Email == "" || request.Password == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Username, email, and password are required"))
		return
	}

	response, err := h.service.registerUser(request)
	if err != nil {
		if err == ErrUserExists {
			errors.WriteHTTPError(w, errors.NewConflictError("User already exists"))
			return
		}
		errors.WriteHTTPError(w, errors.NewInternalServerError("Registration failed", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteHTTPError(w, errors.NewMethodNotAllowedError())
		return
	}

	var request LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Invalid JSON"))
		return
	}

	if request.Email == "" || request.Password == "" {
		errors.WriteHTTPError(w, errors.NewBadRequestError("Email and password are required"))
		return
	}

	response, err := h.service.loginUser(request)
	if err != nil {
		if err == ErrInvalidCredentials {
			errors.WriteHTTPError(w, errors.NewUnauthorizedError("Invalid email or password"))
			return
		}
		errors.WriteHTTPError(w, errors.NewInternalServerError("Login failed", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
