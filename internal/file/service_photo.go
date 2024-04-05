package filehandler

import (
	"fmt"
	"io"
	"mime/multipart"
	upload_query "nearbyassist/internal/db/query/upload"
	"nearbyassist/internal/types"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ServicePhoto struct {
	VendorId  int
	ServiceId int
	Timestamp string
	File      *multipart.FileHeader
}

func NewServicePhoto(vendorId int, serviceId int, file *multipart.FileHeader) *ServicePhoto {
	return &ServicePhoto{
		VendorId:  vendorId,
		ServiceId: serviceId,
		Timestamp: time.Now().Format("2006-01-02_15:04:05"),
		File:      file,
	}
}

func (s *ServicePhoto) SavePhoto(uuid string) (string, error) {
	src, err := s.File.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(s.File.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%s.%s", uuid, mimeType)

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

func (s *ServicePhoto) Upload() (int, error) {
	uuid := uuid.New()
	filename, err := s.SavePhoto(uuid.String())
	if err != nil {
		return 0, err
	}

	fileData := types.ServicePhoto{
		VendorId:  s.VendorId,
		ServiceId: s.ServiceId,
		Url:       fmt.Sprintf("/resource/service/%s", filename),
	}
	id, err := upload_query.UploadServicePhoto(fileData)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
