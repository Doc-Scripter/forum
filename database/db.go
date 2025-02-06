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

	if err=CreateUsersTable(Db);err != nil{
		log.Fatalf("Could not create users table: %e\n", err)
	}

	if err=CreateSessionsTable(Db);err != nil{
		log.Fatalf("Could not create sessions table: %e\n", err)
	}
	
	if err=CreatePostsTable(Db);err != nil{
		log.Fatalf("\nCould not create posts table: %e\n", err)
	}
	if err=CreateCommentsTable(Db);err != nil{
		log.Fatalf("\nCould not create comments table: %e\n", err)
	}


	if err = Db.Ping(); err != nil {
		log.Fatalf("\nCould not test the database connection: %e\n", err)
	}
}

