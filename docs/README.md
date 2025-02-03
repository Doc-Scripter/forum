# Forum Authentication

- This is the complete Documentation for the Authentication of the Forum Application project found in [my Github account](https://github.com/anxielray/Forum-Application.git)
- The Authentication used here compassed the following:

  ## **1. Core Authentication Features**

  The authentication system includes:


  1. **User Registration** (with validation and secure password storage)
     However this is considered a weak mode of security.
  2. **User Login** (with session handling using cookies)
  3. **Session Management** (using cookies with expiration)
     A time lapse of, when user exits from the browser, is valid, when set to 0
  4. **Logout** (clearing session cookies)
  5. **Authorization Checks** (restricting actions for logged-in users)
  6. **Password Encryption** (hashing using bcrypt)
  7. **Unique User Identification** (UUID for better tracking)

## **2. Database Schema for Authentication**

We need a **users** table to store user credentials. Here’s a structured way to define it in SQLite:

<pre class="!overflow-visible"><div class="contain-inline-size rounded-md border-[0.5px] border-token-border-medium relative bg-token-sidebar-surface-primary dark:bg-gray-950"></div></pre>

```sql
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);
```

uuid: A unique user ID (generated via uuid.New().String()) the package is standard in Golang.
username: Unique identifier for users
email: Must be unique, used for login
password: Must be hashed using bcrypt before storing

## **3. Implementing User Registration**

Here’s a function to **register a new user** in Go:

```go
package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User struct
type User struct {
	Username string
	Email    string
	Password string
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

// RegisterUser handles user registration
func RegisterUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	user := User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	// Validate input
	if user.Username == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Hash the password
	if err := user.HashPassword(); err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Generate UUID
	userID := uuid.New().String()

	// Insert into database
	_, err := db.Exec("INSERT INTO users (uuid, username, email, password) VALUES (?, ?, ?, ?)",
		userID, user.Username, user.Email, user.Password)
	if err != nil {
		http.Error(w, "User registration failed", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User registered successfully!")
}
```

## 4. Implementing User Login

Login involves:

- Checking if the user exists
- Verifying the password with bcrypt
- Creating a session using cookies

```go
func LoginUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Validate input
	if email == "" || password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Retrieve user from DB
	var dbPassword, userID string
	err := db.QueryRow("SELECT password, uuid FROM users WHERE email = ?", email).Scan(&dbPassword, &userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Create session cookie
	sessionToken := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   3600, // 1 hour session
	})

	fmt.Fprintln(w, "Login successful!")
}
```

## 5. Implementing Logout

- Logging out involves clearing the session cookie:

```go
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1, // Expire immediately
	})

	fmt.Fprintln(w, "Logged out successfully!")
}
```

## 6. Authentication Middleware

- Middleware ensures only logged-in users can access certain pages:

```go
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the session is valid (this would require session management in DB)
		if cookie.Value == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
```

## 7. Next Steps & Additional Features

🔹 Session Storage

- Right now, session tokens are just cookies, but ideally, you should store session tokens in the database and validate them.

```sql
CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    session_token TEXT NOT NULL,
    expiry TIMESTAMP NOT NULL
);
```

- Each time a user logs in, their session should be saved, and during authentication, check if their session token is valid.
  🔹 Reset Password
- You can implement password reset functionality by:

  Generating a reset token (UUID)
  Sending it to the user's email
  Creating a form for resetting passwords

🔹 Account Verification

- To prevent fake accounts, send a verification email with a unique activation link.
  🔹 Role-Based Access Control (RBAC)
- You can define user roles (admin, moderator, user) and restrict access based on roles.
