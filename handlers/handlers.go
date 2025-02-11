package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	d "forum/database"
	e "forum/Error"
)

type ErrorData struct {
	Msg string
	Code int
}

type ProfileData struct {
	UUID string
	Username string
	Email    string
}

type Post struct {
	CreatedAt time.Time `json:"created_at"`
	Category  string    `json:"category"`
	Likes     int       `json:"likes"`
	Title     string    `json:"title"`
	Dislikes  int       `json:"dislikes"`
	Comments  string    `json:"comments"`
	Content   string    `json:"content"`
	Post_ID   int       `json:"post_id"`
}

var InternalError = ErrorData{
	Msg: "An unexpected error occurred. The error seems to be on our end. Hang tight",
    Code: http.StatusInternalServerError,
}

var PageNotFound = ErrorData{
	Msg: "The Page you are trying to access does not seem to exist",
    Code: http.StatusNotFound,
}

var BadRequest = ErrorData{
	Msg: "That's a Bad Request",
    Code: http.StatusBadRequest,
}

var Unauthorized = ErrorData{
    Msg: "You are not authorized to view this page",
    Code: http.StatusUnauthorized,
}

var Forbidden = ErrorData{
    Msg: "You are not allowed to perform this action",
    Code: http.StatusForbidden,
}

var MethodNotAllowed = ErrorData{
	Msg: "The HTTP method you used is not allowed for this route",
    Code: http.StatusMethodNotAllowed,
}

func ErrorPage(Error error, Err ErrorData, w http.ResponseWriter, r *http.Request) {

	e.LogError(Error)
	tmpl, err := template.ParseFiles("./web/templates/error.html")
		if err != nil {
			e.LogError(err)
		}
		if err = tmpl.Execute(w, Err); err != nil {
			e.LogError(err)
		}
}

// serve the login form
func LandingPage(w http.ResponseWriter, r *http.Request) {
	if bl, _ := ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if !bl {
		
		if r.Method == http.MethodGet{
			ErrorPage(nil, InternalError, w, r)
		}
		tmpl, err := template.ParseFiles("./web/templates/index.html")
		errD := InternalError
		if err != nil {
			ErrorPage(err, errD, w, r)
			return
		}

		if err = tmpl.Execute(w, nil); err != nil {
			ErrorPage(err, errD, w, r)
			return
		}
	}
}

// serve the Homepage

func getUserDetails(w http.ResponseWriter, r *http.Request) (ProfileData, error) {
	var PD ProfileData

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Profile Section: No session cookie found")
		ErrorPage(err, InternalError, w, r)
		return ProfileData{}, err
	}

	var userID string

	err = d.Db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println("Session not found in DB:", err)
		return ProfileData{}, err
	}

	query := `
	SELECT  uuid, username, email FROM users WHERE id = ?`

	err = d.Db.QueryRow(query, userID).Scan(&PD.UUID, &PD.Username, &PD.Email)
	if err != nil {
		return ProfileData{}, err
	}
	return PD, nil
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, MethodNotAllowed, w, r)
		return
	}

	Pd, err := getUserDetails(w, r)
	if err != nil {
		ErrorPage(err, InternalError, w, r)
		return
	}
	tmpl, err := template.ParseFiles("./web/templates/home.html")
	if err != nil {
		ErrorPage(err, InternalError, w, r)
		return
	}

	if err = tmpl.Execute(w, Pd); err != nil {
		ErrorPage(err, InternalError, w, r)
		return
	}
}

// serve the login form
func Login(w http.ResponseWriter, r *http.Request) {
	if bl, _ := ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else if !bl {

		tmpl, err := template.ParseFiles("./web/templates/login.html")
		if err != nil {
			ErrorPage(err, InternalError, w, r)
			return
		}
		if err = tmpl.Execute(w, nil); err != nil {
			ErrorPage(err, InternalError, w, r)
			return
		}
	}
}

// serve the registration form
func Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrorPage(nil, MethodNotAllowed, w, r)
        return
	}

	tmpl, err := template.ParseFiles("./web/templates/register.html")
	if err != nil {
		log.Fatal(err)
	}
	if err = tmpl.Execute(w, nil); err != nil {
		ErrorPage(err, InternalError, w, r)
        return
	}
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, InternalError, w, r)
		return
	}
	rows, err := d.Db.Query("SELECT category,title,content,created_at,post_id FROM posts")
	if err != nil {
		ErrorPage(err, InternalError, w, r)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var eachPost Post
		// var comments sql.NullString
		err := rows.Scan(&eachPost.Category, &eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, &eachPost.Post_ID)
		if err != nil {
			fmt.Println(err)
			ErrorPage(err, InternalError, w, r)
			return
		}

		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", eachPost.Post_ID).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			ErrorPage(err, InternalError, w, r)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", eachPost.Post_ID).Scan(&dislikeCount)
		if err != nil {
			fmt.Println(err)
			ErrorPage(err, InternalError, w, r)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		// http.Error(w, "could not marshal posts", http.StatusInternalServerError)
		ErrorPage(err, InternalError, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// http.Error(w, "method not allowed", http.StatusBadRequest)
		ErrorPage(nil, MethodNotAllowed, w, r)
		return
	}
	fmt.Println("creating post")
	r.ParseForm()
	category := r.FormValue("category")
	content := r.FormValue("content")
	title := r.FormValue("title")

	_, err := d.Db.Exec("INSERT INTO posts (category, content, title) VALUES ($1, $2, $3)", category, content, title)
	fmt.Println(err)
	if err != nil {
		// http.Error(w, "could not insert post", http.StatusInternalServerError)
		ErrorPage(err, MethodNotAllowed, w, r)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
	fmt.Println("post created")
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	
	if r.Method != http.MethodPost {
		// http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		ErrorPage(nil, MethodNotAllowed, w, r)
		return
	}

	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}
	fmt.Println(string(str))
	err := json.Unmarshal(str, &postID)

	if err != nil {
		fmt.Println("could not unmarshal post id")
		ErrorPage(err, BadRequest, w, r)
		return
	}


	PostNumID, err := strconv.Atoi(postID.Post_id)
	if err != nil {
		// http.Error(w, "Invalid post id", http.StatusBadRequest)
		ErrorPage(err, BadRequest, w, r)
		return
	}

	// Check if the user has already liked or disliked the post
	var likeDislike string

	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND post_id = ?", PostNumID).Scan(&likeDislike)
	if err == sql.ErrNoRows {
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id ", PostNumID).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				err = d.Db.QueryRow("SELECT * FROM likes_dislikes WHERE post_id = ?", PostNumID).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,post_id) VALUES ('like',?)", PostNumID)
					if err != nil {
						fmt.Println("Failed to like post", err)
						http.Error(w, "Failed to like post", http.StatusInternalServerError)
						return
					}
				} else {
					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? ", PostNumID)
					if err != nil {
						fmt.Println("Failed to like post", err)
						ErrorPage(err, InternalError, w, r)
						return
					}
				}

			} else {

				fmt.Println("Failed to query post", err)
				ErrorPage(err, InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? ", PostNumID)
		if err != nil {
			fmt.Println("Failed to like post", err)
			ErrorPage(err, InternalError, w, r)
			return
		}

	} else if err != nil {

		fmt.Println("Failed to check if user has liked post", err)
		ErrorPage(err, InternalError, w, r)
		return
	} else if likeDislike == "like" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? ", PostNumID)
		if err != nil {
			fmt.Println("Failed to minus like", err)
			ErrorPage(err, InternalError, w, r)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	
	str, err := io.ReadAll(r.Body)
	if err != nil {
        ErrorPage(err, BadRequest, w, r)
		return
    }
	var postID struct {
		Post_id string `json:"post_id"`
	}

	err = json.Unmarshal(str, &postID)
	if err != nil {
	
		http.Error(w, "could not unmarshal post id", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		ErrorPage(nil, BadRequest, w, r)
		return
	}

	PostNumID, err := strconv.Atoi(postID.Post_id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid post id", http.StatusBadRequest)
		return
	}

	// Check if the user has already liked or disliked the post
	var likeDislike string
	
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ?", PostNumID).Scan(&likeDislike)
	if err == sql.ErrNoRows {
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'slike' AND post_id ", PostNumID).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("had  liked it")
				err = d.Db.QueryRow("SELECT * FROM likes_dislikes WHERE post_id = ?", PostNumID).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,post_id) VALUES ('dislike',?)", PostNumID)
					if err != nil {
						fmt.Println("Failed to dislike post", err)
						ErrorPage(err, InternalError, w, r)
						return
					}
				} else {
					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? ", PostNumID)
					if err != nil {
						fmt.Println("Failed to dislike post", err)
						ErrorPage(err, InternalError, w, r)
						return
					}
				}

			} else {

				fmt.Println("Failed to query post", err)
				ErrorPage(err, InternalError, w, r)
				return
			}
		}

	} else if err != nil {
		fmt.Println("Failed to check if user has disliked post", err)
		http.Error(w, "Failed to check if user has disliked post", http.StatusInternalServerError)
		return
	} else if likeDislike == "dislike" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? ", PostNumID)
		if err != nil {
			fmt.Println("Failed to minus dislike", err)
			ErrorPage(err, InternalError, w, r)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
