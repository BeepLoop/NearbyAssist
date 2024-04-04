package filehandler

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type servicePhoto struct {
	VendorId  int
	ServiceId int
	Timestamp string
	File      *multipart.FileHeader
}

func NewServicePhoto(vendorId int, serviceId int, file *multipart.FileHeader) *servicePhoto {
	return &servicePhoto{
		VendorId:  vendorId,
		ServiceId: serviceId,
		Timestamp: time.Now().Format("2006-01-02_15:04:05"),
		File:      file,
	}
}

func (s *servicePhoto) SavePhoto() (string, error) {
	src, err := s.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(s.File.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%d_%d_%s.%s", s.VendorId, s.ServiceId, s.Timestamp, mimeType)

	// create the file in the server
	dist, err := os.Create("store/service/" + filename)
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
