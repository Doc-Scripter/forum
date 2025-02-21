package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func CreateLikesDislikesTable(db *sql.DB) error {

	if db == nil {
		defer db.Close()
		return fmt.Errorf("nil database connection")
	}

	query := `
	CREATE TABLE IF NOT EXISTS likes_dislikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER ,
		comment_id INTEGER ,
		user_uuid TEXT NOT NULL,
		like_dislike TEXT NOT NULL DEFAULT '' ,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid)
		FOREIGN KEY (post_id) REFERENCES posts(post_id)
		FOREIGN KEY (comment_id) REFERENCES comments(comment_id)

	);
	`

	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}
