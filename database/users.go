package database

import (
    "database/sql"
    "fmt"

    _ "github.com/mattn/go-sqlite3"
)

func CreateUsersTable(db *sql.DB) error {

    if db == nil {
    
        return fmt.Errorf("nil database connection")
    }

    query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
