package Error

import (
	"log"
	"os"
)

var ErrorLogger *log.Logger

func init() {

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	ErrorLogger = log.New(file, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// ===== This function will take an error err and display ot in the app.log file found in the home directory ====
func LogError(err error) {
	if err != nil {
		ErrorLogger.Println(err)
	}
}
