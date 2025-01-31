package auth

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"regexp"
	"database/sql"
	"net/http"
)

// User struct to store user details
type User struct {
	UUID     string
	Username string
	Email    string
	Password string
}

// Validate email format
func isValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// Check if a username or email already exists in the database
func credentialExists(db *sql.DB, credential string) bool {
	query := `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`
	var count int
	err := db.QueryRow(query, credential, credential).Scan(&count)
	if err != nil {
		log.Printf("Database error: %s", err)
		return false
	}
	return count > 0
}

// HashPassword hashes the user's password before storing it
func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// =========Handle user registration========================
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// Get user input
	user := User{
		UUID:     uuid.New().String(),
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Validate input fields
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	if !isValidEmail(user.Email) {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		return
	}

	// Check if username or email already exists
	if credentialExists(Db, user.Username) || credentialExists(Db, user.Email) {
		http.Error(w, "Username or email already in use", http.StatusConflict)
		http.Redirect(w, r, "/signup", http.StatusFound)
		return
	}

	if err =  user.HashPassword(); err != nil {
		fmt.Fprint(w, "Failed to hash password", http.StatusInternalServerError)
        return
	}

	// Insert new user into the database
	query := `INSERT INTO users (uuid, username, email, password) VALUES (?, ?, ?, ?)`
	_, err = Db.Exec(query, user.UUID, user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Success response
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "User registered successfully")

}