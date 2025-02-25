package handlers

import (
	m "forum/models"
	"net/http"
	"fmt"
	e "forum/Error"
	"os"
)

//===== The handler will process the binary files that will be posted and be fetched from the database =====
func ImageHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		ErrorPage(nil, m.ErrorsData.MethodNotAllowed, w, r)
		return
	}
	imagePath := "web/uploads/" + r.URL.Path[len("/image/web/uploads/"):]

	file, err := os.Open(imagePath)
	if err != nil {
		ErrorPage(fmt.Errorf("|image handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		ErrorPage(fmt.Errorf("|image handler| ---> {%v}", err), m.ErrorsData.InternalError, w, r)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename="+fileInfo.Name())

	http.ServeFile(w, r, imagePath)
	e.LOGGER(fmt.Sprintf("[SUCCESS]: Served image: %s", imagePath), nil)
}
