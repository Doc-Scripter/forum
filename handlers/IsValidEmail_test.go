package handlers

import (
	"testing"
)

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
