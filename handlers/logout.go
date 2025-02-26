package handlers

import (
	e "forum/Error"
	d "forum/database"
	m "forum/models"
	"fmt"
	"net/http"
	"time"
)

//====== LogoutUser removes session from the database and clears the cookie =========

func LogoutUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	if r.URL.Path != "/logout" {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|logout handler| ---> user tried to access the logout page in a wrong url"), m.ErrorsData.PageNotFound, w, r)
        return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		e.LOGGER("[ERROR]", fmt.Errorf("|logout handler| ---> {%v}", err))
		return
	}

	_, err = d.Db.Exec("DELETE FROM sessions WHERE session_token = ?", cookie.Value)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|logout handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	SetSessionCookie(w, "", time.Now().Add(-time.Hour))

	defer d.Db.Close()
	e.LOGGER("[SUCCESS]: Closed the database successfully to prevent any action after logging out!", nil)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ==================This function will protected the private endpoints===============
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// headers to prevent caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

	// Your protected content logic here
	w.Write([]byte("This is protected content.")) //I will have a deeper look at this later
}
