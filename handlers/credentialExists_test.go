package handlers

import (
	"database/sql"
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

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
