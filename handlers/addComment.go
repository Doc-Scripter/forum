package handlers

import (
	"fmt"
	"net/http"
	"strings"

	m "forum/models"
	u "forum/utils"

	d "forum/database"
	e "forum/Error"
)

// ==== This function will receive a comment addition request to a post through method of POST. It then proceeds to the addition of the comment to the database ====
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorPage(fmt.Errorf("|add comment handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	r.ParseForm()
	comment := strings.TrimSpace(r.FormValue("add-comment"))
	post_id := r.FormValue("post_id")

	_, err = d.Db.Exec("INSERT INTO comments (user_uuid,post_id,content) VALUES (?,?,?)", Profile.Uuid, post_id, comment)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorPage(fmt.Errorf("|add comment handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	e.LOGGER(fmt.Sprintf("[SUCCESS]: User %s added a comment to post %s", Profile.Username, post_id), nil)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
