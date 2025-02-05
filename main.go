package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	datab "forum/database"
	handler "forum/handlers"
)

func init() {

	//check the number of arguments
	if len(os.Args) != 1 {
		log.Fatal("\nUsage: go run main.go")
	}

	//start the database connection
	datab.StartDbConn()
	
	//create the tables in the db
	datab.CreateUsersTable(datab.Db)
	datab.CreateSessionsTable(datab.Db)
	datab.CreateCategoriesTable(datab.Db)
	datab.CreatePostsTable(datab.Db)
	datab.CreateCommentsTable(datab.Db)
	datab.CreateLikesTable(datab.Db)
	datab.CreateDislikesTable(datab.Db)
	datab.CreatePostCategoriesTable(datab.Db)
}


func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", handler.HomePage)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Starting server on: %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
	
	defer datab.Db.Close()
}
