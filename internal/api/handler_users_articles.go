package api

import (
	"net/http"

	"github.com/GitIBB/pursuit/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) handlerUserArticles(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Optional: support pagination
	limit := 100 // or parse from query
	offset := 0  // or parse from query

	articles, err := cfg.db.GetArticlesByUserId(r.Context(), database.GetArticlesByUserIdParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user's articles", err)
		return
	}

	respondWithJSON(w, http.StatusOK, articles)
}
