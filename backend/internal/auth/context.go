package auth

import (
	"context"
	"net/http"
)

// contextKey is an unexported type for context keys to prevent collisions
// with keys defined in other packages.
type contextKey string

// userIDKey is the context key for the user ID.
const userIDKey contextKey = "user_id"

// SetUserIDInContext returns a new context with the user ID set.
func SetUserIDInContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext retrieves the user ID from the context.
// Returns empty string if not found.
func GetUserIDFromContext(ctx context.Context) string {
	if userID, ok := ctx.Value(userIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserIDFromRequest is a convenience function to get user ID from request context.
// Returns empty string if not found.
func GetUserIDFromRequest(r *http.Request) string {
	return GetUserIDFromContext(r.Context())
}
