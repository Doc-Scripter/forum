package main

import (
	"log"
	"net/http"

	"web-forum/web"
)

func main() {

	http.HandleFunc("/", web.Home)

	fileServer := http.FileServer(http.Dir("./ui/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// Start a new web server
	// 8080 TCP network address to listen on
	log.Println("Starting server on: 8080")
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
