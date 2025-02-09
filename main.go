package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	han "forum/handlers"
	r "forum/routers"
)

func init() {
	// check the number of arguments
	// check the number of arguments
	if len(os.Args) != 1 {
		log.Fatal("\nUsage: go run main.go")
	}

	// start the database connection
	han.StartDbConnection()
}

func main() {
	mux, err := r.Routers()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "33333"
	}

	fmt.Printf("Starting server on: %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}

	defer han.Db.Close()
}
