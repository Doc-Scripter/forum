package auth

import (
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate session and retrieve user ID
		valid, userID := ValidateSession(r)

		if !valid {
			// If session is invalid, redirect to login page
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// If valid, add the User-ID to the request header for downstream handlers
		r.Header.Set("User-ID", userID)

		// Continue with the next handler in the middleware chain
		next.ServeHTTP(w, r)
	}
}

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// Check for session cookie
// 		cookie, err := r.Cookie("session_token")
// 		if err != nil {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		// Check if the session is valid (this would require session management in DB)
// 		if cookie.Value == "" {
// 			http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }
