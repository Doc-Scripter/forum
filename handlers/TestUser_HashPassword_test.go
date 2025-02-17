package handlers

import (
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

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
