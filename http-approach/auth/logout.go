package auth

import(
	"time"
	"net/http"
)

// LogoutUser removes session from the database and clears the cookie
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	// Remove session from the database
	_, err = Db.Exec("DELETE FROM sessions WHERE session_token = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}