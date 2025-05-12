package api

import (
	"fmt"
	"net/http"
)

func (cfg *APIConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) { // Handler function for the /admin/metrics endpoint
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`
	<html>

<body>
	<h1>Welcome, Pursuit admin</h1>
	<p>Pursuit has been visited %d times!</p>
</body>

</html>
	`, cfg.fileserverHits.Load()))) // Respond with the number of hits to the file server
}

func (cfg *APIConfig) middlewareMetricsInc(next http.Handler) http.Handler { // Middleware function to increment the file server hits counter
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
