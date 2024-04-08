package models

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"nearbyassist/internal/db"
	"os"
	"strings"
	"time"
)

type ServicePhotoModel struct {
	Model
	UpdateableModel
	ServiceId int    `json:"serviceId" db:"serviceId"`
	VendorId  int    `json:"vendorId" db:"vendorId"`
	Url       string `json:"url" db:"url"`
}

func NewServicePhotoModel(vendorId, serviceId int) *ServicePhotoModel {
	return &ServicePhotoModel{
		ServiceId: serviceId,
		VendorId:  vendorId,
	}
}

func (s *ServicePhotoModel) Create() (int, error) {
	return 0, nil
}

func (s *ServicePhotoModel) Update(id int) error {
	return nil
}

func (s *ServicePhotoModel) Delete(id int) error {
	return nil
}

func (s *ServicePhotoModel) SaveToDisk(uuid string, file *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(file.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%s.%s", uuid, mimeType)

	dist, err := os.Create("store/service/" + filename)
	if err != nil {
		return "", err
	}
	defer dist.Close()

	_, err = io.Copy(dist, src)
	if err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return filename, nil
}

func (s *ServicePhotoModel) SaveToDb(filename string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	s.Url = "/resource/service/" + filename

	query := `
        INSERT INTO
            ServicePhoto (vendorId, serviceId, url)
        VALUES 
            (:vendorId, :serviceId, :url)
    `

	res, err := db.Connection.NamedExecContext(ctx, query, s)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}
