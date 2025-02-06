package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreateSessionsTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
        session_token TEXT UNIQUE NOT NULL,
        expires_at DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`
    
	_, err := db.Exec(query)
	return err
}
