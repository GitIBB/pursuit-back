package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/GitIBB/pursuit/internal/database"
	"github.com/google/uuid"
)

type Article struct { // struct to hold article data
	ID        uuid.UUID   `json:"id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	UserID    uuid.UUID   `json:"user_id"`
	Title     string      `json:"title"`
	Body      ArticleBody `json:"body"`
	ImageUrl  string      `json:"image_url"`
}

type ArticleBody struct { // struct to hold article body data
	Introduction string `json:"introduction"`
	MainBody     string `json:"main_body"`
	End          string `json:"end"`
}

// Handler function to create a new article
func (cfg *APIConfig) handlerArticlesCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string      `json:"title"`
		ArticleBody ArticleBody `json:"article_body"`
	}

	// Retrieve the user ID from the context
	userID, ok := r.Context().Value("userID").(uuid.UUID)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized: missing user ID", nil)
		return
	}

	// Parse the form data (for image upload)
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10MB
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form data", err)
		return
	}

	// Get the image file from the form
	file, handler, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		respondWithError(w, http.StatusBadRequest, "Failed to retrieve image file", err)
		return
	}
	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	// image type checking logic
	var imageURL string
	if file != nil {
		// Validate the file type
		buffer := make([]byte, 512) // Read the first 512 bytes of the file
		_, err := file.Read(buffer)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to read image file", err)
			return
		}

		// Reset the file pointer after reading
		_, err = file.Seek(0, 0)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to reset file pointer", err)
			return
		}

		// Detect the MIME type
		mimeType := http.DetectContentType(buffer)
		if mimeType != "image/jpeg" && mimeType != "image/png" {
			respondWithError(w, http.StatusBadRequest, "Invalid image format. Only JPEG and PNG are valid", err)
			return
		}

		// Validate the file extension
		allowedExtensions := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		// Get the file extension
		fileExtension := handler.Filename[len(handler.Filename)-4:] // Get the last 4 characters of the filename
		if !allowedExtensions[fileExtension] {
			respondWithError(w, http.StatusBadRequest, "Invalid file extension. Only .jpg, .jpeg and .png are valid", err)
			return
		}

		// Save the image to the uploads folder
		imagePath := "github.com/pursuit/uploads/" + handler.Filename
		dst, err := os.Create(imagePath)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save image", err)
			return
		}
		defer dst.Close()

		// Copy the file content to the destination
		_, err = io.Copy(dst, file)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to save image", err)
			return

		}
		// Set the image URL to the saved path
		imageURL = imagePath // Set the image URL to the saved path
	}

	decoder := json.NewDecoder(r.Body) // Create a new JSON decoder for the request body
	params := parameters{}             // Create a new instance of the parameters struct
	err = decoder.Decode(&params)      // Decode the request body into the parameters struct
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
		UserID: userID,
		Title:  params.Title,
		Body:   bodyJSON, // save the marshaled JSON body
		ImageUrl: sql.NullString{
			String: imageURL,
			Valid:  imageURL != "", // Set Valid to true if imageURL is not empty
		},
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create article", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Article{ // Create a new response instance containing the article data
		ID:        article.ID,
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
		UserID:    article.UserID,
		Title:     params.Title,
		Body:      cleanedBody,
		ImageUrl:  imageURL,
	})
}

func validateArticle(body ArticleBody) (ArticleBody, error) { // Validate the article body
	// Check for empty fields in the article body
	if len(body.Introduction) == 0 {
		return body, errors.New("introduction cannot be empty")
	}
	if len(body.MainBody) == 0 {
		return body, errors.New("main body cannot be empty")
	}
	if len(body.End) == 0 {
		return body, errors.New("end cannot be empty")
	}

	// alternatively, set a max length here

	return body, nil // Return the validated article body
}
