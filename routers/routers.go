package routers

import(
	"net/http"
	handler "forum/handlers"
)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	scriptServer := http.FileServer(http.Dir("web/scripts/"))
    mux.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", scriptServer))

	mux.HandleFunc("/", handler.Login)
	mux.HandleFunc("/posts", handler.PostsHandler)
	mux.HandleFunc("/create-post", handler.CreatePostsHandler)


	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/registration", handler.RegisterUser)
	
	mux.Handle("/logging", handler.AuthMiddleware(http.HandlerFunc(handler.AuthenticateUserCredentialsLogin)))
	mux.HandleFunc("/logout", handler.LogoutUser)


	return mux, nil
}
