package auth

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {

	//filedata
	Userfile, _ := FileData()

	// Handle login submissions
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	// Get username and password
	username := r.FormValue("username")
	password := r.FormValue("password")
	var us, psw string

	// Validate the credentials
	scanner := bufio.NewScanner(Userfile)
	authenticated := false

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) < 2 {
			continue
		}
		us = parts[0][2:]
		psw = parts[1][2:]

		if us == username && psw == password {
			authenticated = true
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %s", err)
	}

	if authenticated {
		fmt.Println("Authenticated")
		// Redirect to the dashboard on successful login
		http.Redirect(w, r, "/dashboard", http.StatusMovedPermanently)
	} else {

		// If authentication fails
		fmt.Fprintf(w, "Invalid username or password. Try loggin in.\n")
	}
}
