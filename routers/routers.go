package routers

import(
	"net/http"
	handler "forum/handlers"
)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	mux.HandleFunc("/", handler.LandingPage)
	mux.HandleFunc("/login", handler.Login)
	mux.Handle("/logging", handler.AuthMiddleware(http.HandlerFunc(handler.AuthenticateUserCredentialsLogin)))
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/registration", handler.RegisterUser)
	mux.HandleFunc("/logout", handler.LogoutUser)

	return mux, nil
}
