package handlers

import (
	"database/sql"
	"encoding/json"
	d "forum/database"
	m "forum/models"
	"fmt"
	e "forum/Error"
	u "forum/utils"
	"io"
	"net/http"
)


// ==== This function will handle liking a post and adding it to the database ====
func LikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	str, _ := io.ReadAll(r.Body)
	var postID struct {
		Post_id string `json:"post_id"`
	}

	err = json.Unmarshal(str, &postID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.BadRequest, w, r)
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
						w.WriteHeader(http.StatusInternalServerError)
						ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
						return
					}
				} else {
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
						return
					}
				}

			} else {
				w.WriteHeader(http.StatusInternalServerError)
				ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'like' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "like" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ErrorPage(fmt.Errorf("|like post handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
	}

	e.LOGGER(fmt.Sprintf("[SUCCESS]: User %s has liked the post: post_id(%v)", Profile.Username, postID.Post_id), nil)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
