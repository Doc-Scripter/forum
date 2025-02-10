package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreatePostsTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS posts (
        post_id INTEGER  PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        category  TEXT NOT NULL,
		likes INTEGER DEFAULT 0,
		comments TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
    _, err := db.Exec(query)
    return err
}
