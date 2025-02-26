package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	e "forum/Error"
	d "forum/database"
	r "forum/routes"
)

func init() {
	if len(os.Args) != 1 {
		log.Fatal("\nUsage: go run main.go")
	}
}

func main() {

	mux, err := r.Routers()
	if err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|main package| ---> {%v}", err))
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "33333"
	}

	reset := "\033[0m"
	boldWhite := "\033[1;37m"
	red := "\033[1;31m"
	blue := "\033[1;34m"
	brown := "\033[0;33m"

	fmt.Println(brown + "╔═══════════════════════════════════════════════════╗" + reset)
	fmt.Println(brown + "║" + red + " 🚀 Server is starting...         " + reset + brown + "                 ║" + reset)
	fmt.Printf(brown+"║ "+boldWhite+"Forum running on port --}  "+blue+"http://localhost:%s"+reset+brown+" ║\n", port)
	fmt.Println(brown + "╚═══════════════════════════════════════════════════╝" + reset)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		e.LOGGER("[ERROR]", fmt.Errorf("|main package| ---> {%v}", err))
		return
	}

	defer d.Db.Close()
}
