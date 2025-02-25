package handlers

import (
	"net/http"
	"time"
)

// ==== This function will set up a cookie using the UUID(sessionToken) and sets the expiration time to expiresAt and sets the access of the cookie to html only ====
func SetSessionCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
	})
}
