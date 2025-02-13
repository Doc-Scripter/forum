package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func CreateLikesDislikesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS likes_dislikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER  NOT NULL ,
		user_uuid TEXT NOT NULL,
		like_dislike TEXT NOT NULL DEFAULT '' ,
		FOREIGN KEY (user_uuid) REFERENCES users(uuid)
		FOREIGN KEY (post_id) REFERENCES posts(post_id)
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
