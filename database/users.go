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
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );`
    _, err := db.Exec(query)
    return err
}
