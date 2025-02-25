package handlers

import (
	"fmt"
	"database/sql"
	"net/http"
	"strings"
	"time"

	e "forum/Error"
	d "forum/database"
	m "forum/models"
	u "forum/utils"

	"github.com/gofrs/uuid"
)

// ===== User struct to store user details= ===
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
		e.LOGGER("[ERROR]", err)
		return
	}

	m.User.Username = strings.TrimSpace(r.FormValue("username"))
	m.User.Email = strings.TrimSpace(r.FormValue("email"))
	m.User.Password = strings.TrimSpace(r.FormValue("password"))

	if m.User.Username == "" || m.User.Email == "" || m.User.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Leading and trailing spaces are not accepted"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if !u.IsValidEmail(m.User.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid email address"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if u.CredentialExists(d.Db, m.User.Username) || u.CredentialExists(d.Db, m.User.Email) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Username or email already in use"))
		http.Redirect(w, r, "/register", http.StatusSeeOther)
		return
	}

	if err = m.User.HashPassword(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|register user handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	// ======= generate a UUID for the user using the google/uuid package =====
	// UUID := uuid.New().String()
	u, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|register user handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	UUID := u.String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert new user into the database
	query := `INSERT INTO users (uuid, username, email, password) VALUES (?, ?, ?, ?)`
	_, err = d.Db.Exec(query, UUID, m.User.Username, m.User.Email, m.User.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|register user handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	var userID string
	err = d.Db.QueryRow("SELECT id FROM users WHERE email = ?", m.User.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|register user handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|register user handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	//grant a session on registration
	_, err = d.Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, UUID, expiresAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|register user handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	SetSessionCookie(w, UUID, expiresAt)

	e.LOGGER(fmt.Sprintf("[SUCCESS]: %s was registered successfully!", m.User.Username), nil)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
