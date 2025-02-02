package main

import(
	"os"
	"net/http"
	"log"
	"database/sql"
	auth "forum/auth"
)

var Db *sql.DB


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
	mux.HandleFunc("/signing", auth.Login)

	mux.HandleFunc("/signup", auth.Register)
	mux.HandleFunc("/registration", auth.RegisterUser)
	mux.Handle("/login", auth.AuthMiddleware(http.HandlerFunc(auth.AuthenticateUserCredentialsLogin)))
	mux.Handle("/logout", auth.AuthMiddleware(http.HandlerFunc(auth.LogoutUser)))

	port :=  os.Getenv("PORT")
	if port == "" {
		port = "41532"
	}

	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
	defer Db.Close()
}
