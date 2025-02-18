package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"os"
	"path/filepath"
	"strconv"
	"text/template"

	e "forum/Error"

	d "forum/database"
	m "forum/models"
)

const (
	maxUploadSize = 10 << 20 // 10MB
	uploadsDir    = "./uploads"
	allowedTypes  = "image/jpeg,image/png,image/gif"
)

// serve the login form

func ErrorPage(Error error, ErrorData m.ErrorData, w http.ResponseWriter, r *http.Request) {
	e.LogError(Error)
	tmpl, err := template.ParseFiles("./web/templates/error.html")
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		e.LogError(err)
		return
	}
	if err = tmpl.Execute(w, ErrorData); err != nil {
		e.LogError(err)
		return
	}
}

// serve the login form
func LandingPage(w http.ResponseWriter, r *http.Request) {

	if bl, _ := ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}
	tmpl, err := template.ParseFiles("./web/templates/index.html")
	errD := m.ErrorsData.InternalError
	if err != nil {
		ErrorPage(err, errD, w, r)
		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		ErrorPage(err, errD, w, r)
		return
	}
}

// serve the Homepage

func getUserDetails(w http.ResponseWriter, r *http.Request) m.ProfileData {
	var Profile m.ProfileData

	cookie, err := r.Cookie("session_token")
	if err != nil {
		fmt.Println("Profile Section: No session cookie found:", err)
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return m.ProfileData{}
	}

	var userID string

	err = d.Db.QueryRow("SELECT user_id FROM sessions WHERE session_token = ?", cookie.Value).Scan(&userID)
	if err != nil {
		fmt.Println("Session not found in DB:", err)
		e.LogError(err)
		return m.ProfileData{}
	}

	query := `
		SELECT  username, email , uuid  FROM users WHERE id = ?`

	err = d.Db.QueryRow(query, userID).Scan(&Profile.Username, &Profile.Email, &Profile.Uuid)
	if err != nil {
		e.LogError(err)
		return m.ProfileData{}
	}
	
	// Generate and set the initials
	Profile.Initials = Profile.GenerateInitials()
	
	return Profile
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	if bl, _ := ValidateSession(r); !bl {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}

	Profile := getUserDetails(w, r)
	
	
	tmpl, err := template.ParseFiles("./web/templates/home.html")
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}

	if err = tmpl.Execute(w, Profile); err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
}

// serve the login form
func Login(w http.ResponseWriter, r *http.Request) {
	if bl, _ := ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/login.html")
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	if err = tmpl.Execute(w, nil); err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
}

// serve the registration form
func Register(w http.ResponseWriter, r *http.Request) {

	if bl, _ := ValidateSession(r); bl {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}

	tmpl, err := template.ParseFiles("./web/templates/register.html")
	if err != nil {
		log.Fatal(err)
	}
	if err = tmpl.Execute(w, nil); err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.InternalError, w, r)
		return
	}
	rows, err := d.Db.Query("SELECT category,title,content,created_at,post_id FROM posts")
	if err != nil {
		fmt.Println(err)
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return

	}
	defer rows.Close()

	var posts []m.Post
	for rows.Next() {
		var eachPost m.Post
		// var comments sql.NullString
		err := rows.Scan(&eachPost.Category, &eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, &eachPost.Post_id)
		if err != nil {
			fmt.Println(err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		commentsCount := 0
		err = d.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", eachPost.Post_id).Scan(&commentsCount)
		if err != nil {
			fmt.Println("unable ro query comments", err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		defer rows.Close()

		eachPost.CommentsCount = commentsCount

		rows, err := d.Db.Query(`SELECT content,likes,dislikes FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
			fmt.Println("unable to query comments", err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		var comments []m.Comment
		for rows.Next() {
			var comment m.Comment
			rows.Scan(&comment.Content, &comment.Likes, &comment.Dislikes)
			comments = append(comments, comment)
		}

		eachPost.Comments = comments

		var likeCount, dislikeCount int

		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", &eachPost.Post_id).Scan(&likeCount)
		if err != nil {
			fmt.Println("unable to query likes and dislikes", err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", &eachPost.Post_id).Scan(&dislikeCount)
		if err != nil {
			fmt.Println("unable to query likes and dislikes", err)

			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		fmt.Println("unable to marshal", err)

		// http.Error(w, "could not marshal posts", http.StatusInternalServerError)
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {

	Profile := getUserDetails(w, r)

	if r.Method != http.MethodPost {
		// http.Error(w, "method not allowed", http.StatusBadRequest)
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}
	fmt.Println("creating post")
	r.ParseForm()
	category := r.FormValue("category")
	content := r.FormValue("content")
	title := r.FormValue("title")

	_, err := d.Db.Exec("INSERT INTO posts (category, content, title, user_uuid) VALUES ($1, $2, $3 ,$4)", category, content, title, Profile.Uuid)
	if err != nil {
		fmt.Println("could not insert posts", err)
		// http.Error(w, "could not insert post", http.StatusInternalServerError)
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return
	}
	// Update the post like count
	// _, err = Db.Exec("UPDATE posts SET post_id = post_id + 1 ")
	// if err != nil {
	// 	fmt.Println("Failed to update post count: ",err)
	// 	http.Error(w, "Failed to update post count", http.StatusInternalServerError)
	// 	return
	// }
	fmt.Println("Post created")
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}
	Profile := getUserDetails(w, r)

	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}
	fmt.Println(string(str))
	err := json.Unmarshal(str, &postID)
	if err != nil {
		fmt.Println("could not unmarshal post id")
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return
	}

	// Check if the user has already liked or disliked the post
	var likeDislike string
	// check if liked
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)

	if err == sql.ErrNoRows {
		// if not liked check if disliked
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				// check if the post exists
				err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					fmt.Println("had not liked it")
					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,post_id,user_uuid) VALUES ('like',?,?)", postID.Post_id, Profile.Uuid)
					if err != nil {
						fmt.Println("Failed to like post", err)
						http.Error(w, "Failed to like post", http.StatusInternalServerError)
						return
					}
				} else {
					fmt.Println("had not liked it")
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
					if err != nil {
						fmt.Println("Failed to like post", err)
						http.Error(w, "Failed to like post", http.StatusInternalServerError)
						return
					}
				}

				// If the user hasn't liked or disliked the post, insert a new like
				fmt.Println("has liked it")

			} else {

				fmt.Println("Failed to query post", err)
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			fmt.Println("Failed to like post", err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {

		fmt.Println("Failed to check if user has liked post", err)
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "like" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			fmt.Println("Failed to minus like", err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}
	Profile := getUserDetails(w, r)
	
	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}
	fmt.Println(string(str))
	err := json.Unmarshal(str, &postID)
	if err != nil {
		fmt.Println("could not unmarshal post id")
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return
	}

	// Check if the user has already liked or disliked the post
	var likeDislike string
	// check if liked
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)

	if err == sql.ErrNoRows {
		// if not liked check if disliked
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					fmt.Println("had not disliked it")
					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,post_id,user_uuid) VALUES ('dislike',?,?)", postID.Post_id, Profile.Uuid)
					if err != nil {
						fmt.Println("Failed to dislike post", err)
						http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
						return
					}
				} else {
					fmt.Println("had not disliked it")
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
					if err != nil {
						fmt.Println("Failed to dislike post", err)
						http.Error(w, "Failed to dislike post", http.StatusInternalServerError)
						return
					}
				}

				// If the user hasn't liked or disliked the post, insert a new like
				fmt.Println("has disliked it")

			} else {

				fmt.Println("Failed to query post", err)
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			fmt.Println("Failed to dislike post", err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {

		fmt.Println("Failed to check if user has disliked post", err)
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "dislike" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			fmt.Println("Failed to minus dislike", err)
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func MyPostHandler(w http.ResponseWriter, r *http.Request) {
	Profile := getUserDetails(w, r)
	
	rows, err := d.Db.Query("SELECT title,content,category,post_id FROM posts WHERE user_uuid = ? ", Profile.Uuid)
	if err != nil {
		fmt.Println("unable to query my posts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var posts []m.Post
	for rows.Next() {
		var eachPost m.Post
		err = rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Category, &eachPost.Post_id)
		if err != nil {
			fmt.Println(fmt.Println("unable to scan my posts", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE  like_dislike = 'like' AND post_id = ?", eachPost.Post_id).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get like count", http.StatusInternalServerError)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ?", eachPost.Post_id).Scan(&dislikeCount)
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

func FavoritesPostHandler(w http.ResponseWriter, r *http.Request) {

	Profile := getUserDetails(w, r)
	
	likedRows, err := d.Db.Query("SELECT post_id FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'like'", Profile.Uuid)
	if err != nil {
		fmt.Println("unable to query my posts", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var posts []m.Post
	for likedRows.Next() {
		var postID int
		err = likedRows.Scan(&postID)
		if err != nil {
			fmt.Println(fmt.Println("unable to scan my posts", err))
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		Postrows, err := d.Db.Query("SELECT category, likes, dislikes, title, content, post_id FROM posts WHERE post_id = ?", postID)
		if err != nil {
			fmt.Println("unable to query my posts", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		var eachPost m.Post

		for Postrows.Next() {

			err = Postrows.Scan(&eachPost.Category, &eachPost.Likes, &eachPost.Dislikes, &eachPost.Title, &eachPost.Content, &eachPost.Post_id)
			if err != nil {
				fmt.Println(fmt.Println("unable to scan my posts", err))
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

		}

		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", eachPost.Post_id).Scan(&likeCount)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "could not get like count", http.StatusInternalServerError)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", eachPost.Post_id).Scan(&dislikeCount)
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}
	
	Profile := getUserDetails(w, r)
	
	r.ParseForm()
	comment := r.FormValue("add-comment")
	post_id := r.FormValue("post_id")

	_, err := d.Db.Exec("INSERT INTO comments (user_uuid,post_id,content) VALUES (?,?,?)", Profile.Uuid, post_id, comment)
	if err != nil {
		fmt.Println("could not insert comment", err)
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	// r.Method = http.MethodGet
	// PostsHandler(w, r)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
		return
	}

	str, _ := io.ReadAll(r.Body)
	if strings.Contains(string(str), "post_id") {

		var postID struct {
			Post_id string `json:"post_id"`
		}
		fmt.Println(string(str))
		err := json.Unmarshal(str, &postID)
		if err != nil {
			fmt.Println("could not unmarshal post id")
			ErrorPage(err, m.ErrorsData.BadRequest, w, r)
			return
		}

	rows, err := d.Db.Query(`SELECT created_at,likes,dislikes,content FROM comments WHERE post_id=?`, postID.Post_id)
	if err != nil {
		fmt.Println("could not query comments", err)
		http.Error(w, "could not get like count", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var comments []m.Comment
	for rows.Next() {
		var eachComment m.Comment
		rows.Scan(&eachComment.CreatedAt, &eachComment.Likes, &eachComment.Dislikes, &eachComment.Content)
		comments = append(comments, eachComment)
	}
	commentsJson, err := json.Marshal(comments)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "could not get like count", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(commentsJson)
}
}

// SaveImage saves image information to the database
func SaveImage(db *sql.DB, img *m.Image) error {
	query := `INSERT INTO images (user_id, post_id, filename, path) 
              VALUES (?, ?, ?, ?)`

	result, err := db.Exec(query, img.UserID, img.PostID, img.Filename, img.Path)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	img.ID = id
	return nil
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authenticated
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	isValid, userIDStr := ValidateSession(r)
	if !isValid {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Convert userID from string to int64
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusInternalServerError)
		return
	}
	// Parse multipart form
	r.ParseMultipartForm(maxUploadSize)
	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file size
	if handler.Size > maxUploadSize {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Validate file type
	if !validateFileType(file) {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
		return
	}

	// Generate unique filename
	fileName, err := generateFileName()
	if err != nil {
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	// Add original file extension
	fileName = fileName + filepath.Ext(handler.Filename)

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		http.Error(w, "Error processing file", http.StatusInternalServerError)
		return
	}

	// Create new file
	filePath := filepath.Join(uploadsDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	// Save to database
	img := &m.Image{
		UserID:   userID,
		Filename: handler.Filename,
		Path:     filePath,
	}

	if err := SaveImage(d.Db, img); err != nil {
		// Clean up file if database save fails
		os.Remove(filePath)
		http.Error(w, "Error saving file info", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"id":%d,"path":"%s"}`, img.ID, img.Path)
}
