package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/GitIBB/pursuit/internal/api"
	"github.com/GitIBB/pursuit/internal/database"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const filepathRoot = "." // sets the root directory for the file server
	const port = "8080"      // sets port for the server to listen on

	godotenv.Load()              // Load environment variables from .env file
	dbURL := os.Getenv("DB_URL") // Get database URL from environment variable
	if dbURL == "" {             // Check if DB_URL is set
		log.Fatal("DB_URL not set") // log error and exit program if DB_URL is not set
	}

	platform := os.Getenv("PLATFORM") // Get platform name from environment variable
	if platform == "" {
		log.Fatal("PLATFORM NOT SET")
	}
	jwtSecret := os.Getenv("JWT_SECRET") // Get JWT secret from environment variable
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET NOT SET")
	}

	dbCon, err := sql.Open("postgres", dbURL)
	if err != nil { // Check if there was an error opening the database connection
		log.Fatal("Error opening database connection:", err) // log error and exit program if there was an error
	}
	dbQueries := database.New(dbCon) // Create a new database connection using the provided URL

	apiCfg := api.NewAPIConfig(dbQueries, platform, jwtSecret)

	// Create a new HTTP server mux (router)
	mux := http.NewServeMux()
	apiCfg.SetupRoutes(mux, "./static") // Setup routes for the API using the provided configuration

	srv := http.Server{ // Server config
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on port: %s\n", port) // log server startup and which port it is listening on
	log.Fatal(srv.ListenAndServe())           // will log error and exit program if server fails
}
