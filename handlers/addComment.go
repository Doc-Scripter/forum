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

// ==== This function will handle adding a comment to a post ====
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		ErrorPage(fmt.Errorf("|add comment handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	r.ParseForm()
	comment := strings.TrimSpace(r.FormValue("add-comment"))
	post_id := r.FormValue("post_id")

	_, err = d.Db.Exec("INSERT INTO comments (user_uuid,post_id,content) VALUES (?,?,?)", Profile.Uuid, post_id, comment)
	if err != nil {
		ErrorPage(fmt.Errorf("|add comment handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	e.LOGGER(fmt.Sprintf("[SUCCESS]: User %s added a comment to post %s", Profile.Username, post_id), nil)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
