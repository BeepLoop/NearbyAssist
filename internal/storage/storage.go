package storage

import (
	"mime/multipart"
	"nearbyassist/internal/config"
)

type Storage interface {
	Initialize() error
	SaveServicePhoto(uuid string, file *multipart.FileHeader) (string, error)
	SaveApplicationProof(uuid string, file *multipart.FileHeader) (string, error)
	SaveSystemComplaint(uuid string, file *multipart.FileHeader) (string, error)
}

func NewStorage(conf *config.Config) Storage {

	switch conf.StorageType {

	case config.STORAGE_DISK:
		return newDiskStorage(conf)

	case config.STORAGE_DUMMY:
		return newDummyStorage()

	default:
		panic("Invalid environment. Cannot initialize storage.")
	}
}
