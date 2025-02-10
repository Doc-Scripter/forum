package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"
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
	Email string
}

type Post struct {
	CreatedAt time.Time `json:"created_at"`
	Category  string    `json:"category"`
	Likes     int       `json:"likes"`
	Title     string    `json:"title"`
	Dislikes  int       `json:"dislikes"`
	Comments  string   `json:"comments"`
	Content   string    `json:"content"`
	ID        int       `json:"id"`
	
}

// serve the Homepage

func getUserDetails(r *http.Request) (ProfileData, error) {
	var PD ProfileData

	cookie, err := r.Cookie("session_token")
    if err != nil {
        fmt.Println("Profile Section: No session cookie found")
        return ProfileData{}, err
    }

    var (
        userID    string
    )

    err = Db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
    if err != nil {
        fmt.Println("Session not found in DB:", err)
        return ProfileData{}, err
    }

	query := `
	SELECT  username, email FROM users WHERE id = ?`

	err = Db.QueryRow(query, userID).Scan(&PD.Username, &PD.Email)
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
	rows, err := Db.Query("SELECT category,title,content,created_at FROM posts")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not get posts", http.StatusInternalServerError)
	}
	defer rows.Close()
	

	var posts []Post
	for rows.Next() {
		var eachPost Post
		// var comments sql.NullString
		err := rows.Scan(&eachPost.Category, &eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, )
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get post", http.StatusInternalServerError)
			return
		}
		// if comments.Valid {
		// 	eachPost.Comments = &comments.String
		// } else {
		// 	eachPost.Comments = nil
		// }
		// Retrieve like and dislike counts from the database
		var likeCount, dislikeCount int
		err = Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", eachPost.ID).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get like count", http.StatusInternalServerError)
			return
		}
		err = Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", eachPost.ID).Scan(&dislikeCount)
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
	title := r.FormValue("title")

	_, err := Db.Exec("INSERT INTO posts (category, content, title) VALUES ($1, $2, $3)", category, content, title)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "could not insert post", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
	fmt.Println("post created")
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {

	// fmt.Println("liking post")
	// fmt.Println(r.Method)
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	userID := r.FormValue("user_id")

	// Check if the user has already liked or disliked the post
	var likeDislike string
	err := Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&likeDislike)
	if err == sql.ErrNoRows {
		// If the user hasn't liked or disliked the post, insert a new like
		_, err = Db.Exec("INSERT INTO likes_dislikes (post_id, user_id, like_dislike) VALUES (?, ?, 'like')", postID, userID)
		if err != nil {
			http.Error(w, "Failed to like post", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Failed to check if user has liked post", http.StatusInternalServerError)
		return
	} else if likeDislike == "like" {
		// If the user has already liked the post, do nothing
		http.Error(w, "You have already liked this post", http.StatusBadRequest)
		return
	} else if likeDislike == "dislike" {
		// If the user has disliked the post, update the like_dislike column to 'like'
		_, err = Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			http.Error(w, "Failed to update like", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println("disliking post")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postID := r.FormValue("post_id")
	userID := r.FormValue("user_id")

	// Check if the user has already liked or disliked the post
	var likeDislike string
	err := Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&likeDislike)
	if err == sql.ErrNoRows {
		// If the user hasn't liked or disliked the post, insert a new dislike
		_, err = Db.Exec("INSERT INTO likes_dislikes (post_id, user_id, like_dislike) VALUES (?, ?, 'dislike')", postID, userID)
		if err != nil {
			http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Failed to check if user has disliked post", http.StatusInternalServerError)
		return
	} else if likeDislike == "dislike" {
		// If the user has already disliked the post, do nothing
		http.Error(w, "You have already disliked this post", http.StatusBadRequest)
		return
	} else if likeDislike == "like" {
		// If the user has liked the post, update the like_dislike column to 'dislike'
		_, err = Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			http.Error(w, "Failed to update dislike", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
