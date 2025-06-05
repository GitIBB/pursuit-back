package api

import (
	"log"
	"net/http"
	"path/filepath"
)

func (cfg *APIConfig) SetupRoutes(mux *http.ServeMux) {
	// Health Check
	mux.HandleFunc("/api/healthz", handlerReadiness) // Register readiness endpoint at /healthz path, delegates handling to the handlerReadiness function

	// Auth endpoints
	mux.Handle("POST /api/login", cfg.middlewareBrowserAwareness(http.HandlerFunc(cfg.handlerLogin))) // Register login endpoint at /login path, delegates handling to the handlerLogin function
	mux.Handle("POST /api/logout", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerLogout)))           // Register logout endpoint at /logout path, delegates handling to the handlerLogout function
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)                                           // Register refresh endpoint at /refresh path, delegates handling to the handlerRefresh function
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)                                             // Register revoke endpoint at /revoke path, delegates handling to the handlerRevoke function

	// User endpoints
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)                                  // Register user creation endpoint at /users path, delegates handling to the handlerUsersCreate function
	mux.Handle("PUT /api/users", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUsersUpdate))) // Register user update endpoint at /users path, delegates handling to the handlerUsersUpdate function
	mux.Handle("GET /api/users", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUsersGet)))    // Register user retrieval endpoint at /users path, delegates handling to the handlerUsersGet function
	mux.Handle("GET /api/me", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerMe)))             // Register user (me) retrieval endpoint at /me path, delegates handling to the handlerMe function

	// Uploads endpoints
	// TEMPORARY USAGE: REPLACE STORAGE IN UPLOADS FOLDER WITH CDN OR BUCKET STORAGE IN PRODUCTION
	uploadsDir, err := filepath.Abs(filepath.Join("..", "..", "uploads"))
	if err != nil {
		log.Fatalf("Failed to get absolute path for uploads directory: %v", err)
	}
	mux.Handle("/api/uploads/", http.StripPrefix("/api/uploads/", http.FileServer(http.Dir(uploadsDir))))
	mux.Handle("POST /api/uploads", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUploads)))

	// Article endpoints
	mux.Handle("POST /api/articles", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerArticlesCreate)))               // Register article creation endpoint at /articles path, delegates handling to the handlerArticlesCreate function
	mux.HandleFunc("GET /api/articles", cfg.handlerArticlesRetrieve)                                                // Register article (all) retrieval endpoint at /articles path, delegates handling to the handlerArticlesRetrieve function
	mux.HandleFunc("GET /api/articles/{articleID}", cfg.handlerArticlesGet)                                         // Register article retrieval endpoint at /articles/{articleID} path, delegates handling to the handlerArticlesGet function
	mux.Handle("DELETE /api/articles/{articleID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerArticlesDelete))) // Register article deletion endpoint at /articles/{articleID} path, delegates handling to the handlerArticlesDelete function
	mux.HandleFunc("GET /api/users/{userID}/articles", cfg.handlerUserArticles)                                     // Register user articles retrieval endpoint at /users/{userID}/articles path, delegates handling to the handlerUserArticles function

	// Category endpoint
	mux.HandleFunc("GET /api/categories", cfg.handlerCategoriesGet) // Register categories retrieval endpoint at /categories path, delegates handling to the handlerCategoriesGet function

	// Admin endpoints
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics) // Register metrics endpoint at /metrics path, delegates handling to the handlerMetrics function
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)    // Register readiness endpoint at /healthz path, delegates handling to the handlerReadiness function

}
