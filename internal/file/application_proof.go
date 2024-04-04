package filehandler

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"
)

type ApplicationProof struct {
	ApplicantId int
	Timestamp   string
	File        *multipart.FileHeader
}

func NewApplicationProof(applicantId int, file *multipart.FileHeader) *ApplicationProof {
	return &ApplicationProof{
		ApplicantId: applicantId,
		Timestamp:   time.Now().Format("2006-01-02_15:04:05"),
		File:        file,
	}
}

func (a *ApplicationProof) SavePhoto() (string, error) {
	src, err := a.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(a.File.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%d_%s.%s", a.ApplicantId, a.Timestamp, mimeType)

	// create the file in the server
	dist, err := os.Create("store/application/" + filename)
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
