package handlers

import (
	"encoding/json"
	d "forum/database"
	"fmt"
	e "forum/Error"
	m "forum/models"
	u "forum/utils"
	"net/http"
)

// ==== This function will handle filtration of the posts based on the ones that have been liked by the current user ====
func FavoritesPostHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	likedRows, err := d.Db.Query("SELECT post_id FROM likes_dislikes WHERE user_uuid = ? AND like_dislike = 'like'", Profile.Uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	var posts []m.Post
	for likedRows.Next() {

		var postID int
		err = likedRows.Scan(&postID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		Postrows, err := d.Db.Query("SELECT title, content, post_id, filename,filepath FROM posts WHERE post_id = ?", postID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
		defer Postrows.Close()

		var eachPost m.Post

		for Postrows.Next() {

			err = Postrows.Scan(&eachPost.Title, &eachPost.Content, &eachPost.Post_id, &eachPost.Filename, &eachPost.Filepath)
			eachPost.Seperate_Categories()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
				return
			}

		}

		commentsCount := 0
		err = d.Db.QueryRow("SELECT COUNT(*) FROM comments WHERE post_id = ?", eachPost.Post_id).Scan(&commentsCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		eachPost.CommentsCount = commentsCount

		rows, err := d.Db.Query(`SELECT content FROM comments WHERE post_id = ?`, eachPost.Post_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
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
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		err = d.Db.QueryRow("SELECT COUNT(*) FROM likes_dislikes WHERE post_id = ? AND like_dislike = 'dislike'", eachPost.Post_id).Scan(&dislikeCount)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		eachPost.Likes = likeCount
		eachPost.Dislikes = dislikeCount

		posts = append(posts, eachPost)
	}
	postsJson, err := json.Marshal(posts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|favorite post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(postsJson)
	e.LOGGER("[SUCCESS]: Fetching favorite posts was a success!", nil)
}
