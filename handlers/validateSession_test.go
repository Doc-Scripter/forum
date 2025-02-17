package handlers

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

type Datastore struct {
	Db *sql.DB
}

func (d *Datastore) ValidateSession(r *http.Request) (bool, string) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false, ""
	}

	var userID string
	var expiresAt time.Time
	err = d.Db.QueryRow("SELECT user_id, expires_at FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ""
		}
		return false, ""
	}

	if time.Now().After(expiresAt) {
		return false, ""
	}

	return true, userID
}

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
