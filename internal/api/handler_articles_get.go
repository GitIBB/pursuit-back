package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/GitIBB/pursuit/internal/database"
	"github.com/google/uuid"
)

func (cfg *APIConfig) handlerArticlesGet(w http.ResponseWriter, r *http.Request) {
	articleIDString := r.PathValue("articleID")   // Extract the articleID from the URL path
	articleID, err := uuid.Parse(articleIDString) // Parse the articleID string to a UUID
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid article ID", err)
		return
	}

	// Use the GetArticle function to retrieve the article by its ID
	dbArticle, err := cfg.db.GetArticle(r.Context(), articleID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve article", err)
		return
	}

	// Fetch the username using GetUserByID
	user, err := cfg.db.GetUserByID(r.Context(), dbArticle.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user", err)
		return
	}

	// After fetching dbArticle
	category, err := cfg.db.GetCategoryByID(r.Context(), dbArticle.CategoryID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve category", err)
		return
	}

	var body ArticleBody                        // Initialize an empty ArticleBody struct
	err = json.Unmarshal(dbArticle.Body, &body) // Unmarshal the article body from JSON into the struct
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to unmarshal article body", err)
		return
	}

	// Handle sql.NullString for ImageUrl
	imageUrl := ""
	if dbArticle.ImageUrl.Valid {
		imageUrl = dbArticle.ImageUrl.String
	}

	respondWithJSON(w, http.StatusOK, Article{
		ID:        dbArticle.ID,
		CreatedAt: dbArticle.CreatedAt,
		UpdatedAt: dbArticle.UpdatedAt,
		UserID:    dbArticle.UserID,
		Title:     dbArticle.Title,
		Body:      body,
		ImageUrl:  imageUrl,
		Username:  user.Username,
		Category:  category.Name,
	})
}

func (cfg *APIConfig) handlerArticlesRetrieve(w http.ResponseWriter, r *http.Request) { // Handler function to retrieve all articles with pagination
	// parse query parameters for pagination
	page := r.URL.Query().Get("page")   // Get the page query parameter from the URL
	limit := r.URL.Query().Get("limit") // Get the limit query parameter from the URL

	// Default values for pagination
	pageNum := 1
	limitNum := 10

	// Convert query parameters to integers
	if page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			pageNum = p // Set pageNum to the parsed value if it's a valid positive integer
		}
	}
	if limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			limitNum = l // <-- Use the parsed value!
		}
	}

	//calculate the offset
	offset := (pageNum - 1) * limitNum // Calculate the offset for pagination

	// Fetch the total number of articles
	totalRecords, err := cfg.db.GetTotalArticlesCount(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve total articles count", err)
		return
	}

	// Calculate total pages
	totalPages := (int(totalRecords) + limitNum - 1) / limitNum // Round up

	// Create getArticlesParams struct to pass to GetArticles function
	params := database.GetArticlesParams{
		Limit:  int32(limitNum),
		Offset: int32(offset),
	}

	dbArticles, err := cfg.db.GetArticles(r.Context(), params) // Retrieve all articles from the database
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve articles", err)
		return
	}

	authorID := uuid.Nil                           // Initialize authorID to uuid.Nil
	authorIDString := r.URL.Query().Get("user_id") // Get the author_id query parameter from the URL
	if authorIDString != "" {                      // If author_id is provided
		authorID, err = uuid.Parse(authorIDString) // Parse the author_id string to a UUID
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid author ID", err)
			return
		}
	}

	articles := []Article{}                // Initialize an empty slice of Article
	for _, dbArticle := range dbArticles { // Iterate over the retrieved articles
		if authorID != uuid.Nil && dbArticle.UserID != authorID { // Check if the article's author ID matches the provided author ID
			continue // If not, skip this article
		}

		// Unmarshal the article body from JSON into the ArticleBody struct
		var body ArticleBody
		err = json.Unmarshal(dbArticle.Body, &body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to unmarshal article body", err)
			return
		}

		// Fetch the category using GetCategoryByID
		category, err := cfg.db.GetCategoryByID(r.Context(), dbArticle.CategoryID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to retrieve category", err)
			return
		}

		// Handle sql.NullString for ImageUrl
		imageUrl := ""
		if dbArticle.ImageUrl.Valid {
			imageUrl = dbArticle.ImageUrl.String
		}

		articles = append(articles, Article{
			ID:        dbArticle.ID,
			CreatedAt: dbArticle.CreatedAt,
			UpdatedAt: dbArticle.UpdatedAt,
			UserID:    dbArticle.UserID,
			Title:     dbArticle.Title,
			Body:      body,
			ImageUrl:  imageUrl,
			Username:  dbArticle.Username,
			Category:  category.Name,
		})
	}
	// Create the response with metadata
	type Metadata struct {
		CurrentPage  int `json:"current_page"`
		TotalPages   int `json:"total_pages"`
		TotalRecords int `json:"total_records"`
	}
	type Response struct {
		Metadata Metadata  `json:"metadata"`
		Articles []Article `json:"articles"`
	}

	respondWithJSON(w, http.StatusOK, Response{
		Metadata: Metadata{
			CurrentPage:  pageNum,
			TotalPages:   totalPages,
			TotalRecords: int(totalRecords),
		},
		Articles: articles,
	})
}
