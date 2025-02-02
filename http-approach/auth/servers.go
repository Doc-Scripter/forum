package auth

import(
	"log"
	"net/http"
	"text/template"
)

//serve the Homepage
func HomePage(rw http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err!= nil {
		log.Fatal(err)
	}
	t.Execute(rw, nil)
}

//serve the login form
func Login(rw http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/login.html")
	if err!= nil {
		log.Fatal(err)
	}
	t.Execute(rw, nil)
}

//serve the registration form
func Register(rw http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("templates/register.html")
	if err!= nil {
		log.Fatal(err)
	}
	t.Execute(rw, nil)
}