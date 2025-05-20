package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/GitIBB/pursuit/internal/auth"
)

func (cfg *APIConfig) middlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get the token from the "Authorization" header
		authHeader := r.Header.Get("Authorization")
		var token string

		if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			// If no Authorization header, try to get the token from the "auth-token" cookie
			cookie, err := r.Cookie("auth-token")
			if err != nil {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}
			token = cookie.Value
		}

		// Validate the token
		userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
		if err != nil {
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		// Add the user ID to the request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
