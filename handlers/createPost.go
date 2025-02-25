package handlers

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	d "forum/database"
	m "forum/models"
	e "forum/Error"
	u "forum/utils"
)

// ===== The handler will, through method POST, add a post to the database based on  a specific user ID and add the post to the database =====
func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	if err := r.ParseMultipartForm(u.MaxUploadSize); err != nil {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
	}

	var category string

	if !u.ValidateCategory(r.Form["category"]) {
		w.WriteHeader(http.StatusBadRequest)
		ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.BadRequest, w, r)
		return
	} else {
		category = u.CombineCategory(r.Form["category"])
	}
	content := strings.TrimSpace(html.EscapeString(r.FormValue("content")))
	title := strings.TrimSpace(html.EscapeString(r.FormValue("title")))

	if content == "" || title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content and Title cannot be empty"))
		ErrorPage(fmt.Errorf("|create post handler|--> Could not create an empty post"), m.ErrorsData.BadRequest, w, r)
		return
	}

	var img m.Image
	file, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			img.Filename = ""
			img.Path = ""
		} else {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
	} else {

		defer file.Close()

		if handler.Size > u.MaxUploadSize {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		if !u.ValidateFileType(file) {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler - validate file type|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		fileName, err := u.GenerateFileName()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		fileName = fileName + filepath.Ext(handler.Filename)

		if err := os.MkdirAll(u.UploadsDir, 0o755); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		filePath := filepath.Join(u.UploadsDir, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		modifyFilename := strings.Fields((handler.Filename))

		img.Path = strings.Join(modifyFilename, "_")
		img.Path = filePath
	}
	_, err = d.Db.Exec("INSERT INTO posts (category, content, title, user_uuid ,filename,filepath) VALUES ($1, $2, $3, $4, $5, $6)", category, content, title, Profile.Uuid, img.Filename, img.Path)
	if err != nil {
		os.Remove(img.Path)
		w.WriteHeader(http.StatusBadRequest)
		ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.BadRequest, w, r)
		return

	}

	e.LOGGER(fmt.Sprintf("[SUCCESS]: User %s created a new post", Profile.Username), nil)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func ValidateCategory(str []string) bool {
	categories := map[string]bool{
		"All Categories": true,
		"Technology":     true,
		"Health":         true,
		"Math":           true,
		"Games":          true,
		"Science":        true,
		"Religion":       true,
		"Education":      true,
		"Politics":       true,
		"Fashion":        true,
		"Lifestyle":      true,
		"Sports":         true,
	}

	for _, s := range str {
		if !categories[s] {
			return false
		}
	}
	return true
}

