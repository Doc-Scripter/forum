package web

import (
	"log"
	"net/http"
	"text/template"
)

// Define home handler function which writes a byte slice
func Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "page not found", 404)
	}
	t, err := template.ParseFiles("templates/home.html")
	if err != nil {
		log.Fatal("error parsing html")
	}
	t.Execute(w, r)
}
