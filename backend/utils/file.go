package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func UploadFile(r *http.Request, keyName string, allowedFormats []string) (string, *ErrorResponse) {
	file, header, err := r.FormFile(keyName)
	if err != nil {
		return "", &ErrorResponse{Type: "fileMissing", Detail: err.Error()}
	}
	defer file.Close()

	if err := validateFileFormat(header, allowedFormats); err != nil {
		return "", &ErrorResponse{Type: "validateFormat", Detail: err.Error()}
	}

	if err := os.MkdirAll("uploads", 0755); err != nil {
		return "", &ErrorResponse{Type: "MkdirAll", Detail: err.Error()}
	}

	fileName := rand.Text() + header.Filename
	path := filepath.Join("uploads", fileName)

	dst, err := os.Create(path)
	if err != nil {
		return "", &ErrorResponse{Type: "createOs", Detail: err.Error()}
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", &ErrorResponse{Type: "ioCopy", Detail: err.Error()}
	}

	return fileName, nil
}

func validateFileFormat(header *multipart.FileHeader, allowedFormats []string) error {
	fileFormat := filepath.Ext(header.Filename)

	for _, format := range allowedFormats {
		if format == fileFormat {
			return nil
		}
	}

	return fmt.Errorf("format : '%s' is not allowed. Must be in: %s", fileFormat, strings.Join(allowedFormats, ", "))
}
