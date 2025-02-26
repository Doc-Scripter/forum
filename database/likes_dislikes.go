package database

import (
	"database/sql"
	"fmt"
	e "forum/Error"

	_ "github.com/mattn/go-sqlite3"
)

// ==== The function will create the likes and the dislikes table in the database =====
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
	e.LOGGER("[SUCCESS]: Created the likes and dislikes table", nil)
	return nil
}
