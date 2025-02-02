package auth

import (
	"time"
	"net/http"
	"database/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUserCredentialsLogin(w http.ResponseWriter, r *http.Request) {
	

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		Cover(w, r)
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		Cover(w, r)
	}

	// Get username and password
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		Cover(w, r)
	}

	// Retrieve user from DB using their email
	var dbPassword, userID string
	err = Db.QueryRow("SELECT password, uuid FROM users WHERE email = ?", email).Scan(&dbPassword, &userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		Cover(w, r)
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		Cover(w, r)
	}

	// Compare passwords
	if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		Cover(w, r)
	}

	// Generate a session token
	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Store session in the database
	_, err = Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, sessionToken, expiresAt)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// Set the cookie with the session token
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   false,
	})

	w.WriteHeader(http.StatusCreated)
	HomePage(w, r)
}
