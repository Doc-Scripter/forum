package database

import (
	"database/sql" // Imported for side effects
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB is a global variable to hold the database connection.
var DB *sql.DB

// InitDB initializes the SQLite database and creates tables if they don't exist.
func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./database.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables
	if err := createUserTable(); err != nil {
		return fmt.Errorf("failed to create user table: %w", err)
	}
	if err := createProductTable(); err != nil {
		return fmt.Errorf("failed to create product table: %w", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// createUserTable creates the users table with a UUID primary key.
func createUserTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL
	);`
	_, err := DB.Exec(query)
	return err
}

// createProductTable creates the products table with a UUID primary key.
func createProductTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		price REAL NOT NULL
	);`
	_, err := DB.Exec(query)
	return err
}