package auth

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var Db *sql.DB

//==========This function creates a 'users' table in the SQLite database===========
func CreateUserTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		uuid TEXT UNIQUE NOT NULL,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);
	`

	_, err := Db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}

	fmt.Println("Table 'users' created successfully!")
}

//==========This function creates a 'sessions' table in the SQLite database===========
func CreateSessionTable() {
	query := `
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
    	user_id INTEGER NOT NULL,
    	session_token TEXT UNIQUE NOT NULL,
    	expires_at DATETIME NOT NULL,
    	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);`

	_, err := Db.Exec(query)
	if err != nil {
		log.Fatal("Error creating table:", err)
	}
	
	fmt.Println("Table 'sessions' created successfully!")
}


//============Starting the connection to the database=============
func StartDBConnection() error {

	var err error

	Db, err = sql.Open("sqlite3", "mydatabase.db")

	if err != nil {
		return err
	}
	
	err = Db.Ping()
	if err != nil {
		return err
	}

	fmt.Println("Connected to SQLite database successfully!")
	return nil

}