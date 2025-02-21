package handlers

import (
	"net/http"
	"strings"

	m "forum/models"
	u "forum/utils"

	d "forum/database"
)

// ==== This function will handle adding a comment to a post ====
func AddCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	Profile, err := u.GetUserDetails(w, r)
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
