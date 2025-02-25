package handlers

import (
	"database/sql"
	"encoding/json"
	d "forum/database"
	m "forum/models"
	e "forum/Error"
	u "forum/utils"
	"fmt"
	"io"
	"net/http"
)

// ==== This function will handle liking a comment ====
func LikeCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	str, _ := io.ReadAll(r.Body)
	var commentId struct {
		Comment_Id string `json:"comment_id"`
	}

	err = json.Unmarshal(str, &commentId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.BadRequest, w, r)
		return
	}

	//======= Check if the user has already liked or disliked the post =========
	var likeDislike string
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
						ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
						return
					}
				} else {
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
					if err != nil {
						ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
						return
					}
				}

			} else {
				ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
		if err != nil {
			ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {
		ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "like" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE comment_id = ? AND user_uuid = ?", commentId.Comment_Id, Profile.Uuid)
		if err != nil {
			ErrorPage(fmt.Errorf("|like comment handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
	}

	e.LOGGER(fmt.Sprintf("[SUCCESS]: User %s liked the comment: comment(%v)", Profile.Username, commentId.Comment_Id), nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
