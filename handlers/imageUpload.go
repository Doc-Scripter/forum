package handlers

import (
	m "forum/models"
	"net/http"
	"fmt"
	e "forum/Error"
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
		ErrorPage(fmt.Errorf("|image handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
	defer file.Close()

	// Get the file info
	fileInfo, err := file.Stat()
	if err != nil {
		ErrorPage(fmt.Errorf("|image handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	// Set the content type and disposition
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename="+fileInfo.Name())

	// Write the image file to the response
	http.ServeFile(w, r, imagePath)
	e.LOGGER(fmt.Sprintf("[SUCCESS]: Served image: %s", imagePath), nil)
}
