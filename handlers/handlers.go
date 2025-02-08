package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// serve the Homepage
func HomePage(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/index.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(rw, nil)
}

// serve the login form
func Login(rw http.ResponseWriter, req *http.Request) {

	if bl, _ := ValidateSession(req); bl {
		HomePage(rw, req)
	} else if !bl {

		tmpl, err := template.ParseFiles("web/templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(rw, nil)
	}
}

// serve the registration form
func Register(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(rw, nil)
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 404)
	}
	rows, err := Db.Query("SELECT * FROM posts")
	if err != nil {
		http.Error(w, "could not get posts", 500)
	}
	var posts []map[string]interface{}
	// posts:=[]post

	// var posts []post
	for rows.Next() {
		var eachPost map[string]interface{}
		err := rows.Scan(eachPost["created_at"], eachPost["category"], eachPost["likes"], eachPost["comments"])
		fmt.Println(eachPost)
		if err != nil {
			http.Error(w, "could not get post", http.StatusInternalServerError)
		}
		posts = append(posts, eachPost)
	}
	fmt.Println(posts)
	postsJson, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "could not marshal posts", http.StatusInternalServerError)
	}

	w.Write(postsJson)
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}
	fmt.Println("creating post")
	r.ParseForm()
	category := r.FormValue("category")
	content := r.FormValue("content")

	_, err := Db.Exec("INSERT INTO posts (category, content) VALUES ($1, $2)", category, content)
	if err != nil {
		http.Error(w, "could not insert post", http.StatusInternalServerError)
		return
	}
	fmt.Println("post created")
}