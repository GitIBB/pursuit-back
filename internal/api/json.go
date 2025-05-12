package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) { // Function to respond with an error message in JSON format
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct { // Struct to hold the error message
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{ // Create a JSON response with the error message
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) { // Function to respond with a JSON payload
	// Set the content type to application/json
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload) // Marshal the payload to JSON, data is a byte slice containing the JSON representation of whatever was in the payload
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code) // Set the HTTP status code
	w.Write(data)       // Write the JSON data to the response
}
