package handlers

import (
	"fmt"
	"log"
	"time"
	"strings"
	"regexp"
	"net/http"
	"database/sql"
	
	
	m "forum/models"
	e "forum/Error"
	d "forum/database"
	"golang.org/x/crypto/bcrypt"
)

// ==============Validate email format==========
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ===== Check if a username or email (credential) already exists in a database db =====
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

/*=== ValidateSession checks if a session token is valid. The function takes a pointer to the request
and returns a boolean value and a user_ID of type string based on the session_token found in the 
cookie present in the header, within the request =====*/
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


// ==== The function will take  a response writer w and a request r, for displaying errors when encountered. The function returns a struct of a user model ====
func GetUserDetails(w http.ResponseWriter, r *http.Request) m.ProfileData {
	
	var (
		Profile m.ProfileData
		userID string
	)
	
	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Profile Section: No session cookie found:", err)
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return m.ProfileData{}
	}

	err = d.Db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println("Session not found in DB:", err)
		e.LogError(err)
		return m.ProfileData{}
	}

	query := `
		SELECT  username, email , uuid  FROM users WHERE id = ?`

	err = d.Db.QueryRow(query, userID).Scan(&Profile.Username, &Profile.Email, &Profile.Uuid)
	if err != nil {
		e.LogError(err)
		return m.ProfileData{}
	}
	
	Profile.Initials = Profile.GenerateInitials()
	
	return Profile
}

//=====The function to make all the categories as a string to be stored into the database===========
func CombineCategory(category []string) string{

	return strings.Join(category, ", ")
}
