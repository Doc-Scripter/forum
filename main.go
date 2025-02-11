package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	d "forum/database"
	r "forum/routers"
)

func init() {
	// check the number of arguments
	if len(os.Args) != 1 {
		log.Fatal("\nUsage: go run main.go")
	}

	// start the database connection
	err := d.StartDbConnection()
	if err != nil {
		log.Fatal(err)
	}
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

	defer d.Db.Close()
}
