package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (cfg *APIConfig) handlerUploads(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form (10 MB max)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Failed to parse form data", err)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil && err != http.ErrMissingFile {
		respondWithError(w, http.StatusBadRequest, "Failed to retrieve image file", err)
		return
	}
	if file == nil {
		respondWithError(w, http.StatusBadRequest, "No image file provided", nil)
		return
	}
	defer file.Close()

	// Validate the file type (read first 512 bytes)
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to read image file", err)
		return
	}
	// Reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset file pointer", err)
		return
	}

	mimeType := http.DetectContentType(buffer)
	if mimeType != "image/jpeg" && mimeType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid image format. Only JPEG and PNG are valid", nil)
		return
	}

	// Validate file extension
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if !allowedExtensions[ext] {
		respondWithError(w, http.StatusBadRequest, "Invalid file extension. Only .jpg, .jpeg, and .png are valid", nil)
		return
	}

	// Ensure uploads directory exists in project root
	projectRoot, err := filepath.Abs(filepath.Join(".", "..", ".."))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to resolve project root", err)
		return
	}
	uploadsDir := filepath.Join(projectRoot, "uploads")
	os.MkdirAll(uploadsDir, os.ModePerm)
	imagePath := filepath.Join(uploadsDir, handler.Filename)

	dst, err := os.Create(imagePath)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save image", err)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to save image", err)
		return
	}

	// Respond with the file URL
	url := "/api/uploads/" + handler.Filename
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf(`{"url":"%s"}`, url)))
}
