package api

import (
	"context"
	"net/http"
	"strings"
)

// Context key type to avoid conflicts
type contextKey string

const isBrowserKey contextKey = "is_browser"

// BrowserAwareness middleware detects and tracks browser clients
func (cfg *APIConfig) middlewareBrowserAwareness(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Detect if the client is a browser based on the User-Agent header
		userAgent := r.Header.Get("User-Agent")
		isBrowser := strings.Contains(userAgent, "Mozilla") || strings.Contains(userAgent, "Chrome") || strings.Contains(userAgent, "Safari")

		// Store the result in the request context
		ctx := context.WithValue(r.Context(), isBrowserKey, isBrowser)

		// Call the next handler with the enriched context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to retrieve the browser status from the context
func IsBrowser(r *http.Request) bool {
	isBrowser, ok := r.Context().Value(isBrowserKey).(bool)
	return ok && isBrowser
}
