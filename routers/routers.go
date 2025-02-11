package routers

import (
	"net/http"
	han "forum/handlers"
	

)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	scriptServer := http.FileServer(http.Dir("web/scripts/"))
	mux.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", scriptServer))

	mux.HandleFunc("/", han.LandingPage)
	mux.HandleFunc("/login", han.Login)
	mux.HandleFunc("/logging", han.AuthenticateUserCredentialsLogin)
	// mux.Handle("/logging", han.AuthMiddleware(http.HandlerFunc(han.AuthenticateUserCredentialsLogin)))
	mux.HandleFunc("/posts", han.PostsHandler)
	mux.HandleFunc("/create-post", han.CreatePostsHandler)
	mux.HandleFunc("/home", han.HomePage)

	mux.HandleFunc("/register", han.Register)
	mux.HandleFunc("/registration", han.RegisterUser)
	mux.HandleFunc("/logout", han.LogoutUser)

	return mux, nil
}
