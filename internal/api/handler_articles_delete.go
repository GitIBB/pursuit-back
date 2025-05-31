package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) handlerArticlesDelete(w http.ResponseWriter, r *http.Request) {

	// Extract the article ID from the URL
	articleIDString := r.PathValue("articleID")
	articleID, err := uuid.Parse(articleIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid article ID", err)
		return
	}

	// Retrieve the user ID from the context
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: missing user ID", nil)
		return
	}

	// Retrieve the article from the database using the article ID
	dbArticle, err := cfg.db.GetArticle(r.Context(), articleID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Article not found", err)
		return
	}
	// Check if the user ID from the token matches the user ID of the article
	if dbArticle.UserID != userID {
		respondWithError(w, http.StatusForbidden, "You can not delete this article", nil)
		return
	}

	// Delete the article from the database
	err = cfg.db.DeleteArticle(r.Context(), articleID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete article", err)
		return
	}

	w.WriteHeader(http.StatusNoContent) // Respond with 204 No Content
}
