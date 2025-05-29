package api

import (
	"net/http"
	"time"
)

func (cfg *APIConfig) handlerLogout(w http.ResponseWriter, r *http.Request) {
	// Clear the auth-token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Expire the cookie immediately
		HttpOnly: true,
		Secure:   false,                 // Set to true in production with HTTPS
		SameSite: http.SameSiteNoneMode, // Allow cross-site usage
	})

	// Respond with a success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}
