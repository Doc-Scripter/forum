package handlers

import (
	"net/http"
	u "forum/utils"
)

/*===== AuthMiddleware is a middleware function that checks if a user is authenticated before allowing them to access certain routes.====*/
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		valid, userID := u.ValidateSession(r)

		if !valid {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.Header.Set("User-ID", userID)
		next.ServeHTTP(w, r)
	}
}
