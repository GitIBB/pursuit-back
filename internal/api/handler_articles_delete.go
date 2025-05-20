package api

import (
	"net/http"

	"github.com/GitIBB/pursuit/internal/auth"
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

	// Extract the token from the Authorization header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find JWT", err)
	}

	// Validate the JWT token and extract the user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid JWT", err)
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
