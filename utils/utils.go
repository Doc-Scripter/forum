package utils

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	MaxUploadSize = 10 << 20 // 10MB
	UploadsDir    = "./web/uploads"
	AllowedTypes  = "image/jpeg,image/png,image/gif"
)

func ValidateFileType(file multipart.File) bool {
	// Read the first 512 bytes to detect content type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return false
	}

	// Reset the file pointer
	file.Seek(0, 0)

	// Get content type and check if it's allowed
	contentType := http.DetectContentType(buffer)
	return strings.Contains(AllowedTypes, contentType)
}

func GenerateFileName() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
