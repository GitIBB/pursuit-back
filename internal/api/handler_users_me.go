package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) handlerMe(w http.ResponseWriter, r *http.Request) {
	// Retrieve the user ID from the context
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: missing user ID", nil)
		return
	}

	// retrieve user from pursuitdb using userID
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	// response with user data
	respondWithJSON(w, http.StatusOK, user)
}
