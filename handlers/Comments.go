package handlers

import (
	"encoding/json"
	d "forum/database"
	m "forum/models"
	"io"
	"net/http"
	"strings"
)

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
