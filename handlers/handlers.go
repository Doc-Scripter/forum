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
	CreatedAt time.Time `json:"created_at"`
	Category  string    `json:"category"`
	Likes     int       `json:"likes"`
	Title     string    `json:"title"`
	Dislikes  int       `json:"dislikes"`
	Comments  []Comment    `json:"comments"`
	Content   string    `json:"content"`
	Post_ID   int       `json:"post_id"`
}
type Comment struct {
	CommentID int       `json:"comment_id"`
	PostID    int       `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
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
		fmt.Printf("Profile Section: %v\n", err)
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
	rows, err := d.Db.Query("SELECT category,title,content,created_at,post_id FROM posts")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not get posts", http.StatusInternalServerError)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var eachPost Post
		// var comments sql.NullString
		err := rows.Scan(&eachPost.Category, &eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, &eachPost.Post_ID)
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
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", eachPost.Post_ID).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get like count", http.StatusInternalServerError)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", eachPost.Post_ID).Scan(&dislikeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get dislike count", http.StatusInternalServerError)
			return
		}
		// Get comments for a specific post
		commentRows, err := d.Db.Query(`
		SELECT comment_id, content, created_at
		FROM comments
		WHERE post_id = ?
		ORDER BY created_at DESC`,
		eachPost.Post_ID,
	)
	if	err != nil {
		fmt.Println(err)
		http.Error(w, "could not get comments", http.StatusInternalServerError)
		return
	}
	defer commentRows.Close()
	var comments []Comment
	for commentRows.Next() {
	var comment Comment
	err := commentRows.Scan(&comment.CommentID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not scan comment", http.StatusInternalServerError)
		return
	}
	comments = append(comments, comment)
}
        eachPost.Comments = comments
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

	_, err := d.Db.Exec("INSERT INTO posts (category, content, title) VALUES ($1, $2, $3)", category, content, title)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "could not insert post", http.StatusInternalServerError)
		return
	}
	// Update the post like count
	// _, err = Db.Exec("UPDATE posts SET post_id = post_id + 1 ")
	// if err != nil {
	// 	fmt.Println("Failed to update post count: ",err)
	// 	http.Error(w, "Failed to update post count", http.StatusInternalServerError)
	// 	return
	// }

	http.Redirect(w, r, "/home", http.StatusSeeOther)
	fmt.Println("post created")
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}
	fmt.Println(string(str))
	err := json.Unmarshal(str, &postID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not unmarshal post id", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND post_id = ?", PostNumID).Scan(&likeDislike)
	if err == sql.ErrNoRows {
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id ", PostNumID).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				fmt.Println("had not liked it")
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
						http.Error(w, "Failed to like post", http.StatusInternalServerError)
						return
					}
				}

			} else {

				fmt.Println("Failed to query post", err)
				http.Error(w, "Failed to like post", http.StatusInternalServerError)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? ", PostNumID)
		if err != nil {
			fmt.Println("Failed to like post", err)
			http.Error(w, "Failed to like post", http.StatusInternalServerError)
			return
		}

	} else if err != nil {
		fmt.Println("Failed to check if user has liked post", err)
		http.Error(w, "Failed to check if user has liked post", http.StatusInternalServerError)
		return
	} else if likeDislike == "like" {
		fmt.Println("had liked it")
		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? ", PostNumID)
		if err != nil {
			http.Error(w, "Failed to minus like", http.StatusInternalServerError)
			fmt.Println("Failed to minus like", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}
	fmt.Println(string(str))
	err := json.Unmarshal(str, &postID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not unmarshal post id", http.StatusBadRequest)
		return
	}
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
						http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
						return
					}
				} else {
					// If the user hasn't liked or disliked the post, insert a new like
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? ", PostNumID)
					if err != nil {
						fmt.Println("Failed to dislike post", err)
						http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
						return
					}
				}

			} else {

				fmt.Println("Failed to query post", err)
				http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
				return
			}
		}

	} else if err != nil {
		fmt.Println("Failed to check if user has disliked post", err)
		http.Error(w, "Failed to check if user has disliked post", http.StatusInternalServerError)
		return
	} else if likeDislike == "dislike" {
		fmt.Println("had disliked it")
		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? ", PostNumID)
		if err != nil {
			http.Error(w, "Failed to minus dislike", http.StatusInternalServerError)
			fmt.Println("Failed to minus dislike", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

// Inserts a new comment into the database
func AddComment(db *sql.DB, postID int, content string) (int, error) {
	query := `
	INSERT INTO comments (post_id, content)
	VALUES (?, ?)
	`
	result, err := db.Exec(query, postID, content)
	if err != nil {
		return 0, fmt.Errorf("failed to insert comment: %v", err)
	}
	commentID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}
	return int(commentID), nil
}

// Fetch comments for a specific post
func GetCommentsByPostID(db *sql.DB, postID int) ([]Comment, error) {
	query := `
	SELECT comment_id, post_id, content, created_at
	FROM comments
	WHERE post_id=?
	ORDER BY created_at ASC
	`
	rows, err := db.Query(query, postID)
	if err != nil {
		return nil, fmt.Errorf("failed to query comments %v:", err)
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Content, &comment.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan comment: %v", err)
		}
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return comments, nil
}

func AddCommentHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var request struct {
			PostID  int    `json: "post_id"`
			Content string `json: "content"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Validate input
		if request.PostID <= 0 || request.Content == "" {
			http.Error(w, "Post ID and content are required", http.StatusBadRequest)
			return
		}
		// Add the comment to the database
		commentID, err := AddComment(db, request.PostID, request.Content)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to add comment: %v", err), http.StatusInternalServerError)
			return
		}
		response := struct {
			CommentID int `json: "comment_id"`
		}{
			CommentID: commentID,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
func GetCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		postIDStr := r.URL.Query().Get("post_id")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil || postID <= 0 {
			http.Error(w, "Invalid post ID", http.StatusBadRequest)
			return
		}
		// Retrieve comments from the database
		comments, err := GetCommentsByPostID(db, postID)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to retrieve comments: %v", err), http.StatusInternalServerError)
			return
		}
		// Return comments in the response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(comments)
	}
}
