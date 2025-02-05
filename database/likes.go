package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreateLikesTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS likes (
        like_id TEXT PRIMARY KEY,
        post_id TEXT,
        user_id INTEGER,
        FOREIGN KEY (post_id) REFERENCES posts(post_id),
        FOREIGN KEY (user_id) REFERENCES users(id)
    );`
    _, err := db.Exec(query)
    return err
}
