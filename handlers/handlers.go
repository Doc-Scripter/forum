package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
)

// serve the login form
func LandingPage(rw http.ResponseWriter, req *http.Request) {

	if bl, _ := ValidateSession(req); bl {
		HomePage(rw, req)
	} else if !bl {

		tmpl, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(rw, nil)
	}
}

type Post struct {
	CreatedAt string `json:"created_at"`
	Category  string `json:"category"`
	Likes     int    `json:"likes"`
	Comments  *string    `json:"comments"`
	Content   string `json:"content"`
}
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

		tmpl, err := template.ParseFiles("./web/templates/login.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(rw, nil)
	}
}

// serve the registration form
func Register(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/register.html")
	if err != nil {
		log.Fatal(err)
	}
	tmpl.Execute(rw, nil)
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 404)
	}
	rows, err := Db.Query("SELECT category,content,comments,created_at,likes FROM posts")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not get posts", http.StatusInternalServerError)
	}
	defer rows.Close()
	// var posts []map[string]interface{}
	// posts:=[]post

	var posts []Post
	for rows.Next() {
		var eachPost Post
		var comments sql.NullString
		err := rows.Scan(&eachPost.Category, &eachPost.Content,&comments, &eachPost.CreatedAt, &eachPost.Likes)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get post", http.StatusInternalServerError)
			return
		}
		if comments.Valid {
			eachPost.Comments = &comments.String
		} else {
			eachPost.Comments = nil
		}
		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "could not marshal posts", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
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
	HomePage(w, r)
	fmt.Println("post created")
}