package api

import (
	"log"
	"net/http"
	"time"
)

func (cfg *APIConfig) middlewareLogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Capture response status
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		// Log slow requests
		if duration > time.Duration(5)*time.Second {
			log.Printf("[SLOW REQUEST] %s %s from %s took %v", r.Method, r.URL.Path, r.RemoteAddr, duration)
		}

		// Log error responses (4xx and 5xx)
		if rec.statusCode >= 400 {
			log.Printf("[ERROR RESPONSE] %s %s from %s - Status: %d", r.Method, r.URL.Path, r.RemoteAddr, rec.statusCode)
		}

		// Log unauthorized access attempts to critical endpoints
		if rec.statusCode == http.StatusUnauthorized && isCriticalEndpoint(r.URL.Path) {
			log.Printf("[UNAUTHORIZED ACCESS] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		}

		// Log high-frequency requests or unusual activity (example: specific endpoint)
		if isHighFrequencyEndpoint(r.URL.Path) {
			log.Printf("[HIGH FREQUENCY REQUEST] %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		}
	})
}

// Helper to determine if an endpoint is critical
func isCriticalEndpoint(path string) bool {
	criticalEndpoints := []string{"/admin/metrics", "/admin/reset", "/login"}
	for _, endpoint := range criticalEndpoints {
		if path == endpoint {
			return true
		}
	}
	return false
}

// Helper to determine if an endpoint is high-frequency
func isHighFrequencyEndpoint(path string) bool {
	highFrequencyEndpoints := []string{"/app/", "/api/healthz"}
	for _, endpoint := range highFrequencyEndpoints {
		if path == endpoint {
			return true
		}
	}
	return false
}

// Helper to capture response status
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
