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

	decoder := json.NewDecoder(r.Body) // create a new JSON decoder for request body
	params := parameters{}             // create a new instance of parameters struct
	err := decoder.Decode(&params)     // decode the request body into parameters struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request parameters", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email) // retrieve the user by email
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "User does not exist", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword) // check if provided password matches the stored hashed password
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT( // Create a new JWT token for user
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
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60), // set the expiration time for the refresh token to 60 days
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save refresh token in database", err)
		return
	}

	// check if the client is a browser
	if IsBrowser(r) {
		// set a cookie for browser clients
		http.SetCookie(w, &http.Cookie{
			Name:     "auth-token",
			Value:    accessToken,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(time.Hour), // Set the expiration time for the cookie to match the token expiration
			Secure:   true,                      // set to true if using HTTPS (during production).
			SameSite: http.SameSiteNoneMode,     // allow cross-site usage
		})
	}

	respondWithJSON(w, http.StatusOK, response{ // create a new response instance containing the user data and token
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
