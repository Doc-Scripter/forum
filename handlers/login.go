package handlers

import (
	"time"
	"net/http"
	"database/sql"
	// "fmt"

	d "forum/database"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUserCredentialsLogin(w http.ResponseWriter, r *http.Request) {
	if bl, _ := ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid request method"))
		return
	}

	// Parse both URL-encoded form data and multipart form data
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to process login request"))
		return
	}
	
	// Log the raw form data for debugging
	// fmt.Printf("Form data received: %+v\n", r.Form)
	
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Log the extracted values
	// fmt.Printf("Extracted email: %s, password: %s\n", email, password)

	// Validate input
	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Email and password are required"))
		return
	}

	// Retrieve user from DB
	var dbPassword, userID string
	err = d.Db.QueryRow("SELECT password, id FROM users WHERE email = ?", email).Scan(&dbPassword, &userID)
	if err == sql.ErrNoRows {
		// User not found - return error but don't specify which credential was wrong
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid email or password"))
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred while processing your request"))
		return
	}

	// Compare passwords
	if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid email or password"))
		return
	}

	// Create session
	u, err := uuid.NewV4()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred while creating your session"))
		return
	}

	sessionToken := u.String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Store session
	_, err = d.Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", 
		userID, sessionToken, expiresAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("An error occurred while creating your session"))
		return
	}

	SetSessionCookie(w, sessionToken, expiresAt)
	
	// Return success response
	w.WriteHeader(http.StatusOK)
}
