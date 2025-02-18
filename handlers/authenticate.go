package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"

	d "forum/database"

	"golang.org/x/crypto/bcrypt"
)

// ==============Validate email format==========
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// =============Check if a username or email already exists in the database==========
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

// ==============HashPassword hashes the user's password before storing it=========
func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// ==================ValidateSession checks if a session token is valid=========
func ValidateSession(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("No session cookie found")
		return false, ""
	}

	var (
		userID    string
		expiresAt time.Time
	)

	err = d.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		fmt.Println("Session not found in DB:", err)
		return false, ""
	}

	if time.Now().After(expiresAt) {
		fmt.Println("Session expired")
		return false, ""
	}

	fmt.Println("Session valid for user:", userID)
	return true, userID
}

func validateFileType(file multipart.File) bool {
	// Read the first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// Reset the file pointer
	file.Seek(0, 0)

	// Get content type and check if it's allowed
	contentType := http.DetectContentType(buffer)
	return strings.Contains(allowedTypes, contentType)
}

func generateFileName() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
