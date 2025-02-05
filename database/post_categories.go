package database

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)

func CreatePostCategoriesTable(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS post_categories (
        post_id TEXT,
        category_id TEXT,
        FOREIGN KEY (post_id) REFERENCES posts(post_id),
        FOREIGN KEY (category_id) REFERENCES categories(category_id),
        PRIMARY KEY (post_id, category_id)
    );`
    _, err := db.Exec(query)
    return err
}
