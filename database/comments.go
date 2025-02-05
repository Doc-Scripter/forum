package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreateCommentsTable(db *sql.DB) error {

    query := `
    CREATE TABLE IF NOT EXISTS comments (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts(post_id)
    );`

    if _, err := db.Exec(query); err != nil {
		return err
	}
	return  nil
}

func CreateCommentsReactionLike(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS comments_reaction (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        reaction TEXT NOT NULL,
        FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`
    if _, err := db.Exec(query); err != nil {
        return err
    }
    return nil
}

func CreateCommentsReactionDisLike(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS comments_reaction (
        comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
        reaction TEXT NOT NULL,
        FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );`
    if _, err := db.Exec(query); err != nil {
        return err
    }
    return nil
}