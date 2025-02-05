package database

import(
	"log"
	"database/sql"
)

var Db *sql.DB

func StartDbConn() {

	var err error

	Db, err = sql.Open("sqlite3", "forum_database.db")
	if err != nil {
		log.Fatal(err)
	}

	if err = Db.Ping(); err != nil {
		log.Fatalf("\nCould not test the database connection: %e\n", err)
	}
}

