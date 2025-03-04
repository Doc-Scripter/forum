package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	e "forum/Error"
	d "forum/database"
	m "forum/models"
	u "forum/utils"
	"github.com/gofrs/uuid" // the google package is also allowed
	"golang.org/x/crypto/bcrypt"
)

// ==============This function will be called when a the login submission is done=====================
func AuthenticateUserCredentialsLogin(w http.ResponseWriter, r *http.Request) {

	if bl, _ := u.ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to process login request"))
		e.LOGGER("[ERROR]", fmt.Errorf("|login handler| ---> {%v}", err))
		// ErrorPage(err, m.ErrorsData.InternalError , w, r)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := strings.TrimSpace(r.FormValue("password"))

	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email and password are required"))
		return
	}

	var dbPassword, userID string
	err = d.Db.QueryRow("SELECT password, id FROM users WHERE email = ?", email).Scan(&dbPassword, &userID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid email or password"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred while processing your request"))
		e.LOGGER("[ERROR]", fmt.Errorf("|login handler| ---> {%v}", err))
		return
	}

	//====================Compare hashed passwords=======================
	if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid email or password"))
		return
	}

	//==========================Create a session for a user on logging in=================
	u, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred while logging you in"))
		e.LOGGER("[ERROR]", fmt.Errorf("|login handler| ---> {%v}", err))
		return
	}

	sessionToken := u.String()
	expiresAt := time.Now().Add(24 * time.Hour)

	_, err = d.Db.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|login handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
	SetSessionCookie(w, "", time.Now().Add(-time.Hour))

	//==============Store the created session=================
	_, err = d.Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)",
		userID, sessionToken, expiresAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred while logging  you in"))
		e.LOGGER("[ERROR]", fmt.Errorf("|login handler| ---> {%v}", err))
		return
	}

	SetSessionCookie(w, sessionToken, expiresAt)

	e.LOGGER("[SUCCESS]: Logged in user successfully!", nil)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Redirecting you to home page"))
}
