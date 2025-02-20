package handlers

import (
	"time"
	"strings"
	"net/http"
	"database/sql"

	e "forum/Error"
	m "forum/models"
	d "forum/database"

	"github.com/gofrs/uuid"
)

//=======User struct to store user details===
type User struct {
	ID       int
	UUID     string
	Username string
	Email    string
	Password string
}

// ==== The handler will perform registration on submission of a form from the frontend. Register a new user and generate them a cookie ====
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to process registration request"))
		e.LogError(err)
		return
	}

	user := User{
		Username: strings.TrimSpace(r.FormValue("username")),
		Email:    strings.TrimSpace(r.FormValue("email")),
		Password: strings.TrimSpace(r.FormValue("password")),
	}

	if user.Username == "" || user.Email == "" || user.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Leading and trailing spaces are not accepted"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if !IsValidEmail(user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid email address"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if credentialExists(d.Db, user.Username) || credentialExists(d.Db, user.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username or email already in use"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if err = user.HashPassword(); err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}

	// generate a UUID for the user using the google/uuid package
	// UUID := uuid.New().String()
	u, err := uuid.NewV4()
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}

	UUID := u.String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert new user into the database
	query := `INSERT INTO users (uuid, username, email, password) VALUES (?, ?, ?, ?)`
	_, err = d.Db.Exec(query, UUID, user.Username, user.Email, user.Password)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}

	var userID string
	err = d.Db.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
	} else if err != nil {

		ErrorPage(err, m.ErrorsData.InternalError, w, r)
	}

	//grant a session on registration
	_, err = d.Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, UUID, expiresAt)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}

	SetSessionCookie(w, UUID, expiresAt)

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
