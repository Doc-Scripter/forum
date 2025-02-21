package handlers

import (
	m "forum/models"
	"net/http"
	"os"
)

func ImageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	// Get the image file path from the request URL
	imagePath := "web/uploads/" + r.URL.Path[len("/image/web/uploads/"):]

	// Open the image file
	file, err := os.Open(imagePath)
	if err != nil {
		http.Error(w, "Failed to open image file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Get the file info
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to get file info", http.StatusInternalServerError)
		return
	}

	// Set the content type and disposition
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename="+fileInfo.Name())

	// Write the image file to the response
	http.ServeFile(w, r, imagePath)
}
