package database

import (
	"database/sql"
	"fmt"
	e "forum/Error"
	_ "github.com/mattn/go-sqlite3"
)

//==== The function will create the sessions table in the database =====
func CreateSessionsTable(db *sql.DB) error {

	if db == nil {
		defer db.Close()
		return fmt.Errorf("nil database connection")
	}

	query := `
    CREATE TABLE IF NOT EXISTS sessions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
        session_token TEXT UNIQUE NOT NULL,
        expires_at DATETIME NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	e.LOGGER("[SUCCESS]: Created the Sessions table", nil)
	return nil
}
