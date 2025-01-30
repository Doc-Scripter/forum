package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	auth "forum/auth"

	_ "github.com/mattn/go-sqlite3"
)

// initialize a global database variable
var db *sql.DB

func init() {
	var err error

	// Setup your database connection string here
	dbPath := "./mydatabase.db"
	db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Database connection is not alive: %v", err)
	}
}

func main() {
	// Serve homepage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./templates/login.html")
	})

	// Serve the login form
	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/login.html")
	})
	http.HandleFunc("/login", auth.HandleLogin)

	// serve the registration form
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/registration.html")
	})
	http.HandleFunc("/registration", auth.HandleRegistration)

	// Handle the dashboard route
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("templates/dashboard.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// serve the functionality page of the program
	// http.HandleFunc("/library", lib.HandleLibrary)

	// serve the profile page
	// http.HandleFunc("/profile", prf.HandleProfile)

	fmt.Println("Server is running on localhost:9999")
	if err := http.ListenAndServe(":9999", nil); err != nil {
		log.Fatal(err)
	}

	// defer close the open files
	U, M := auth.FileData()
	defer U.Close()
	defer M.Close()
}
