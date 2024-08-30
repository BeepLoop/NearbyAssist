package models

import (
	"context"
	"path/filepath"
	"time"
)

type ServicePhotoModel struct {
	Model
	UpdateableModel
	ServiceId string `json:"serviceId" db:"serviceId"`
	VendorId  string `json:"vendorId" db:"vendorId"`
	Url       string `json:"url" db:"url"`
}

func NewServicePhotoModel(servicePhotoId, vendorId, serviceId string, filename string) *ServicePhotoModel {
	fileLocation := filepath.Join("/resource/service", filename)

	return &ServicePhotoModel{
		Model:     Model{Id: servicePhotoId},
		ServiceId: serviceId,
		VendorId:  vendorId,
		Url:       fileLocation,
	}
}

func (s *ServicePhotoModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        INSERT INTO
            ServicePhoto (id, vendorId, serviceId, url)
        VALUES
            (:id, :vendorId, :serviceId, :url)
    `

	if _, err := s.Conn.NamedExecContext(ctx, query, s); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return s.Id, nil
}
