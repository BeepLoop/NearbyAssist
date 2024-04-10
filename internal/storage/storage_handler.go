package storage

import (
	"nearbyassist/internal/config"
	"os"
)

type Storage struct {
	ApplicationProofLocation string
	ServicePhotoLocation     string
	storagePermission        os.FileMode
}

func NewStorage(conf *config.Config) *Storage {
	return &Storage{
		ApplicationProofLocation: conf.ApplicationProofLocation,
		ServicePhotoLocation:     conf.ServicePhotoLocation,
		storagePermission:        0777,
	}
}

func (s *Storage) CreateDirectories() error {
	if err := s.createApplicationProofLocation(); err != nil {
		return err
	}

	if err := s.createServicePhotoLocation(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) createApplicationProofLocation() error {
	err := os.MkdirAll(s.ApplicationProofLocation, s.storagePermission)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) createServicePhotoLocation() error {
	err := os.MkdirAll(s.ServicePhotoLocation, s.storagePermission)
	if err != nil {
		return err
	}

	return nil
}
