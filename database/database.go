package database

import (
	"database/sql"
	"fmt"
	"log"

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
		log.Fatalf("\nCould not create User table: %e\n", err)
	}

	if err = CreateLikesDislikesTable(Db); err != nil {
		log.Fatalf("\nCould not create Likes and Dislikes table: %e\n", err)
	}

	if err = CreateSessionsTable(Db); err != nil {
		log.Fatalf("\nCould not create sessions table: %e\n", err)
	}
	
	if err = CreatePostsTable(Db); err != nil {
		log.Fatalf("\nCould not create posts table: %e\n", err)
	}
	if err = CreateCommentsTable(Db); err != nil {
		log.Fatalf("\nCould not create comments table: %e\n", err)
	}

	fmt.Println("Connected to SQLite database successfully!")
	return nil
}
