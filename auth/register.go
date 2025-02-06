package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// User struct to store user details
type User struct {
	ID       int
	UUID     string
	Username string
	Email    string
	Password string
}

// =========Handle user registration========================
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/registration" {
		fmt.Fprint(w, "Error: Bad request!")
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get user input
	user := User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Validate input fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !IsValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Check if username or email already exists
	if credentialExists(Db, user.Username) || credentialExists(Db, user.Email) {
		http.Error(w, "Username or email already in use", http.StatusConflict)
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	if err = user.HashPassword(); err != nil {
		fmt.Fprint(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// generate a UUID for the user
	UUID := uuid.New().String()
	expiresAt := time.Now().Add(24 * time.Hour)

	// Insert new user into the database
	query := `INSERT INTO users (uuid, username, email, password) VALUES (?, ?, ?, ?)`
	_, err = Db.Exec(query, UUID, user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	var userID string
	err = Db.QueryRow("SELECT id FROM users WHERE email = ?", user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		fmt.Println(err)
	} else if err != nil {

		http.Error(w, "Database error", http.StatusInternalServerError)
		fmt.Println(err)
	}

	// Store session in the database
	// _, err = Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, UUID, expiresAt)
	_, err = Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, UUID, expiresAt)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	// Success response
	SetSessionCookie(w, UUID, expiresAt)

	w.WriteHeader(http.StatusCreated)
	HomePage(w, r)
}
