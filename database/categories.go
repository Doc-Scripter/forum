package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreateCategoriesTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS categories (
        category_id TEXT PRIMARY KEY,
        name TEXT UNIQUE NOT NULL
    );`
    _, err := db.Exec(query)
    return err
}
