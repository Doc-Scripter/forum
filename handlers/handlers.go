package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	Comments  []Comment `json:"comments"`
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
	if req.Method != http.MethodGet {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if bl, _ := ValidateSession(req); bl {
		http.Redirect(rw, req, "/home", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/login.html")
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}
	if err := tmpl.Execute(rw, nil); err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

// serve the registration form
func Register(rw http.ResponseWriter, req *http.Request) {
	tmpl, err := template.ParseFiles("./web/templates/register.html")

	if req.Method != http.MethodGet {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error parsing template: %v", err)
		return
	}
	if err := tmpl.Execute(rw, nil); err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		log.Printf("Error executing template: %v", err)
		return
	}
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", 405)
		return
	}
	rows, err := d.Db.Query("SELECT category,title,content,created_at,post_id FROM posts")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not get posts", http.StatusInternalServerError)
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
		commentRows, err := d.Db.Query(`SELECT comment_id, content, created_at FROM comments WHERE post_id = ? ORDER BY created_at DESC`,
			eachPost.Post_ID)
		if err != nil {
			log.Printf("Error fetching comments for post %d: %v", eachPost.Post_ID, err)
			http.Error(w, "could not get comments", http.StatusInternalServerError)
			return
		}
		defer commentRows.Close()

		var comments []Comment
		for commentRows.Next() {
			var comment Comment
			err := commentRows.Scan(&comment.CommentID, &comment.Content, &comment.CreatedAt)
			if err != nil {
				log.Printf("Error scanning comment: %v", err)
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
	// Marshal posts into JSON
	postsJson, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "could not marshal posts", http.StatusInternalServerError)
		return
	}
	// Send JSON RESPONSE
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusBadRequest)
		return
	}

	r.ParseForm()
	category := r.FormValue("category")
	content := r.FormValue("content")
	title := r.FormValue("title")

	if category == "" || content == "" || title == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	_, err := d.Db.Exec("INSERT INTO posts (category, content, title) VALUES ($1, $2, $3)", category, content, title)
	fmt.Println(err)
	if err != nil {
		http.Error(w, "could not insert post", http.StatusInternalServerError)
		log.Printf("Error inserting post: %v", err)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read and parse the request body
	var requestBody struct {
		PostID string `json:"post_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	PostNumID, err := strconv.Atoi(requestBody.PostID)
	if err != nil {
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
			http.Error(w, "Failed to update like", http.StatusInternalServerError)
			fmt.Println("Failed to update like", err)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		fmt.Println(r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Read and parse the request body
	var requestBody struct {
		PostID string `json:"post_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	PostNumID, err := strconv.Atoi(requestBody.PostID)
	if err != nil {
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

// // Inserts a new comment into the database
// func AddComment(db *sql.DB, postID int, content string) (int, error) {
// 	query := `
// 	INSERT INTO comments (post_id, content)
// 	VALUES (?, ?)
// 	`
// 	result, err := db.Exec(query, postID, content)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to insert comment: %v", err)
// 	}
// 	commentID, err := result.LastInsertId()
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
// 	}
// 	return int(commentID), nil
// }

// Fetch comments for a specific post

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request struct {
		PostID  string `json:"post_id"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("Error decoding request body: %v", err)
		return
	}

	// Validate required fields
	if request.PostID == "" || request.Content == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Convert post_id to int
	postID, err := strconv.Atoi(request.PostID)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Validate content length
	if len(request.Content) < 1 {
		http.Error(w, "Content must not be empty", http.StatusBadRequest)
		return
	}

	// Begin transaction
	tx, err := d.Db.Begin()
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback() // Will rollback if not committed

	// Verify post exists
	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE post_id = ?)", postID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	// Insert comment
	result, err := tx.Exec(`
		INSERT INTO comments (post_id, content) 
		VALUES (?, ?)`,
		postID, request.Content,
	)
	if err != nil {
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	// Get the ID of the newly inserted comment
	commentID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get comment ID", http.StatusInternalServerError)
		log.Printf("Error getting last insert ID: %v", err)
		return
	}

	// Get the comment details to return
	var comment Comment
	err = tx.QueryRow(`
		SELECT comment_id, content, created_at 
		FROM comments 
		WHERE comment_id = ?`,
		commentID,
	).Scan(&comment.CommentID, &comment.Content, &comment.CreatedAt)
	if err != nil {
		http.Error(w, "Failed to fetch comment details", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed to commit changes", http.StatusInternalServerError)
		return
	}

	// Return success response with comment details
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// func GetCommentsByPostID(db *sql.DB, postID int) ([]Comment, error) {
// 	query := `
// 	SELECT comment_id, post_id, content, created_at
// 	FROM comments
// 	WHERE post_id=?
// 	ORDER BY created_at ASC
// 	`
// 	rows, err := db.Query(query, postID)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query comments %v:", err)
// 	}
// 	defer rows.Close()

// 	var comments []Comment
// 	for rows.Next() {
// 		var comment Comment
// 		err := rows.Scan(&comment.CommentID, &comment.PostID, &comment.Content, &comment.CreatedAt)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan comment: %v", err)
// 		}
// 		comments = append(comments, comment)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating over rows: %v", err)
// 	}
// 	return comments, nil
// }

// func AddCommentHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// Parse the request body
// 		var request struct {
// 			PostID  int    `json: "post_id"`
// 			Content string `json: "content"`
// 		}
// 		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
// 			http.Error(w, "Invalid request body", http.StatusBadRequest)
// 			return
// 		}
// 		// Validate input
// 		if request.PostID <= 0 || request.Content == "" {
// 			http.Error(w, "Post ID and content are required", http.StatusBadRequest)
// 			return
// 		}
// 		// Add the comment to the database
// 		commentID, err := AddComment(d.Db, request.PostID, request.Content)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("failed to add comment: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		response := struct {
// 			CommentID int `json: "comment_id"`
// 		}{
// 			CommentID: commentID,
// 		}
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(response)
// 	}
// }

// func GetCommentsHandler() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		postIDStr := r.URL.Query().Get("post_id")
// 		postID, err := strconv.Atoi(postIDStr)
// 		if err != nil || postID <= 0 {
// 			http.Error(w, "Invalid post ID", http.StatusBadRequest)
// 			return
// 		}
// 		// Retrieve comments from the database
// 		comments, err := GetCommentsByPostID(d.Db, postID)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("failed to retrieve comments: %v", err), http.StatusInternalServerError)
// 			return
// 		}
// 		// Return comments in the response
// 		w.Header().Set("Content-Type", "application/json")
// 		json.NewEncoder(w).Encode(comments)
// 	}
// }

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get post_id from query parameters
	postIDStr := r.URL.Query().Get("post_id")
	if postIDStr == "" {
		http.Error(w, "Missing post_id parameter", http.StatusBadRequest)
		return
	}

	// Convert post_id to int
	postID, err := strconv.Atoi(postIDStr)
	if err != nil || postID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// Verify post exists
	var exists bool
	err = d.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE post_id = ?)", postID).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Post not found", http.StatusNotFound)
		log.Printf("Error checking post existence: %v", err)
		return
	}

	// Query comments
	rows, err := d.Db.Query(`
        SELECT comment_id, content, created_at 
        FROM comments 
        WHERE post_id = ? 
        ORDER BY created_at DESC`,
		postID,
	)
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		log.Printf("Error querying comments: %v", err)
		return
	}
	defer rows.Close()

	// Scan comments into slice
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err := rows.Scan(
			&comment.CommentID,
			&comment.Content,
			&comment.CreatedAt,
		)
		if err != nil {
			http.Error(w, "Error scanning comments", http.StatusInternalServerError)
			log.Printf("Error scanning comment: %v", err)
			return
		}
		comments = append(comments, comment)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		http.Error(w, "Error processing comments", http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
