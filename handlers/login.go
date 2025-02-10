package handlers

import (
	"fmt"
	"time"
	"net/http"
	"database/sql"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUserCredentialsLogin(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println(r.Method)
	if bl, _ := ValidateSession(r); bl {
		HomePage(w, r)
		return
	}else if !bl{
		if r.Method != http.MethodPost {			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			// http.Redirect(w, r, "/web//", http.StatusFound)
			LandingPage(w, r)
			return
		}
	
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		HomePage(w, r)

	}

	// Get username and password
	email := r.FormValue("email")
	password := r.FormValue("password")
	
	// Validate input
	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		HomePage(w, r)
	}
	
	// Retrieve user from DB using their email
	var dbPassword, userID string
	err = Db.QueryRow("SELECT password, id FROM users WHERE email = ?", email).Scan(&dbPassword, &userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		HomePage(w, r)
		} else if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			HomePage(w, r)
		}
		
		// Compare passwords
		if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			HomePage(w, r)
		}
		
		// Generate a session token
		// sessionToken := uuid.New().String()

		u, err := uuid.NewV4()
		if err != nil {
			fmt.Println("Error generating UUID:", err)
			return
		}
		sessionToken := u.String()
		expiresAt := time.Now().Add(24 * time.Hour)
		
		// Store session in the database
		_, err = Db.Exec("INSERT INTO sessions (user_id, session_token, expires_at) VALUES (?, ?, ?)", userID, sessionToken, expiresAt)
		if err != nil {
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			RegisterUser(w, r)
		}
		
		SetSessionCookie(w, sessionToken, expiresAt)
		
		w.WriteHeader(http.StatusCreated)
		HomePage(w, r)
	}
}