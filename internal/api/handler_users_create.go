package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GitIBB/pursuit/internal/auth"
	"github.com/GitIBB/pursuit/internal/database"
	"github.com/google/uuid"
)

type User struct { // struct to hold user data
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
}

func (cfg *APIConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct { // struct to hold the request parameters
		Password string `json:"password"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}
	type response struct { // struct to hold the response data
		User
	}

	decoder := json.NewDecoder(r.Body) // Create a new JSON decoder for the request body
	params := parameters{}             // Create a new instance of the parameters struct
	err := decoder.Decode(&params)     // Decode the request body into the parameters struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request parameters", err)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password) // Hash the password using the HashPassword function
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to hash password", err)
		return
	}

	createUserParams := database.CreateUserParams{ // Create a new instance of the CreateUserParams struct
		Email:          params.Email,
		Username:       params.Username,
		HashedPassword: hashedPassword,
	}

	user, err := cfg.db.CreateUser(r.Context(), createUserParams) // Call the CreateUser function to create a new user
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{ // Create a new response instance containing the user data
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
			Username:  user.Username,
		},
	})
}
