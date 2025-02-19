package handlers

import (
	"io"
	"fmt"
	"html"
	"strings"
	"net/http"
	"database/sql"
	"encoding/json"

	d "forum/database"
	m "forum/models"
)

// ==== This function handler function will handling the posts that come in and those that leave to be posted  on the homepage and the landing page ====
func PostsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	rows, err := d.Db.Query("SELECT title,content,created_at,post_id FROM posts")
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return

	}

	defer rows.Close()

	var posts []m.Post
	for rows.Next() {
		
		var (
			eachPost m.Post
		)

		err := rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.CreatedAt, &eachPost.Post_id)
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

		rows, err := d.Db.Query(`SELECT content,likes,dislikes FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
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

	Profile := GetUserDetails(w, r)

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	r.ParseForm()
	category := CombineCategory(r.Form["category"])
	content := strings.TrimSpace(html.EscapeString(r.FormValue("content")))
	title := strings.TrimSpace(html.EscapeString(r.FormValue("title")))

	if content == "" || title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content and Title cannot be empty"))
		return
	}

	_, err := d.Db.Exec("INSERT INTO posts (category, content, title, user_uuid) VALUES ($1, $2, $3 ,$4)", category, content, title, Profile.Uuid)
	if err != nil {
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

	Profile := GetUserDetails(w, r)

	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}

	err := json.Unmarshal(str, &postID)
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
	Profile := GetUserDetails(w, r)
	
	str, _ := io.ReadAll(r.Body)
	
	var postID struct {
		Post_id string `json:"post_id"`
	}
	
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
        ErrorPage(nil, m.ErrorsData.BadRequest, w, r)
        return
    }

	Profile := GetUserDetails(w, r)
	
	rows, err := d.Db.Query("SELECT title,content,post_id FROM posts WHERE user_uuid = ? ", Profile.Uuid)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
        return
	}

	var posts []m.Post
	for rows.Next() {
		var eachPost m.Post
		err = rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Post_id)
		eachPost.Seperate_Categories()
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
        	return
		}
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

// ==== This function will handle filtration of the posts based on the ones that have been liked ====
func FavoritesPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
        ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
        return
    }

	Profile := GetUserDetails(w, r)
	
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
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
        	return
		}

		Postrows, err := d.Db.Query("SELECT likes, dislikes, title, content, post_id FROM posts WHERE post_id = ?", postID)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
        	return
		}
		var eachPost m.Post

		for Postrows.Next() {

			err = Postrows.Scan(&eachPost.Likes, &eachPost.Dislikes, &eachPost.Title, &eachPost.Content, &eachPost.Post_id)
			eachPost.Seperate_Categories()
			if err != nil {
				ErrorPage(err, m.ErrorsData.InternalError, w, r)
        		return
			}

		}

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
	
	Profile := GetUserDetails(w, r)
	
	r.ParseForm()
	comment := r.FormValue("add-comment")
	post_id := r.FormValue("post_id")

	_, err := d.Db.Exec("INSERT INTO comments (user_uuid,post_id,content) VALUES (?,?,?)", Profile.Uuid, post_id, comment)
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
		fmt.Println(string(str))
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
				eachComment m.Comment
				likeCount int
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

	Profile := GetUserDetails(w, r)

	str, _ := io.ReadAll(r.Body)
	var commentId struct {
		Comment_Id string `json:"comment_id"`
	}

	err := json.Unmarshal(str, &commentId)
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

	Profile := GetUserDetails(w, r)

	str, _ := io.ReadAll(r.Body)
	var commentId struct {
		Comment_Id string `json:"comment_id"`
	}

	err := json.Unmarshal(str, &commentId)
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
					fmt.Println("had not disliked it")
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
