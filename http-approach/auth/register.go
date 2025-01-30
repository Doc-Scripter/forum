package auth

import (
	"fmt"
	"net/http"
)

//Handle registration
func HandleRegistration(w http.ResponseWriter, r *http.Request) {

	//file data
	_, Membersfile := FileData()

	if r.Method == http.MethodGet {
		// Serve the registration page
		http.ServeFile(w, r, "../templates/registration.html")
		return
	}

	if r.Method == http.MethodPost {

		// Parse form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Validate form inputs
		if username == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Check if the username already exists
		if CredentialExists(username) {
			http.Error(w, fmt.Sprintf("Username %s is already taken. Please choose another.", username), http.StatusConflict)
			return
		}

		// Add user data to the members file
		userData := fmt.Sprintf("u-%s e-%s p-%s\n", username, email, password)
		fmt.Println(userData)

		_, err := Membersfile.WriteString(userData)
		if err != nil {
			http.Error(w, "Failed to save user data", http.StatusInternalServerError)
			return
		}

		// Add user to the system
		if AddUser(username, password) {

			// Redirect to the home page
			http.Redirect(w, r, "/signin", http.StatusSeeOther)
		} else {
			http.Error(w, "Failed to add user to the system", http.StatusInternalServerError)
		}
		return
	}

	// Handle unsupported methods
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
