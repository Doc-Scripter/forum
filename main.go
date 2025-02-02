package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	handler "forum/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler.Home)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	fmt.Printf("Starting server on: %s", port)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
