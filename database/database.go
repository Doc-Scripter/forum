package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// ==== This function will starting the connection to the database=============
func StartDbConnection() error {
	var err error

	Db, err = sql.Open("sqlite3", "data/forum.db")
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

	fmt.Println("Connected to the SQLite database!")
	return nil
}
