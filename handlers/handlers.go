package web

import (
	"log"
	"net/http"
	"text/template"
)

// Define home handler function which writes a byte slice
func HomePage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		http.Error(w, "page not found", 404)
	}

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}

	t, err := template.ParseFiles("web/templates/home.html")
	if err != nil {
		log.Fatalf("\nError parsing html\n")
	}
	t.Execute(w, nil)
}
