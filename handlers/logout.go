package handlers

import (
	"time"
	"net/http"
	d "forum/database"
)

// LogoutUser removes session from the database and clears the cookie

func LogoutUser(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther) //though the route should be the home directory
		return
	}

	_, err = d.Db.Exec("DELETE FROM sessions WHERE session_token = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Failed to log out", http.StatusInternalServerError)
		return
	}

	SetSessionCookie(w, "", time.Now().Add(-time.Hour))

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
