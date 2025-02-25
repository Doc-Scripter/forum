package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	e "forum/Error"
	d "forum/database"
	r "forum/routes"
)

func init() {
	if len(os.Args) != 1 {
		log.Fatal("\nUsage: go run main.go")
	}
}

func main() {

	mux, err := r.Routers()
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|main package| ---> {%v}", err))
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "33333"
	}

	fmt.Printf("Starting server on: %s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|main package| ---> {%v}", err))
		return
	}

	defer d.Db.Close()
}
