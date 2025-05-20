package api

import "net/http"

func (cfg *APIConfig) SetupRoutes(mux *http.ServeMux) {
	// Health Check
	mux.HandleFunc("/api/healthz", handlerReadiness) // Register readiness endpoint at /healthz path, delegates handling to the handlerReadiness function

	// Auth endpoints
	mux.Handle("POST /api/login", cfg.middlewareBrowserAwareness(http.HandlerFunc(cfg.handlerLogin))) // Register login endpoint at /login path, delegates handling to the handlerLogin function
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)                                           // Register refresh endpoint at /refresh path, delegates handling to the handlerRefresh function
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)                                             // Register revoke endpoint at /revoke path, delegates handling to the handlerRevoke function

	// User endpoints
	mux.HandleFunc("POST /api/users", cfg.handlerUsersCreate)                                  // Register user creation endpoint at /users path, delegates handling to the handlerUsersCreate function
	mux.Handle("PUT /api/users", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUsersUpdate))) // Register user update endpoint at /users path, delegates handling to the handlerUsersUpdate function
	mux.Handle("GET /api/users", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerUsersGet)))    // Register user retrieval endpoint at /users path, delegates handling to the handlerUsersGet function

	// Article endpoints
	mux.Handle("POST /api/articles", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerArticlesCreate)))               // Register article creation endpoint at /articles path, delegates handling to the handlerArticlesCreate function
	mux.HandleFunc("GET /api/articles", cfg.handlerArticlesRetrieve)                                                // Register article (all) retrieval endpoint at /articles path, delegates handling to the handlerArticlesRetrieve function
	mux.HandleFunc("GET /api/articles/{articleID}", cfg.handlerArticlesGet)                                         // Register article retrieval endpoint at /articles/{articleID} path, delegates handling to the handlerArticlesGet function
	mux.Handle("DELETE /api/articles/{articleID}", cfg.middlewareAuth(http.HandlerFunc(cfg.handlerArticlesDelete))) // Register article deletion endpoint at /articles/{articleID} path, delegates handling to the handlerArticlesDelete function

	// Admin endpoints
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics) // Register metrics endpoint at /metrics path, delegates handling to the handlerMetrics function
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)    // Register readiness endpoint at /healthz path, delegates handling to the handlerReadiness function

}
