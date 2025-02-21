package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	m "forum/models"
	u "forum/utils"

	d "forum/database"
)



func PostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	rows, err := d.Db.Query("SELECT title,content,created_at,post_id,filename,filepath FROM posts")
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return

	}

	defer rows.Close()

	var posts []m.Post
	for rows.Next() {

		var eachPost m.Post

		err := rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, &eachPost.Post_id, &eachPost.Filename, &eachPost.Filepath)
		eachPost.Seperate_Categories()
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		commentsCount := 0
		err = d.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", eachPost.Post_id).Scan(&commentsCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		defer rows.Close()

		eachPost.CommentsCount = commentsCount

		rows, err := d.Db.Query(`SELECT content FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		var comments []m.Comment
		for rows.Next() {
			var comment m.Comment
			rows.Scan(&comment.Content)
			comments = append(comments, comment)
		}

		eachPost.Comments = comments

		var likeCount, dislikeCount int

		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", &eachPost.Post_id).Scan(&likeCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", &eachPost.Post_id).Scan(&dislikeCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

// ==== This function will handle post creation and insertion of the post into the database ====
func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	// Parse form with multipart support. This is needed when an image is provided.
	if err := r.ParseMultipartForm(u.MaxUploadSize); err != nil {
		// If parsing fails, attempt a standard form parse which may work if no file is present.
		fmt.Println("ParseMultipartForm error:", err)
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
	}

	fmt.Println("this is the form--> ", r.Form)
	category := u.CombineCategory(r.Form["category"])
	content := strings.TrimSpace(html.EscapeString(r.FormValue("content")))
	title := strings.TrimSpace(html.EscapeString(r.FormValue("title")))

	if content == "" || title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content and Title cannot be empty"))
		return
	}

	var img m.Image
	// Attempt to retrieve the file. If no image is uploaded, proceed without processing.
	file, handler, err := r.FormFile("image")
	if err != nil {
		// Check if error is due to missing file.
		if err == http.ErrMissingFile {
			fmt.Println("No image uploaded, continuing without image")
			// Leave img fields as empty
			img.Filename = ""
			img.Path = ""
		} else {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
	} else {

		defer file.Close()

		// Validate file size
		if handler.Size > u.MaxUploadSize {
			http.Error(w, "File too large", http.StatusBadRequest)
			return
		}

		// Validate file type
		if !u.ValidateFileType(file) {
			http.Error(w, "Invalid file type", http.StatusBadRequest)
			return
		}

		// Generate unique filename
		fileName, err := u.GenerateFileName()
		if err != nil {
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}

		// Add original file extension
		fileName = fileName + filepath.Ext(handler.Filename)

		fmt.Println("filename: ", fileName)
		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll(u.UploadsDir, 0o755); err != nil {
			http.Error(w, "Error processing file", http.StatusInternalServerError)
			return
		}

		// Create new file
		filePath := filepath.Join(u.UploadsDir, fileName)
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

		modifyFilename := strings.Fields((handler.Filename))

		// Save to database
		img.Path = strings.Join(modifyFilename, "_")
		img.Path = filePath
	}
	_, err = d.Db.Exec("INSERT INTO posts (category, content, title, user_uuid ,filename,filepath) VALUES ($1, $2, $3, $4, $5, $6)", category, content, title, Profile.Uuid, img.Filename, img.Path)
	if err != nil {
		os.Remove(img.Path)
		fmt.Println("could not insert posts", err)
		// http.Error(w, "could not insert post", http.StatusInternalServerError)
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return

	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// ==== This function will handle liking a post ====
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}

	err = json.Unmarshal(str, &postID)
	if err != nil {
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return
	}

	var likeDislike string
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)

	if err == sql.ErrNoRows {

		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {

				err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,post_id,user_uuid) VALUES ('like',?,?)", postID.Post_id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				} else {
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				}

			} else {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {

		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "like" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// ==== This function will handle disliking a post ====
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	str, _ := io.ReadAll(r.Body)

	var postID struct {
		Post_id string `json:"post_id"`
	}

	err = json.Unmarshal(str, &postID)
	if err != nil {
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

					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,post_id,user_uuid) VALUES ('dislike',?,?)", postID.Post_id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				} else {
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				}

			} else {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "dislike" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// ==== This function will handle the filtration of specific user post ====
func MyPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	rows, err := d.Db.Query("SELECT title,content,post_id,created_at,filename,filepath FROM posts WHERE user_uuid = ? ", Profile.Uuid)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	defer rows.Close()

	var posts []m.Post
	for rows.Next() {
		var eachPost m.Post
		err = rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Post_id, &eachPost.CreatedAt, &eachPost.Filename, &eachPost.Filepath)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		eachPost.Seperate_Categories()
		rows, err := d.Db.Query(`SELECT content FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		var comments []m.Comment
		for rows.Next() {
			var comment m.Comment
			rows.Scan(&comment.Content)
			comments = append(comments, comment)
		}

		eachPost.Comments = comments
		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE  like_dislike = 'like' AND post_id = ?", eachPost.Post_id).Scan(&likeCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ?", eachPost.Post_id).Scan(&dislikeCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		commentsCount := 0
		err = d.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", eachPost.Post_id).Scan(&commentsCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		defer rows.Close()

		eachPost.CommentsCount = commentsCount
		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}

	// fmt.Println(Profile.Uuid)
	postsJson, err := json.Marshal(posts)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

// ==== This function will handle filtration of the posts based on the ones that have been liked ====
func FavoritesPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	likedRows, err := d.Db.Query("SELECT post_id FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'like'", Profile.Uuid)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var posts []m.Post
	for likedRows.Next() {
		var postID int
		err = likedRows.Scan(&postID)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		Postrows, err := d.Db.Query("SELECT title, content, post_id, filename,filepath FROM posts WHERE post_id = ?", postID)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		defer Postrows.Close()

		var eachPost m.Post

		for Postrows.Next() {

			err = Postrows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Post_id, &eachPost.Filename, &eachPost.Filepath)
			eachPost.Seperate_Categories()
			if err != nil {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}

		}
		commentsCount := 0
		err = d.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", eachPost.Post_id).Scan(&commentsCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		eachPost.CommentsCount = commentsCount

		rows, err := d.Db.Query(`SELECT content FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		var comments []m.Comment
		for rows.Next() {
			var comment m.Comment
			rows.Scan(&comment.Content)
			comments = append(comments, comment)
		}

		eachPost.Comments = comments

		var likeCount, dislikeCount int
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'like'", eachPost.Post_id).Scan(&likeCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", eachPost.Post_id).Scan(&dislikeCount)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}

// ==== This function will handle adding a comment to a post ====
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	r.ParseForm()
	comment := strings.TrimSpace(r.FormValue("add-comment"))
	post_id := r.FormValue("post_id")

	_, err = d.Db.Exec("INSERT INTO comments (user_uuid,post_id,content) VALUES (?,?,?)", Profile.Uuid, post_id, comment)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

// ==== This function will handle comments when they are toggled to be displayed ====
func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	str, _ := io.ReadAll(r.Body)
	if strings.Contains(string(str), "post_id") {

		var postID struct {
			Post_id string `json:"post_id"`
		}
		err := json.Unmarshal(str, &postID)
		if err != nil {
			ErrorPage(err, m.ErrorsData.BadRequest, w, r)
			return
		}

		rows, err := d.Db.Query(`SELECT comment_id, created_at, content FROM comments WHERE post_id=?`, postID.Post_id)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		defer rows.Close()
		var comments []m.Comment
		for rows.Next() {

			var (
				eachComment  m.Comment
				likeCount    int
				dislikeCount int
			)

			rows.Scan(&eachComment.Comment_id, &eachComment.CreatedAt, &eachComment.Content)
			eachComment.Post_id = postID.Post_id

			err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE  like_dislike = 'like' AND comment_id = ?", eachComment.Comment_id).Scan(&likeCount)
			if err != nil {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
			err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE like_dislike = 'dislike' AND comment_id = ?", eachComment.Comment_id).Scan(&dislikeCount)
			if err != nil {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}

			eachComment.Likes = likeCount
			eachComment.Dislikes = dislikeCount

			comments = append(comments, eachComment)
		}
		commentsJson, err := json.Marshal(comments)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(commentsJson)
	} else {

		var commentID struct {
			Comment_Id string `json:"comment_id"`
		}

		err := json.Unmarshal(str, &commentID)
		if err != nil {
			ErrorPage(err, m.ErrorsData.BadRequest, w, r)
			return
		}

		postID := ""
		err = d.Db.QueryRow(`SELECT post_id FROM  comments WHERE comment_id=?`, commentID.Comment_Id).Scan(&postID)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		rows, err := d.Db.Query(`SELECT comment_id,created_at,content FROM comments WHERE post_id=?`, postID)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		defer rows.Close()
		var comments []m.Comment
		for rows.Next() {
			var eachComment m.Comment
			rows.Scan(&eachComment.Comment_id, &eachComment.CreatedAt, &eachComment.Content)
			eachComment.Post_id = postID

			var likeCount, dislikeCount int
			err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE  like_dislike = 'like' AND comment_id = ?", eachComment.Comment_id).Scan(&likeCount)
			if err != nil {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
			err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE like_dislike = 'dislike' AND comment_id = ?", eachComment.Comment_id).Scan(&dislikeCount)
			if err != nil {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}

			eachComment.Likes = likeCount
			eachComment.Dislikes = dislikeCount
			comments = append(comments, eachComment)
		}

		commentsJson, err := json.Marshal(comments)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(commentsJson)
	}
}

// ==== This function will handle liking a comment ====
func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	str, _ := io.ReadAll(r.Body)
	var commentId struct {
		Comment_Id string `json:"comment_id"`
	}

	err = json.Unmarshal(str, &commentId)
	if err != nil {
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return
	}

	// Check if the user has already liked or disliked the post
	var likeDislike string
	// check if liked
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid).Scan(&likeDislike)

	if err == sql.ErrNoRows {

		// if not liked check if disliked
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				// check if the post exists
				err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,comment_id,user_uuid) VALUES ('like',?,?)", commentId.Comment_Id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				} else {
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				}

			} else {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "like" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// ==== This function will handle disliking a poscomment ====
func DislikeCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile,err := u.GetUserDetails(w, r)
	if err != nil {
	ErrorPage(err, m.ErrorsData.InternalError, w, r)
	return
	}

	str, _ := io.ReadAll(r.Body)
	var commentId struct {
		Comment_Id string `json:"comment_id"`
	}

	err = json.Unmarshal(str, &commentId)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}

	// Check if the user has already liked or disliked the post
	var likeDislike string
	// check if liked
	err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'dislike' AND comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid).Scan(&likeDislike)

	if err == sql.ErrNoRows {
		// if not liked check if disliked
		err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE like_dislike = 'like' AND comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid).Scan(&likeDislike)
		if err != nil {
			if err == sql.ErrNoRows {
				// check if the post exists
				err = d.Db.QueryRow("SELECT like_dislike FROM likes_dislikes WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid).Scan(&likeDislike)
				if err == sql.ErrNoRows {

					_, err = d.Db.Exec("INSERT INTO  likes_dislikes (like_dislike,comment_id,user_uuid) VALUES ('dislike',?,?)", commentId.Comment_Id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				} else {
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
					if err != nil {
						ErrorPage(err, m.ErrorsData.InternalError, w, r)
						return
					}
				}

			} else {

				ErrorPage(err, m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {

		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "dislike" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	// Get the image file path from the request URL
	imagePath := "web/uploads/" + r.URL.Path[len("/image/web/uploads/"):]

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "Failed to open image file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get the file info
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	// Set the content type and disposition
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename="+fileInfo.Name())

	// Write the image file to the response
	http.ServeFile(w, r, imagePath)
}
