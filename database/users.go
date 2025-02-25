package database

import (
	"database/sql"
	e "forum/Error"
	_ "github.com/mattn/go-sqlite3"
)

//==== The function will create the users table in the database =====
func CreateUsersTable(db *sql.DB) error {
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
	e.LOGGER("[SUCCESS]: Created the users table", nil)
	return nil
}
