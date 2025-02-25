package Error

import (
	"log"
	"fmt"
	"path/filepath"
	"os"
)

var ErrorLogger *log.Logger


//==== Add the LOGGER folder to the root directory and creation of the logging file. Declaration of our logger ===
func init() {

	logFolder := "LOGGING"
	loggingFile := "app.log"

	loggingFolderPath := filepath.Join(logFolder)
	loggingFilePath := filepath.Join(loggingFolderPath, loggingFile)

	if err := os.MkdirAll(loggingFolderPath, os.ModePerm); err != nil {
		fmt.Println("[LOGGER FILE]: Error creating logging folder:", err)
		return
	}

	logFile, err := os.Create(loggingFilePath)
	if err != nil {
		fmt.Println("[LOGGER FILE]: Error creating log file:", err)
		return
	}
	defer logFile.Close()

	file, err := os.OpenFile("LOGGING/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("[LOGGER FILE]: Failed to open log file: %v", err)
	}

	ErrorLogger = log.New(file, "LOGGER: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// ===== This function will take an error err and display it in the app.log file found in the LOGGER folder ====
func LOGGER(log_type string, err error) {
	if err != nil {
		ErrorLogger.Println(log_type, ": ", err)
	}else {
		ErrorLogger.Println(log_type)
	}
}
