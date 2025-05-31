package api

import (
	"encoding/json"
	"net/http"

	"github.com/GitIBB/pursuit/internal/auth"
	"github.com/GitIBB/pursuit/internal/database"
	"github.com/google/uuid"
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

	// Retrieve the user ID from the context
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: missing user ID", nil)
		return
	}

	decoder := json.NewDecoder(r.Body) // new JSON decoder for request body
	params := parameters{}             // new instance of parameters struct
	err := decoder.Decode(&params)     // decode request body into parameters struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password) // hash password
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{ // update user data in db
		ID:             userID,
		Email:          params.Email,
		Username:       params.Username,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{ // respond with updated user data
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
		},
	})
}
