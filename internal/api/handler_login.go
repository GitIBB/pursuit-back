package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GitIBB/pursuit/internal/auth"
	"github.com/GitIBB/pursuit/internal/database"
)

// handlerLogin handles user login requests
func (cfg *APIConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	type response struct { // struct to hold the response data
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body) // Create a new JSON decoder for the request body
	params := parameters{}             // Create a new instance of the parameters struct
	err := decoder.Decode(&params)     // Decode the request body into the parameters struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email) // Call the GetUserByEmail function to retrieve the user by email
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword) // Check if the provided password matches the stored hashed password
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT( // Create a new JWT token for the user
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create access JWT", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken() // Create a new refresh token for the user
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create refresh token", err)
		return
	}
	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60), // Set the expiration time for the refresh token to 60 days
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save refresh token in database", err)
		return
	}

	// Check if the client is a browser
	if IsBrowser(r) {
		// Set a cookie for browser clients
		http.SetCookie(w, &http.Cookie{
			Name:     "auth-token",
			Value:    accessToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour), // Set the expiration time for the cookie to match the token expiration
			Secure:   true,                      // Set to true if using HTTPS
		})
	}

	respondWithJSON(w, http.StatusOK, response{ // Create a new response instance containing the user data and token
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Username:  user.Username,
			Email:     user.Email,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
