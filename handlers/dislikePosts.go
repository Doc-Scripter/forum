package handlers

import (
	"fmt"
	"database/sql"
	"encoding/json"
	d "forum/database"
	e "forum/Error"
	m "forum/models"
	u "forum/utils"
	"io"
	"net/http"
)

// ==== This function will handle disliking a post ====
func DislikePostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		ErrorPage(fmt.Errorf("|dislike post handler| --> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	str, _ := io.ReadAll(r.Body)

	var postID struct {
		Post_id string `json:"post_id"`
	}

	err = json.Unmarshal(str, &postID)
	if err != nil {
		ErrorPage(fmt.Errorf("|dislike post handler|---> {%v}", err), m.ErrorsData.BadRequest, w, r)
		return
	}

	//=========== Check if the user has already liked or disliked the post ===============
	var likeDislike string
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
						ErrorPage(fmt.Errorf("|dislike post handler| --> {%v}", err), m.ErrorsData.InternalError, w, r)
						return
					}
				} else {
					_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
					if err != nil {
						ErrorPage(fmt.Errorf("|dislike post handler| --> {%v}", err), m.ErrorsData.InternalError, w, r)
						return
					}
				}

			} else {
				ErrorPage(fmt.Errorf("|dislike post handler| --> {%v}", err), m.ErrorsData.InternalError, w, r)
				return
			}
		}
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = 'dislike' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			ErrorPage(fmt.Errorf("|dislike post handler| --> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

	} else if err != nil {
		ErrorPage(fmt.Errorf("|dislike post handler| --> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	} else if likeDislike == "dislike" {

		// If the user has already liked the post, minus the like
		_, err = d.Db.Exec("UPDATE likes_dislikes SET like_dislike = '' WHERE post_id = ? AND user_uuid = ?", postID.Post_id, Profile.Uuid)
		if err != nil {
			ErrorPage(fmt.Errorf("|dislike post handler| --> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	e.LOGGER(fmt.Sprintf("[SUCCESS]: User %s has disliked the post: post_id(%v)", Profile.Username, postID.Post_id), nil)
}
