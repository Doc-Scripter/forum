package main

import(
	"os"
	"net/http"
	"log"
	auth "forum/auth"
)




func init() {

	if len(os.Args) != 1 {
		log.Fatal("\nUsage: go run main.go")
	}

	auth.StartDBConnection()
	auth.CreateUserTable()
	auth.CreateSessionTable()
}

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/", auth.Login)
	mux.HandleFunc("/register", auth.Register)

	mux.HandleFunc("/login", auth.Login)

	mux.HandleFunc("/registration", auth.RegisterUser)

	// mux.Handle("/", auth.AuthMiddleware(http.HandlerFunc(auth.Login)))

	// mux.Handle("/login", auth.AuthMiddleware(http.HandlerFunc(auth.Login)))

	// mux.Handle("/register", auth.AuthMiddleware(http.HandlerFunc(auth.Register)))
	
	mux.Handle("/logging", auth.AuthMiddleware(http.HandlerFunc(auth.AuthenticateUserCredentialsLogin)))

	mux.HandleFunc("/logout", auth.LogoutUser)

	port :=  os.Getenv("PORT")
	if port == "" {
		port = "41532"
	}

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
	defer auth.Db.Close()
}
