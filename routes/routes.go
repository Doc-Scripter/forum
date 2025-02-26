package routers

import (
	handler "forum/handlers"
	"net/http"
)

// ====== The routers function will handle trafficking and map each endpoint to its corresponding handler ====
func Routers() (*http.ServeMux, error) {

	mux := http.NewServeMux()

	// ==== static files servers  (CSS, JS) ====
	fileServer := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	scriptServer := http.FileServer(http.Dir("web/scripts/"))
	mux.Handle("/web/scripts/", http.StripPrefix("/web/scripts/", scriptServer))

	// ==== authentication endpoints =====
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/logging", handler.AuthenticateUserCredentialsLogin)
	mux.HandleFunc("/register", handler.Register)
	mux.HandleFunc("/registration", handler.RegisterUser)
	mux.HandleFunc("/logout", handler.LogoutUser)

	// === endpoints available for all users ===
	mux.HandleFunc("/posts", handler.PostsHandler)
	mux.HandleFunc("/comments", handler.CommentHandler)
	mux.HandleFunc("/image/web/uploads/", handler.ImageHandler)

	// === endpoints only available for logged in users ====
	mux.Handle("/create-post", handler.AuthMiddleware(http.HandlerFunc(handler.CreatePostsHandler)))
	mux.Handle("/likes", handler.AuthMiddleware(http.HandlerFunc(handler.LikePostHandler)))
	mux.Handle("/dislikes", handler.AuthMiddleware(http.HandlerFunc(handler.DislikePostHandler)))
	mux.Handle("/likesComment", handler.AuthMiddleware(http.HandlerFunc(handler.LikeCommentHandler)))
	mux.Handle("/dislikesComment", handler.AuthMiddleware(http.HandlerFunc(handler.DislikeCommentHandler)))
	mux.Handle("/addcomment", handler.AuthMiddleware(http.HandlerFunc(handler.AddCommentHandler)))
	mux.Handle("/myPosts", handler.AuthMiddleware(http.HandlerFunc(handler.MyPostHandler)))
	mux.Handle("/favorites", handler.AuthMiddleware(http.HandlerFunc(handler.FavoritesPostHandler)))

	// ==== endpoint navigation in the application ====
	mux.HandleFunc("/home", handler.HomePage)
	mux.HandleFunc("/", handler.LandingPage)

	return mux, nil
}
