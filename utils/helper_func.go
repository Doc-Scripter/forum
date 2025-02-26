package utils

import (
	"database/sql"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sort"
	"time"

	e "forum/Error"
	d "forum/database"
	m "forum/models"
)

// ==============Validate email format==========
func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// ===== Check if a username or email (credential) already exists in a database db =====
func CredentialExists(db *sql.DB, credential string) bool {
	query := `SELECT COUNT(*) FROM users WHERE username = ? OR email = ?`
	var count int
	err := db.QueryRow(query, credential, credential).Scan(&count)
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|credential exist| ---> {%v}", err))
		return false
	}
	return count > 0
}

/*
=== ValidateSession checks if a session token is valid. The function takes a pointer to the request
and returns a boolean value and a user_ID of type string based on the session_token found in the
cookie present in the header, within the request =====
*/
func ValidateSession(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|validate session| ---> no session cookie found"))
		return false, ""
	}

	var (
		userID    string
		expiresAt time.Time
	)

	err = d.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|validate session| ---> {%v}", err))
		return false, ""
	}

	if time.Now().After(expiresAt) {
		e.LOGGER("[ERROR]", fmt.Errorf("session expired for user %s", userID))
		return false, ""
	}

	e.LOGGER(fmt.Sprintf("[SUCCESS]: Session valid for user: %s", userID), nil)
	return true, userID
}

// ==== The function will take  a response writer w and a request r, for displaying errors when encountered. The function returns a struct of a user model ====
func GetUserDetails(w http.ResponseWriter, r *http.Request) (m.ProfileData, error) {

	var (
		Profile m.ProfileData
		userID  string
	)

	cookie, err := r.Cookie("session_token")
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("profile Section: No session cookie found: %v", err))
		return m.ProfileData{}, err
	}

	err = d.Db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("session not found in DB: %v", err))
		return m.ProfileData{}, err
	}

	query := `
		SELECT  username, email , uuid  FROM users WHERE id = ?`

	err = d.Db.QueryRow(query, userID).Scan(&Profile.Username, &Profile.Email, &Profile.Uuid)
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("session not found in DB: %v", err))
		return m.ProfileData{}, err
	}

	Profile.Initials = Profile.GenerateInitials()

	e.LOGGER("[SUCCESS]: Retrieved User details successfully", nil)
	return Profile, nil
}

// =====The function to make all the categories as a string to be stored into the database===========
func CombineCategory(category []string) string {

	e.LOGGER("[SUCCESS]: Combined the categories as a string to be stored into the database", nil)
	return strings.Join(category, ", ")
}


//===== The function will be called to validate the values of the categories from the frontend ======
func ValidateCategory(str []string) bool {
	categories := []string{"All Categories", "Technology", "Health", "Math", "Nature", "Science", "Religion", "Education", "Politics", "Fashion", "Lifestyle", "Sports", "Arts"}

	for _, s := range str {
		for i, v := range categories {

			if s == v{
				return true
			}else {
				i++
			}
		}
	}
	return false
}

// ==== The function will sort the array of comments or posts by time before they are martialled into a json object =====
func OrderComments(comments []m.Comment) []m.Comment{
	sort.Slice(comments, func(i, j int) bool {
        return comments[i].CreatedAt.After(comments[j].CreatedAt)
    })
    return comments
}

func OrderPosts(posts []m.Post) []m.Post{
	sort.Slice(posts, func(i, j int) bool {
        return posts[i].CreatedAt.After(posts[j].CreatedAt)
    })
    return posts
}