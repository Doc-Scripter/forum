package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func CreateLikesDislikesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS likes_dislikes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		post_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		like_dislike TEXT NOT NULL CHECK(like_dislike IN ('like', 'dislike')),
		FOREIGN KEY (post_id) REFERENCES posts(post_id),
		FOREIGN KEY (user_id) REFERENCES users(id),
		UNIQUE (post_id, user_id)
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}