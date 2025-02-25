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
	u "forum/utils"
)

func CreatePostsHandler(w http.ResponseWriter, r *http.Request) {
	Profile, err := u.GetUserDetails(w, r)
	if err != nil {
		ErrorPage(err, m.ErrorsData.InternalError, w, r)
		return
	}

	if r.Method != http.MethodPost {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}

	// Parse form with multipart support. This is needed when an image is provided.
	if err := r.ParseMultipartForm(u.MaxUploadSize); err != nil {
		// If parsing fails, attempt a standard form parse which may work if no file is present.
		fmt.Println("ParseMultipartForm error:", err)
		if err := r.ParseForm(); err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	}

	var category string

	if !ValidateCategory(r.Form["category"]) {
		w.WriteHeader(http.StatusBadRequest)
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return
	} else {
		category = u.CombineCategory(r.Form["category"])
	}
	content := strings.TrimSpace(html.EscapeString(r.FormValue("content")))
	title := strings.TrimSpace(html.EscapeString(r.FormValue("title")))

	if content == "" || title == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Content and Title cannot be empty"))
		return
	}

	var img m.Image
	// Attempt to retrieve the file. If no image is uploaded, proceed without processing.
	file, handler, err := r.FormFile("image")
	if err != nil {
		// Check if error is due to missing file.
		if err == http.ErrMissingFile {
			// Leave img fields as empty
			img.Filename = ""
			img.Path = ""
		} else {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
	} else {

		defer file.Close()

		// Validate file size
		if handler.Size > u.MaxUploadSize {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		// Validate file type
		if !u.ValidateFileType(file) {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		// Generate unique filename
		fileName, err := u.GenerateFileName()
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		// Add original file extension
		fileName = fileName + filepath.Ext(handler.Filename)

		fmt.Println("filename: ", fileName)
		// Create uploads directory if it doesn't exist
		if err := os.MkdirAll(u.UploadsDir, 0o755); err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		// Create new file
		filePath := filepath.Join(u.UploadsDir, fileName)
		dst, err := os.Create(filePath)
		if err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}
		defer dst.Close()

		// Copy file contents
		if _, err := io.Copy(dst, file); err != nil {
			ErrorPage(err, m.ErrorsData.InternalError, w, r)
			return
		}

		modifyFilename := strings.Fields((handler.Filename))

		// Save to database
		img.Path = strings.Join(modifyFilename, "_")
		img.Path = filePath
	}
	_, err = d.Db.Exec("INSERT INTO posts (category, content, title, user_uuid ,filename,filepath) VALUES ($1, $2, $3, $4, $5, $6)", category, content, title, Profile.Uuid, img.Filename, img.Path)
	if err != nil {
		os.Remove(img.Path)
		ErrorPage(err, m.ErrorsData.BadRequest, w, r)
		return

	}
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
