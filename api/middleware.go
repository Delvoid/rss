package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/Delvoid/go_rss/internal/database"
)

// ErrNoAPIKey is returned when no API key is found
// in the Request
var (
	ErrNoAPIKey      = "no API key provided"
	ErrInvalidFormat = "invalid Authorization header format"
	ErrInvalidAPI    = "invalid API key"
)

func (cfg *APIConfig) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("Authorization")
		if apiKey == "" {
			RespondWithError(w, ErrNoAPIKey, http.StatusUnauthorized)
			return
		}

		// The Authorization header should be in the format: "ApiKey <actual-api-key>"
		parts := strings.Split(apiKey, " ")
		if len(parts) != 2 || parts[0] != "ApiKey" {
			RespondWithError(w, ErrInvalidFormat, http.StatusUnauthorized)
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), parts[1])
		if err != nil {
			RespondWithError(w, ErrInvalidAPI, http.StatusUnauthorized)
			return
		}

		// Add the user to the request context
		ctx := r.Context()
		ctx = WithUser(ctx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// WithUser adds a user to the context
func WithUser(ctx context.Context, user database.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(r *http.Request) (database.User, bool) {
	user, ok := r.Context().Value(userContextKey).(database.User)
	return user, ok
}

type contextKey string

const userContextKey contextKey = "user"
