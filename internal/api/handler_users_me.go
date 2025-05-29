package api

import (
	"net/http"

	"github.com/GitIBB/pursuit/internal/auth"
)

func (cfg *APIConfig) handlerMe(w http.ResponseWriter, r *http.Request) {
	// Extract the token from the cookie
	cookie, err := r.Cookie("auth-token")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Auth token cookie is missing", nil)
		return
	}

	// Validate the JWT token and extract the user ID
	userID, err := auth.ValidateJWT(cookie.Value, cfg.jwtSecret) //
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
		return
	}

	// Retrieve the user from the database using the user ID
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	// Respond with the user's information
	respondWithJSON(w, http.StatusOK, user)
}
