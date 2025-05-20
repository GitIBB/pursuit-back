package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) { // Function to respond with an error message in JSON format
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}

	if isProduction() && code > 499 { // In production, do not expose error details
		msg = "Internal Server Error"
	}

	type errorResponse struct { // Struct to hold the error message
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{ // Create a JSON response with the error message
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) { // Function to respond with a JSON payload
	data, err := json.Marshal(payload) // Marshal the payload to JSON, data is a byte slice containing the JSON representation of whatever was in the payload
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json") // Set the content type to application/json
	w.WriteHeader(code)                                // Set the HTTP status code
	w.Write(data)                                      // Write the JSON data to the response
}

func isProduction() bool {
	return os.Getenv("PLATFORM") == "prod"
}
