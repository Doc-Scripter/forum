package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"

	d "forum/database"
	"forum/models"
)

// serve the login form
func LandingPage(rw http.ResponseWriter, req *http.Request) {
	if bl, _ := ValidateSession(req); bl {
		fmt.Println("valid user")
		http.Redirect(rw, req, "/home", http.StatusSeeOther)
		return
	} else if !bl {

		tmpl, err := template.ParseFiles("./web/templates/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(rw, nil)
	}
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

func getUserDetails(w http.ResponseWriter, r *http.Request) error {
	// var PD models.ProfileData

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Profile Section: No session cookie found")
		ErrorPage(err, InternalError, w, r)
		return  err
	}

	var userID string

	err = d.Db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println("Session not found in DB:", err)
		return err
	}

	query := `
	SELECT  username, email , uuid  FROM users WHERE id = ?`

	err = d.Db.QueryRow(query, userID).Scan(&models.Profile.Username, &models.Profile.Email, &models.Profile.Uuid)
	if err != nil {
		return err
	}
	return nil
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
		return
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
	rows, err := d.Db.Query("SELECT category,title,content,created_at,user_uuid FROM posts")
	if err != nil {
		fmt.Println(err)
		ErrorPage(err, InternalError, w, r)
		return
		
	}
	defer rows.Close()

	var posts []models.Post
	for rows.Next() {
		var eachPost models.Post
		// var comments sql.NullString
		err := rows.Scan(&eachPost.Category, &eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, &eachPost.User_uuid)
		if err != nil {
			fmt.Println(err)
			ErrorPage(err, InternalError, w, r)
			return
		}

		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'like'", models.Profile.Uuid).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			ErrorPage(err, InternalError, w, r)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'dislike'", models.Profile.Uuid).Scan(&dislikeCount)
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
	// w.Header().Set("Content-Type", "application/json")
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

	_, err := d.Db.Exec("INSERT INTO posts (category, content, title ,user_uuid) VALUES ($1, $2, $3 ,$4)", category, content, title, models.Profile.Uuid)
	fmt.Println("could not insert posts", err)
	if err != nil {
		// http.Error(w, "could not insert post", http.StatusInternalServerError)
		ErrorPage(err, MethodNotAllowed, w, r)
		return
	}
	// Update the post like count
	// _, err = Db.Exec("UPDATE posts SET post_id = post_id + 1 ")
	// if err != nil {
	// 	fmt.Println("Failed to update post count: ",err)
	// 	http.Error(w, "Failed to update post count", http.StatusInternalServerError)
	// 	return
	// }

	fmt.Println("post created")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
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

	// Check if the user has already liked or disliked the post
	var likeDislike string
	err := d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND user_uuid = ?", models.Profile.Uuid).Scan(&likeDislike)
	if err == sql.ErrNoRows {
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND user_uuid = ?", models.Profile.Uuid).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("had not liked it")
				err = d.Db.QueryRow("SELECT * FROM likes_dislikes WHERE user_uuid = ?", models.Profile.Uuid).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,user_uuid) VALUES ('like',?)", models.Profile.Uuid)
					if err != nil {
						fmt.Println("Failed to like post", err)
						http.Error(w, "Failed to like post", http.StatusInternalServerError)
						return
					}
				} else {
					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE user_uuid = ?", models.Profile.Uuid)
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
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE user_uuid = ?", models.Profile.Uuid)
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
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE user_uuid = ?", models.Profile.Uuid)
		if err != nil {
			fmt.Println("Failed to minus like", err)
			ErrorPage(err, InternalError, w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, BadRequest, w, r)
		return
	}

	// Check if the user has already liked or disliked the post
	var likeDislike string
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ?", PostNumID).Scan(&likeDislike)
	if err == sql.ErrNoRows {
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'slike' AND user_uuid = ?", models.Profile.Uuid).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("had  liked it")
				err = d.Db.QueryRow("SELECT * FROM likes_dislikes WHERE user_uuid = ?", models.Profile.Uuid).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,user_uuid) VALUES ('dislike',?)", models.Profile.Uuid)
					if err != nil {
						fmt.Println("Failed to dislike post", err)
						ErrorPage(err, InternalError, w, r)
						return
					}
				} else {
					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE user_uuid = ?", models.Profile.Uuid)
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
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE user_uuid = ?", models.Profile.Uuid)
		if err != nil {
			fmt.Println("Failed to minus dislike", err)
			ErrorPage(err, InternalError, w, r)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func MyPostHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := d.Db.Query("SELECT title,content,category FROM posts WHERE user_uuid = ?", models.Profile.Uuid)
	if err != nil {
		fmt.Println("unable to query my posts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var posts []models.Post
	for rows.Next() {
		var eachPost models.Post
		err = rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Category)
		if err != nil {
			fmt.Println(fmt.Println("unable to scan my posts", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'like'", models.Profile.Uuid).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get like count", http.StatusInternalServerError)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'dislike'", models.Profile.Uuid).Scan(&dislikeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get dislike count", http.StatusInternalServerError)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("unable to marshal my posts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(postsJson)
}

func LikedPostHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := d.Db.Query("SELECT title,content,category FROM posts WHERE user_uuid = (SELECT user_uuid FROM likes_dislikes WHERE like_dislike = 'like')", models.Profile.Uuid)
	if err != nil {
		fmt.Println("unable to query my posts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	var posts []models.Post
	for rows.Next() {
		var eachPost models.Post
		err = rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Category)
		if err != nil {
			fmt.Println(fmt.Println("unable to scan my posts", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'like'", models.Profile.Uuid).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get like count", http.StatusInternalServerError)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'dislike'", models.Profile.Uuid).Scan(&dislikeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get dislike count", http.StatusInternalServerError)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("unable to marshal my posts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write(postsJson)
}
