package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// startDbConnectionWithPath establishes a database connection with a specified path
// This is a helper function that can be used for testing with different paths
func startDbConnectionWithCustomPath(dbPath string) error {
	var err error

	Db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	err = Db.Ping()
	if err != nil {
		return err
	}

	if err = CreateUsersTable(Db); err != nil {
		return err
	}

	if err = CreateLikesDislikesTable(Db); err != nil {
		return err
	}

	if err = CreateSessionsTable(Db); err != nil {
		return err
	}

	if err = CreatePostsTable(Db); err != nil {
		return err
	}
	if err = CreateCommentsTable(Db); err != nil {
		return err
	}

	return nil
}

// ==== This function will start the connection to the database=============
func StartDbConnection() error {
	err := startDbConnectionWithCustomPath("forum.db")
	if err != nil {
		return err
	}
	
	fmt.Println("Connected to the SQLite database!")
	return nil
}