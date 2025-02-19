package handlers
import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)
// TestCredentialExists tests the credentialExists function
func TestCredentialExists(t *testing.T) {
	// Create a new mock database
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	// Test cases
	tests := []struct {
		name        string
		credential  string
		mockQuery   string
		mockArgs    []driver.Value // Use driver.Value instead of interface{}
		mockCount   int
		expectError bool
		expected    bool
	}{
		{
			name:        "Credential exists",
			credential:  "testuser",
			mockQuery:   `SELECT COUNT\(\*\) FROM users WHERE username = \? OR email = \?`,
			mockArgs:    []driver.Value{"testuser", "testuser"}, // Use driver.Value
			mockCount:   1,
			expectError: false,
			expected:    true,
		},
		{
			name:        "Credential does not exist",
			credential:  "nonexistent",
			mockQuery:   `SELECT COUNT\(\*\) FROM users WHERE username = \? OR email = \?`,
			mockArgs:    []driver.Value{"nonexistent", "nonexistent"}, // Use driver.Value
			mockCount:   0,
			expectError: false,
			expected:    false,
		},
		{
			name:        "Database error",
			credential:  "testuser",
			mockQuery:   `SELECT COUNT\(\*\) FROM users WHERE username = \? OR email = \?`,
			mockArgs:    []driver.Value{"testuser", "testuser"}, // Use driver.Value
			mockCount:   0,
			expectError: true,
			expected:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the expected database query and result
			rows := sqlmock.NewRows([]string{"count"}).AddRow(tt.mockCount)
			if tt.expectError {
				mock.ExpectQuery(tt.mockQuery).WithArgs(tt.mockArgs...).WillReturnError(sql.ErrConnDone)
			} else {
				mock.ExpectQuery(tt.mockQuery).WithArgs(tt.mockArgs...).WillReturnRows(rows)
			}
			// Call the function
			result := credentialExists(db, tt.credential)
			// Assert the result
			assert.Equal(t, tt.expected, result)
			// Ensure all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
// TestIsValidEmail tests the IsValidEmail function
func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name  string
		email string
		want  bool
	}{
		{
			name:  "valid email",
			email: "test@example.com",
			want:  true,
		},
		{
			name:  "valid email with subdomain",
			email: "test@sub.example.com",
			want:  true,
		},
		{
			name:  "valid email with plus sign",
			email: "test+alias@example.com",
			want:  true,
		},
		{
			name:  "valid email with dot in local part",
			email: "test.alias@example.com",
			want:  true,
		},
		{
			name:  "valid email with underscore",
			email: "test_alias@example.com",
			want:  true,
		},
		{
			name:  "invalid email missing @",
			email: "testexample.com",
			want:  false,
		},
		{
			name:  "invalid email missing domain",
			email: "test@",
			want:  false,
		},
		{
			name:  "invalid email missing local part",
			email: "@example.com",
			want:  false,
		},
		{
			name:  "invalid email with space",
			email: "test @example.com",
			want:  false,
		},
		{
			name:  "invalid email with invalid characters",
			email: "test!@example.com",
			want:  false,
		},
		{
			name:  "invalid email with too short domain",
			email: "test@example.c",
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValidEmail(tt.email)
			if got != tt.want {
				t.Errorf("IsValidEmail(%q) = %v, want %v", tt.email, got, tt.want)
			}
		})
	}
}
// TestUser_HashPassword tests the HashPassword method of the User struct
func TestUser_HashPassword(t *testing.T) {
	type fields struct {
		ID       int
		UUID     string
		Username string
		Email    string
		Password string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Successful hashing",
			fields: fields{
				ID:       1,
				UUID:     "some-uuid",
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantErr: false,
		},
		{
			name: "Empty password",
			fields: fields{
				ID:       2,
				UUID:     "another-uuid",
				Username: "anotheruser",
				Email:    "another@example.com",
				Password: "",
			},
			wantErr: false, // bcrypt handles empty strings without error
		},
		{
			name: "Very long password",
			fields: fields{
				ID:       3,
				UUID:     "long-uuid",
				Username: "longuser",
				Email:    "long@example.com",
				// Reduced the length to be within bcrypt's limit (72 bytes)
				Password: strings.Repeat("a", 72),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:       tt.fields.ID,
				UUID:     tt.fields.UUID,
				Username: tt.fields.Username,
				Email:    tt.fields.Email,
				Password: tt.fields.Password,
			}
			originalPassword := u.Password // Store the original password for later comparison
			err := u.HashPassword()
			if (err != nil) != tt.wantErr {
				t.Errorf("User.HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return // Stop if there's an unexpected error
			}
			// Verify that the password has been hashed, but only if no error was expected
			if !tt.wantErr {
				if u.Password == originalPassword {
					t.Errorf("User.HashPassword() password was not hashed")
				}
				// Optional: You can also verify that the hashed password is a valid bcrypt hash
				if originalPassword != "" { //only check if the original passsword is not empty
					err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(originalPassword))
					if err != nil {
						t.Errorf("User.HashPassword() produced an invalid bcrypt hash: %v", err)
					}
				}
			}
		})
	}
}
// Datastore represents a mock datastore for testing
type Datastore struct {
	Db *sql.DB
}
// ValidateSession validates the session token from the request
func (ds *Datastore) ValidateSession(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, ""
	}
	var userID string
	var expiresAt time.Time
	err = ds.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		return false, ""
	}
	if time.Now().After(expiresAt) {
		return false, ""
	}
	return true, userID
}
// TestValidateSession tests the ValidateSession method of the Datastore struct
func TestValidateSession(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	datastore := &Datastore{Db: db} // Initialize Datastore with mock DB
	testUserID := "test-user-id"
	testSessionToken := "test-session-token"
	futureTime := time.Now().Add(time.Hour)
	pastTime := time.Now().Add(-time.Hour)
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		want1     string
		mockSetup func(mock sqlmock.Sqlmock)
	}{
		{
			name: "No session cookie",
			args: args{
				r: httptest.NewRequest("GET", "/", nil),
			},
			want:  false,
			want1: "",
			mockSetup: func(mock sqlmock.Sqlmock) {},
		},
		{
			name: "Session not found in DB",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
					return req
				}(),
			},
			want:  false,
			want1: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT user_id, expires_at FROM sessions WHERE session_token = ?").
					WithArgs(testSessionToken).
					WillReturnError(sql.ErrNoRows)
			},
		},
		{
			name: "Session expired",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
					return req
				}(),
			},
			want:  false,
			want1: "",
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT user_id, expires_at FROM sessions WHERE session_token = ?").
					WithArgs(testSessionToken).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "expires_at"}).AddRow(testUserID, pastTime))
			},
		},
		{
			name: "Valid session",
			args: args{
				r: func() *http.Request {
					req := httptest.NewRequest("GET", "/", nil)
					req.AddCookie(&http.Cookie{Name: "session_token", Value: testSessionToken})
					return req
				}(),
			},
			want:  true,
			want1: testUserID,
			mockSetup: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery("SELECT user_id, expires_at FROM sessions WHERE session_token = ?").
					WithArgs(testSessionToken).
					WillReturnRows(sqlmock.NewRows([]string{"user_id", "expires_at"}).AddRow(testUserID, futureTime))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup(mock)
			got, got1 := datastore.ValidateSession(tt.args.r) // Use initialized datastore
			if got != tt.want {
				t.Errorf("ValidateSession() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ValidateSession() got1 = %v, want %v", got1, tt.want1)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}