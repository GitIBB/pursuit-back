package api

import (
	"encoding/json"
	"net/http"

	"github.com/GitIBB/pursuit/internal/auth"
	"github.com/GitIBB/pursuit/internal/database"
)

func (cfg *APIConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header) // Get the bearer token from the request
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret) // Validate the JWT token
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Failed to validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body) // Create a new JSON decoder for the request body
	params := parameters{}             // Create a new instance of the parameters struct
	err = decoder.Decode(&params)      // Decode the request body into the parameters struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password) // Hash the password
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          params.Email,
		Username:       params.Username,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{ // Respond with the updated user data
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
		},
	})
}
