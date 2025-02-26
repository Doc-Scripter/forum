package database

import (
	"database/sql"
	"fmt"
	e "forum/Error"
	_ "github.com/mattn/go-sqlite3"
)

// ==== The function will create the posts table in the database =====
func CreatePostsTable(db *sql.DB) error {

	if db == nil {
		defer db.Close()
		return fmt.Errorf("nil database connection")
	}

	query := `
    CREATE TABLE IF NOT EXISTS posts (
        post_id INTEGER  PRIMARY KEY AUTOINCREMENT DEFAULT 0,
		user_uuid TEXT NOT NULL,
        title TEXT NOT NULL,
        filename TEXT DEFAULT '',
        content TEXT NOT NULL,
        filepath TEXT DEFAULT '',
        comments INTEGER DEFAULT 0,
        category  TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_uuid) REFERENCES users(uuid)
    );`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	e.LOGGER("[SUCCESS]: Created the posts table", nil)
	return nil
}
