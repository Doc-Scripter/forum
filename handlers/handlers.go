package handlers

import (
	"fmt"
	"log"
	"net/http"
	"database/sql"
	"encoding/json"
	"text/template"
	d "forum/database"
)

// serve the login form
func LandingPage(rw http.ResponseWriter, req *http.Request) {
	if bl, _ := ValidateSession(req); bl {
		http.Redirect(rw, req, "/home", http.StatusSeeOther)
	} else if !bl {

		tmpl, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(rw, nil)
	}
}

type ProfileData struct {
	Username string
	Email    string
}

type Post struct {
	CreatedAt string  `json:"created_at"`
	Category  string  `json:"category"`
	Likes     int     `json:"likes"`
	Comments  *string `json:"comments"`
	Content   string  `json:"content"`
}

// serve the Homepage

func getUserDetails(r *http.Request) (ProfileData, error) {
	var PD ProfileData

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Profile Section: No session cookie found")
		return ProfileData{}, err
	}

	var userID string

	err = d.Db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println("Session not found in DB:", err)
		return ProfileData{}, err
	}

	query := `
	SELECT  username, email FROM users WHERE id = ?`

	err = d.Db.QueryRow(query, userID).Scan(&PD.Username, &PD.Email)
	if err != nil {
		return ProfileData{}, err
	}
	return PD, nil
}

func HomePage(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		fmt.Fprint(rw, "Bad request", http.StatusBadRequest)
	}

	Pd, err := getUserDetails(req)
	fmt.Println(Pd)
	if err != nil {
		fmt.Printf("Profile Section: %e\n", err)
	}
	tmpl, err := template.ParseFiles("./web/templates/home.html")
	if err != nil {
		log.Fatal(err)
	}

	tmpl.Execute(rw, Pd)
}

// serve the login form
func Login(rw http.ResponseWriter, req *http.Request) {
	if bl, _ := ValidateSession(req); bl {
		http.Redirect(rw, req, "/home", http.StatusSeeOther)
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
	rows, err := d.Db.Query("SELECT category,content,comments,created_at,likes FROM posts")
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
		err := rows.Scan(&eachPost.Category, &eachPost.Content, &comments, &eachPost.CreatedAt, &eachPost.Likes)
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

	_, err := d.Db.Exec("INSERT INTO posts (category, content) VALUES ($1, $2)", category, content)
	if err != nil {
		http.Error(w, "could not insert post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	fmt.Println("post created")
}
