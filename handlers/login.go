package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func AuthenticateUserCredentialsLogin(w http.ResponseWriter, r *http.Request) {
	
	if bl, _ := ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}else if !bl{

		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	
	    err := r.ParseForm()
	    if err != nil {
	    	http.Error(w, "Unable to parse form", http.StatusBadRequest)
	    	http.Redirect(w, r, "/", http.StatusSeeOther)
    
	    }
    
	    // Get username and password
	    email := r.FormValue("email")
	    password := r.FormValue("password")
	    
	    // Validate input
	    if email == "" || password == "" {
	    	http.Error(w, "Email and password are required", http.StatusBadRequest)
	    	http.Redirect(w, r, "/", http.StatusSeeOther)
	    }
	    
	    // Retrieve user from DB using their email
	    var dbPassword, userID string
	    err = Db.QueryRow("SELECT password, id FROM users WHERE email = ?", email).Scan(&dbPassword, &userID)
	    if err == sql.ErrNoRows {
	    	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	    	http.Redirect(w, r, "/login", http.StatusSeeOther)
	    } else if err != nil {
	    	http.Error(w, "Database error", http.StatusInternalServerError)
	    	http.Redirect(w, r, "/login", http.StatusSeeOther)
	    }
	    	
	    // Compare passwords
	    if err = bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
	    	http.Error(w, "Invalid email or password", http.StatusUnauthorized)
	    	http.Redirect(w, r, "/login", http.StatusSeeOther)
	    }
    
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
	    	return
	    }
	    	
	    SetSessionCookie(w, sessionToken, expiresAt)
	    	
	    http.Redirect(w, r, "/home", http.StatusSeeOther)
	    //w.WriteHeader(http.StatusCreated)
	}
}
