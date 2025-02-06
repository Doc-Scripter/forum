package main

import (
	"fmt"
	"forum/auth"
	"forum/database"
	handler "forum/handlers"
	"log"
	"net/http"
	"os"
)

func init() {

	//check the number of arguments
	if len(os.Args) != 1 {
		log.Fatal("\nUsage: go run main.go")
	}

	//start the database connection
	database.StartDbConn()

}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", auth.Login)
	mux.HandleFunc("/posts", handler.PostsHandler)

	mux.HandleFunc("/register", auth.Register)

	mux.HandleFunc("/login", auth.Login)

	mux.HandleFunc("/registration", auth.RegisterUser)

	// Serve static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/web/static/", http.StripPrefix("/web/static/", fileServer))

	mux.Handle("/logging", auth.AuthMiddleware(http.HandlerFunc(auth.AuthenticateUserCredentialsLogin)))

	mux.HandleFunc("/logout", auth.LogoutUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "41532"
	}

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}

	port = os.Getenv("PORT")
	if port == "" {
		port = "33333"
	}

	fmt.Printf("Starting server on: %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}

	defer database.Db.Close()
}
