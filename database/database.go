package database

import (
	"database/sql"
	"fmt"
	"path/filepath"
	e "forum/Error"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// ==== The creation of the database folder and the database file ====
func init() {

	dataFolder := "data"
	databaseFile := "forum.db"

	databaseFolderPath := filepath.Join(dataFolder)
	databaseFilePath := filepath.Join(databaseFolderPath, databaseFile)

	if err := os.MkdirAll(databaseFolderPath, os.ModePerm); err != nil {
		fmt.Println("[DATABASE]: Error creating database folder:", err)
		os.Exit(1)
	}
	
	db_file, err := os.Create(databaseFilePath)
	if err != nil {
		fmt.Println("[DATABASE]: Error creating database file:", err)
		os.Exit(1)
	}
	db_file.Close()

	err = StartDbConnection(databaseFilePath)
	if err != nil {
		e.LOGGER("[ERROR]", err)
		os.Exit(1)
	}
	e.LOGGER("[SUCCESS]: Created the logging file and the database file!", nil)
}



var Db *sql.DB


// ==== This function will starting the connection to the database using the SQLite3 driver that works with CGO =====
func StartDbConnection(database_file_path string) error {


	var err error

	Db, err = sql.Open("sqlite3", database_file_path)
	if err != nil {
		e.LOGGER("[DATABASE ERROR]", err)
	}

	err = Db.Ping()
	if err != nil {
		e.LOGGER("[DATABASE ERROR]", err)
	}

	if err = CreateUsersTable(Db); err != nil {
		e.LOGGER("[ERROR]", err)
	}

	if err = CreateLikesDislikesTable(Db); err != nil {
		e.LOGGER("[ERROR]", err)
	}

	if err = CreateSessionsTable(Db); err != nil {
		e.LOGGER("[ERROR]", err)
	}

	if err = CreatePostsTable(Db); err != nil {
		e.LOGGER("[ERROR]", err)
	}
	if err = CreateCommentsTable(Db); err != nil {
		e.LOGGER("[ERROR]", err)
	}

	e.LOGGER("[SUCCESS]: Connected to the SQLite database!", nil)
	return nil
}
