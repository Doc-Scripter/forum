	package routers

	import(
		"net/http"
		han "forum/handlers"
	)

	func Routers() (*http.ServeMux, error) {
		mux := http.NewServeMux()
	
		fileServer := http.FileServer(http.Dir("web/static/"))
		mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	
		mux.HandleFunc("/", han.HomePage)
	
		return mux, nil
	}
	