package handlers

import (
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		valid, userID := ValidateSession(r)
		
		if !valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		r.Header.Set("User-ID", userID)
		next.ServeHTTP(w, r)
	}
}
