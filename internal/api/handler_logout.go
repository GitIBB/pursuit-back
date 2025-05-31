package api

import (
	"net/http"
	"time"
)

func (cfg *APIConfig) handlerLogout(w http.ResponseWriter, r *http.Request) {
	// clear auth-token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth-token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // expire immediately
		HttpOnly: true,
		Secure:   true, // Set to true in production with HTTPS
		SameSite: http.SameSiteNoneMode,
	})

	// respond with success message
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Logged out successfully"}`))
}
