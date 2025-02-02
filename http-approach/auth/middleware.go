package auth

import(
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		valid, userID := ValidateSession(r)
		if !valid {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		r.Header.Set("User-ID", userID)

		next(w, r)
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