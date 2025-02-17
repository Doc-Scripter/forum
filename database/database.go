package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// ============Starting the connection to the database=============
func StartDbConnection() error {
	var err error

	Db, err = sql.Open("sqlite3", "forum.db")
	if err != nil {
		return err
	}

	err = Db.Ping()
	if err != nil {
		return err
	}

	if err = CreateUsersTable(Db); err != nil {
		fmt.Println("Could not create User table: ", err)
	}

	if err = CreateLikesDislikesTable(Db); err != nil {
		fmt.Println("Could not create Likes and Dislikes table: ", err)
	}

	if err = CreateSessionsTable(Db); err != nil {
		fmt.Println("Could not create sessions table: %", err)
	}

	if err = CreatePostsTable(Db); err != nil {
		fmt.Println("Could not create posts table: ", err)
	}
	if err = CreateCommentsTable(Db); err != nil {
		fmt.Println("Could not create comments table: ", err)
	}
	if err = CreateImageTable(Db); err != nil {
		fmt.Println("Could not create image table")
	}
	fmt.Println("Connected to SQLite database successfully!")
	return nil
}
