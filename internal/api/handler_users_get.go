package api

import (
	"net/http"
	"strings"

	"github.com/GitIBB/pursuit/internal/auth"
)

func (cfg *APIConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization") // Get the Authorization header from the request
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is missing", nil)
		return
	}

	// Extract the token from the Bearer <token> format
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		respondWithError(w, http.StatusUnauthorized, "Invalid Authorization header format", nil)
		return
	}

	token := tokenParts[1] // Extract the token from the header
	if token == "" {
		respondWithError(w, http.StatusUnauthorized, "Token is missing", nil)
		return
	}

	// Validate the JWT token and extract the user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret) // Validate the JWT token and extract the user ID
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Retrieve the user from the database using the user ID
	user, err := cfg.db.GetUserByID(r.Context(), userID) // Call the GetUser function to retrieve the user by ID
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, user)

}
