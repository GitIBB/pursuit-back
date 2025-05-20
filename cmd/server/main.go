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

// corsMiddleware is a middleware function that adds CORS (Cross-Origin Resource Sharing)
// headers to the response.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Adjust origin as needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	const port = "8080" // sets port for the server to listen on

	godotenv.Load("/home/hermit/workspace/pursuit/pursuit-back/.env") // Load environment variables from .env file
	dbURL := os.Getenv("DB_URL")                                      // Get database URL from environment variable
	if dbURL == "" {
		log.Fatal("DB_URL not set")
	}

	platform := os.Getenv("PLATFORM") // Get platform name from environment variable
	if platform == "" {
		log.Fatal("PLATFORM NOT SET")
	}

	jwtSecret := os.Getenv("JWT_SECRET") // Get JWT secret from environment variable
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET NOT SET")
	}

	dbCon, err := sql.Open("postgres", dbURL) // Open a connection to the PostgreSQL database using the provided URL
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}
	dbQueries := database.New(dbCon) // Create a new database connection using the provided URL

	apiCfg := api.NewAPIConfig(dbQueries, platform, jwtSecret) // Create a new instance of the apiConfig struct

	mux := http.NewServeMux()      // Create a new HTTP server mux (router)
	apiCfg.SetupRoutes(mux)        // Setup routes for the API using the provided configuration
	handler := corsMiddleware(mux) // Apply CORS middleware to the mux

	srv := http.Server{ // Server config
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Serving on port: %s\n", port) // log server startup and which port it is listening on
	log.Fatal(srv.ListenAndServe())           // will log error and exit program if server fails
}
