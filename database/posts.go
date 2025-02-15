package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CreatePostsTable(db *sql.DB) error {
	query := `
    CREATE TABLE IF NOT EXISTS posts (
        post_id INTEGER  PRIMARY KEY AUTOINCREMENT DEFAULT 0,
		user_uuid TEXT NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        comments INTEGER DEFAULT 0,
        category  TEXT NOT NULL,
		likes INTEGER DEFAULT 0,
        dislikes INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_uuid) REFERENCES users(uuid)
    );`
	_, err := db.Exec(query)
	return err
}
