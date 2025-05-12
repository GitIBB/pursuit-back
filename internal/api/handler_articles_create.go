package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GitIBB/pursuit/internal/auth"
	"github.com/GitIBB/pursuit/internal/database"
	"github.com/google/uuid"
)

type Article struct { // struct to hold article data
	ID        uuid.UUID `json:"id"`
	CreatedAT time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
}

// Handler function to create a new article
func (cfg *APIConfig) handlerArticlesCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header) // Extract the token from the request header
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret) // Validate the JWT token
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body) // Create a new JSON decoder for the request body
	params := parameters{}             // Create a new instance of the parameters struct
	err = decoder.Decode(&params)      // Decode the request body into the parameters struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request parameters", err)
		return
	}

	cleaned, err := validateArticle(params.Body) // Validate the article body
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to validate request parameters", err)
	}

	// Create a new instance of the CreateArticleParams struct
	article, err := cfg.db.CreateArticle(r.Context(), database.CreateArticleParams{
		UserID: userID,
		Body:   cleaned,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create article", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Article{ // Create a new response instance containing the article data
		ID:        article.ID,
		CreatedAT: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		UserID:    article.UserID,
		Title:     params.Title,
		Body:      cleaned,
	})
}

func validateArticle(body string) (string, error) { // Validate the article body
	// Check for empty body
	if len(body) == 0 {
		return "", errors.New("article body cannot be empty")
	}
	// Check for length
	// Set a maximum length for the article body
	const maxLength = 15000
	if len(body) > maxLength {
		return "", errors.New("article is too long, please limit to 15,000 characters")
	}
	return body, nil

}
