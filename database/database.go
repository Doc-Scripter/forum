package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// ==== This function will starting the connection to the database=============
func StartDbConnection() error {
	var err error

	if _,err := os.Stat("data"); os.IsNotExist(err) {
		if err = os.Mkdir("data", 0766); err != nil {
			return err
		}
	}

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
