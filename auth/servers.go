package auth

import(
	"log"
	"net/http"
	"text/template"
)

//serve the Homepage
func HomePage(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err!= nil {
		log.Fatal(err)
	}
	tmpl.Execute(rw, nil)
}

//serve the login form
func Login(rw http.ResponseWriter, req *http.Request) {
	
	if bl, _ := ValidateSession(req); bl {
		HomePage(rw, req)
	}else if !bl {

		
		tmpl, err := template.ParseFiles("templates/login.html")
		if err!= nil {
			log.Fatal(err)
		}
		tmpl.Execute(rw, nil)
	}
}

//serve the registration form
func Register(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("templates/register.html")
	if err!= nil {
		log.Fatal(err)
	}
	tmpl.Execute(rw, nil)
}
