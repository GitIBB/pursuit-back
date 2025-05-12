package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) handlerArticlesGet(w http.ResponseWriter, r *http.Request) {
	articleIDString := r.PathValue("articleID")
	articleID, err := uuid.Parse(articleIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid article ID", err)
		return
	}
	dbArticle, err := cfg.db.GetArticle(r.Context(), articleID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve article", err)
		return
	}
	respondWithJSON(w, http.StatusOK, Article{
		ID:        dbArticle.ID,
		CreatedAT: dbArticle.CreatedAt,
		UpdatedAt: dbArticle.UpdatedAt,
		UserID:    dbArticle.UserID,
		Title:     dbArticle.Title,
		Body:      dbArticle.Body,
	})
}

func (cfg *APIConfig) handlerArticlesRetrieve(w http.ResponseWriter, r *http.Request) {
	dbArticles, err := cfg.db.GetArticles(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve articles", err)
		return
	}
	articles := []Article{}
	for _, dbArticle := range dbArticles {
		articles = append(articles, Article{
			ID:        dbArticle.ID,
			CreatedAT: dbArticle.CreatedAt,
			UpdatedAt: dbArticle.UpdatedAt,
			UserID:    dbArticle.UserID,
			Title:     dbArticle.Title,
			Body:      dbArticle.Body,
		})
	}
	respondWithJSON(w, http.StatusOK, articles)
}
