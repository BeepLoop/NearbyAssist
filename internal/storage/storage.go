package storage

import (
	"nearbyassist/internal/config"
)

type Storage interface {
	Initialize() error
	SaveServicePhoto(file []byte, filename string) (string, error)
	SaveApplicationProof(file []byte, filename string) (string, error)
	SaveSystemComplaint(file []byte, filename string) (string, error)
	SaveFrontId(file []byte, filename string) (string, error)
	SaveBackId(file []byte, filename string) (string, error)
	SaveFace(file []byte, filename string) (string, error)
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
