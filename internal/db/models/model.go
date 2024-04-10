package models

import (
	"mime/multipart"
	"nearbyassist/internal/db"
	"nearbyassist/internal/storage"
)

type FileModelInterface interface {
	SaveToDisk(uuid string, file *multipart.FileHeader) (string, error)
	SaveToDb(filename string) (int, error)
}

type ModelOperation interface {
	Create() (int, error)
	Update(id int) error
	Delete(id int) error
}

type Locatable interface {
	GetGeolocation() (*GeoSpatialModel, error)
}

type Model struct {
	Id        int    `json:"id" db:"id"`
	CreatedAt string `json:"createdAt" db:"createdAt"`
	Db        *db.DB
	Disk      *storage.Storage
}

type UpdateableModel struct {
	UpdatedAt string `json:"updatedAt" db:"updatedAt"`
}
