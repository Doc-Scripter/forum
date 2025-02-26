package handlers

import (
	"encoding/json"
	"fmt"
	e "forum/Error"
	d "forum/database"
	m "forum/models"
	u "forum/utils"
	"net/http"
)

// ==== This function will handle the filtration of specific user post ====
func MyPostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	rows, err := d.Db.Query("SELECT title,content,post_id,created_at,filename,filepath FROM posts WHERE user_uuid = ? ", Profile.Uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
	defer rows.Close()

	var posts []m.Post
	for rows.Next() {
		var eachPost m.Post
		err = rows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Post_id, &eachPost.CreatedAt, &eachPost.Filename, &eachPost.Filepath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
		eachPost.Seperate_Categories()
		rows, err := d.Db.Query(`SELECT content FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
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
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE like_dislike = 'dislike' AND post_id = ?", eachPost.Post_id).Scan(&dislikeCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
		commentsCount := 0
		err = d.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", eachPost.Post_id).Scan(&commentsCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
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
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|my-post Handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	e.LOGGER("[SUCCESS]: Fetching personal posts was a success!", nil)
	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
}
