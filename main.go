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

	 // Serve static files
	fileServer := http.FileServer(http.Dir("./web/static"))
	mux.Handle("/web/static/", http.StripPrefix("/web/static/", fileServer))

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
