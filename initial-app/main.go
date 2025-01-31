package main

import (
	"html/template"
	"log"
	"net/http"
)

// Define home handler function which writes a byte slice
func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "page not found", 404)
	}
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal("error parsing html")
	}
	t.Execute(w, r)
}

func main() {
	// Initialize a new servemux then register home function as handler for "/"
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// Start a new web server
	// 8080 TCP network address to listen on
	log.Println("Starting server on: 8080")
	err := http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
