package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

/*
FileSaver saves the uploaded file to the server.
returns the filename of the saved file or an error if any.
*/
func FileSaver(file *multipart.FileHeader, vendorId, serviceId int) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	timestamp := time.Now().Format("2006-01-02_15:04:05")
	mimeType := strings.Split(file.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%d_%d_%s.%s", vendorId, serviceId, timestamp, mimeType)

	// create the file in the server
	dist, err := os.Create("store/" + filename)
	if err != nil {
		return "", err
	}
	defer dist.Close()

	// copy the uploaded file to the opened file
	_, err = io.Copy(dist, src)
	if err != nil {
		return "", err
	}

	return filename, nil
}
