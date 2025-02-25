package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

// startDbConnectionWithPath establishes a database connection with a specified path
// This is a helper function that can be used for testing with different paths
func startDbConnectionWithCustomPath(dbPath string) error {
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

<<<<<<< HEAD
	return nil
}

// ==== This function will start the connection to the database=============
func StartDbConnection() error {
	err := startDbConnectionWithCustomPath("forum.db")
	if err != nil {
		return err
	}
	
=======
>>>>>>> master
	fmt.Println("Connected to the SQLite database!")
	return nil
}