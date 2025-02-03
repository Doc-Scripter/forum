package database

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// User represents a user in the database.
type User struct {
	ID   string
	Name string
}

// CreateUser inserts a new user into the database.
func CreateUser(user *User) error {
	// Generate a new UUID for the user
	user.ID = uuid.New().String()

	query := `INSERT INTO users (id, name) VALUES (?, ?)`
	_, err := DB.Exec(query, user.ID, user.Name)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUser retrieves a user by ID.
func GetUser(id string) (*User, error) {
	user := &User{}
	query := `SELECT id, name FROM users WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&user.ID, &user.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}