package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *APIConfig) handlerCategoriesGet(w http.ResponseWriter, r *http.Request) {
	// Retrieve the categories from the database
	categories, err := cfg.db.GetCategories(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve categories", err)
		return
	}

	type CategoryResponse struct {
		ID   uuid.UUID `json:"id"`
		Name string    `json:"name"`
	}

	catResp := make([]CategoryResponse, len(categories))
	for i, cat := range categories {
		catResp[i] = CategoryResponse{
			ID:   cat.ID,
			Name: cat.Name,
		}

	}

	// Respond with the categories in JSON format
	respondWithJSON(w, http.StatusOK, catResp)
}
