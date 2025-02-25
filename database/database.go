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

	if _, err := os.Stat(databaseFolderPath); os.IsNotExist(err) {
		if err := os.MkdirAll(databaseFolderPath, os.ModePerm); err != nil {
			fmt.Println("[DATABASE]: Error creating database folder:", err)
			os.Exit(1)
		}
		fmt.Println("[DATABASE]: Database folder created successfully.")
	}

	if _, err := os.Stat(databaseFilePath); os.IsNotExist(err) {
		dbFile, err := os.Create(databaseFilePath)

		if err != nil {
			fmt.Println("[DATABASE]: Error creating database file:", err)
			os.Exit(1)
		}

		dbFile.Close()
		fmt.Println("[DATABASE]: Database file created successfully.")
	} else {
		fmt.Println("[DATABASE]: Database file already exists. Skipping creation.")
	}

	err := StartDbConnection(databaseFilePath)
	if err != nil {
		e.LOGGER("[ERROR]", err)
		os.Exit(1)
	}

	e.LOGGER("[SUCCESS]: Database setup completed successfully!", nil)
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
