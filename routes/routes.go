package routers

import(
	"net/http"
	handler "forum/handlers"
)

func Routers() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	scriptServer := http.FileServer(http.Dir("web/scripts/"))
    mux.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", scriptServer))

	mux.HandleFunc("/", handler.LandingPage)
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/logging", handler.AuthenticateUserCredentialsLogin)
	// mux.Handle("/logging", handler.AuthMiddleware(http.HandlerFunc(handler.AuthenticateUserCredentialsLogin)))
	mux.HandleFunc("/posts", handler.PostsHandler)
	mux.HandleFunc("/create-post", handler.CreatePostsHandler)
    mux.HandleFunc("/home", handler.HomePage)
	mux.HandleFunc("/likes", handler.LikePostHandler)
	mux.HandleFunc("/dislikes", handler.DislikePostHandler)

	mux.HandleFunc("/likesComment", handler.LikeCommentHandler)
	mux.HandleFunc("/dislikesComment", handler.DislikeCommentHandler)

	mux.HandleFunc("/addcomment", handler.AddCommentHandler)
	mux.HandleFunc("/comments", handler.CommentHandler)


	mux.HandleFunc("/myPosts", handler.MyPostHandler)
	mux.HandleFunc("/favorites", handler.FavoritesPostHandler)
	
	

	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/registration", handler.RegisterUser)
	mux.HandleFunc("/logout", handler.LogoutUser)

	return mux, nil
}
