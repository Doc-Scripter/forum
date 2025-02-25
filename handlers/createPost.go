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

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {

	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	if err := r.ParseMultipartForm(u.MaxUploadSize); err != nil {
		fmt.Println("ParseMultipartForm error:", err)
		if err := r.ParseForm(); err != nil {
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
	}

	var category string

	if !ValidateCategory(r.Form["category"]) {
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
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
	} else {

		defer file.Close()

		if handler.Size > u.MaxUploadSize {
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		if !u.ValidateFileType(file) {
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		fileName, err := u.GenerateFileName()
		if err != nil {
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		fileName = fileName + filepath.Ext(handler.Filename)

		fmt.Println("filename: ", fileName)
		if err := os.MkdirAll(u.UploadsDir, 0o755); err != nil {
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}

		filePath := filepath.Join(u.UploadsDir, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.InternalError, w, r)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
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
		ErrorPage(fmt.Errorf("|create post handler|--> {%v}", err), m.ErrorsData.BadRequest, w, r)
		return

	}

	e.LOGGER(fmt.Sprintf("[SUCCESS]: User %s created a new post", Profile.Username), nil)
	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func ValidateCategory(str []string) bool {
	categories := []string{"All Categories", "Technology", "Health", "Math", "Games", "Science", "Religion", "Education", "Politics", "Fashion", "Lifestyle", "Sports"}

	for i, s := range str {
		for _, v := range categories {


			if s == v && i==len(str)-1{
				fmt.Println(str)

				fmt.Println(v)
				return true
			}
		}
	}
	return false
}
