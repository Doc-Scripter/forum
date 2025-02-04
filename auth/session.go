package auth

import(
	"time"
	"net/http"
)



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