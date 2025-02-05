package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreateUsersTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT UNIQUE NOT NULL,
        username TEXT NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`
    if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
