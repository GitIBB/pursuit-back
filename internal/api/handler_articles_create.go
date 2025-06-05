package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/GitIBB/pursuit/internal/database"
	"github.com/google/uuid"
)

type Article struct { // struct to hold article data
	ID        uuid.UUID   `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	UserID    uuid.UUID   `json:"user_id"`
	Category  string      `json:"category"` // Category of the article, can be retrieved from the database category table
	Title     string      `json:"title"`
	Body      ArticleBody `json:"body"`
	ImageUrl  string      `json:"image_url"`
	Username  string      `json:"username"` // Username of the author, can be retrieved from the database user table
}

type ArticleBody struct { // struct to hold article body data
	Headers map[string]string          `json:"headers"`
	Content map[string]json.RawMessage `json:"content"`
	Images  map[string]string          `json:"images"`
}

// Handler function to create a new article
func (cfg *APIConfig) handlerArticlesCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string      `json:"title"`
		ArticleBody ArticleBody `json:"article_body"`
		ImageUrl    string      `json:"image_url"`
		CategoryID  uuid.UUID   `json:"category_id"`
	}

	// Retrieve the user ID from the context
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: missing user ID", nil)
		return
	}

	decoder := json.NewDecoder(r.Body) // Create a new JSON decoder for the request body
	params := parameters{}             // Create a new instance of the parameters struct
	err := decoder.Decode(&params)     // Decode the request body into the parameters struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to decode request parameters", err)
		return
	}

	cleanedBody, err := validateArticle(params.ArticleBody) // Validate the article body
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to validate request parameters", err)
	}

	bodyJSON, err := json.Marshal(cleanedBody) // Marshal the cleaned body to JSON
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal article body", err)
		return
	}

	// Save the article to the database
	article, err := cfg.db.CreateArticle(r.Context(), database.CreateArticleParams{
		UserID:     userID,            // save the user ID from the context
		CategoryID: params.CategoryID, // save the category ID
		Title:      params.Title,      // save the title
		Body:       bodyJSON,          // save the marshaled JSON body
		ImageUrl: sql.NullString{
			String: params.ImageUrl,       // save the image URL
			Valid:  params.ImageUrl != "", // Set Valid to true if imageURL is not empty
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create article", err)
		return
	}

	// retrieve username from the database
	user, err := cfg.db.GetUserByID(r.Context(), article.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch username", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Article{ // Create a new response instance containing the article data
		ID:        article.ID,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		UserID:    article.UserID,
		Title:     params.Title,
		Body:      cleanedBody,
		ImageUrl:  params.ImageUrl,
		Username:  user.Username, // Username is retrieved from the database
	})
}

func validateArticle(body ArticleBody) (ArticleBody, error) { // Function to validate the article body
	if body.Content == nil {
		return body, errors.New("content is required")
	}
	requiredSections := []string{"introduction", "mainBody", "conclusion"}
	for _, section := range requiredSections {
		if _, ok := body.Content[section]; !ok {
			return body, errors.New(section + " is required in content")
		}
	}
	return body, nil
}
