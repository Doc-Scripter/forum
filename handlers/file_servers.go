package handlers

import (
	"net/http"
	"text/template"
	"fmt"

	e "forum/Error"
	m "forum/models"
	u "forum/utils"
)

/*
	==== The function handler serves the error page with relevant error messages. The function will log the actual problem (Error) to the

logging file and then, serve the error page with the error message and the error code in the ErrorData object ====
*/
func ErrorPage(Error error, ErrorData m.ErrorData, w http.ResponseWriter, r *http.Request) {
	
	e.LOGGER("[ERROR]", Error)
	tmpl, err := template.ParseFiles("./web/templates/error.html")
	
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		e.LOGGER("[ERROR]", fmt.Errorf("|error page server| ---> {%v}", err))
		return
	}
	
	if err = tmpl.Execute(w, ErrorData); err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|error page server| ---> {%v}", err))
		return
	}
}

// ==== The function handler serves the landing page of the web application ====
func LandingPage(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/" {

		if bl, _ := u.ValidateSession(r); bl {
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}

		if r.Method != http.MethodGet {
			ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
			return
		}

		tmpl, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			ErrorPage(fmt.Errorf("|landing page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		if err = tmpl.Execute(w, nil); err != nil {
			ErrorPage(fmt.Errorf("|landing page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

	} else {
		ErrorPage(nil, m.ErrorsData.PageNotFound, w, r)
		return
	}
}

// ==== The function handler serves the home page of the web application ====
func HomePage(w http.ResponseWriter, r *http.Request) {

	if bl, _ := u.ValidateSession(r); !bl {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		ErrorPage(fmt.Errorf("|home page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/home.html")
	if err != nil {
		ErrorPage(fmt.Errorf("|home page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	Profile.Category = m.Category

	if err = tmpl.Execute(w, Profile); err != nil {
		ErrorPage(fmt.Errorf("|home page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
}

// ==== This function will serve the login form of the application ====
func Login(w http.ResponseWriter, r *http.Request) {
	if bl, _ := u.ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/login.html")
	if err != nil {
		ErrorPage(fmt.Errorf("|login page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
	if err = tmpl.Execute(w, nil); err != nil {
		ErrorPage(fmt.Errorf("|login page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
}

// ==== This function will serve the registration form of the application ====
func Register(w http.ResponseWriter, r *http.Request) {

	if bl, _ := u.ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/register.html")
	if err != nil {
		ErrorPage(fmt.Errorf("|register page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
	if err = tmpl.Execute(w, nil); err != nil {
		ErrorPage(fmt.Errorf("|register page server| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
}
