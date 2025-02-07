package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreateCommentsTable(db *sql.DB) error {

    query := `
    CREATE TABLE IF NOT EXISTS comments (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(post_id)
    );`

    if _, err := db.Exec(query); err != nil {
		return err
	}
	return  nil
}

