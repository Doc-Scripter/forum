package main

import (
	"log"
	"net/http"

	"web-forum/web"
)

func main() {
	// Initialize a new servemux then register home function as handler for "/"
	mux := http.NewServeMux()
	mux.HandleFunc("/", web.Home)

	// Start a new web server
	// 8080 TCP network address to listen on
	log.Println("Starting server on: 8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
