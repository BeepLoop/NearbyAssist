package storage

import (
	"context"
	"nearbyassist/internal/config"
	"os"
	"path/filepath"
	"time"
)

type DiskStorage struct {
	ApplicationProofLocation string
	ServicePhotoLocation     string
	SystemComplaintLocation  string
	storagePermission        os.FileMode
}

func newDiskStorage(conf *config.Config) *DiskStorage {
	return &DiskStorage{
		ApplicationProofLocation: conf.ApplicationProofLocation,
		ServicePhotoLocation:     conf.ServicePhotoLocation,
		SystemComplaintLocation:  conf.SystemComplaintLocation,
		storagePermission:        0777,
	}
}

func (s *DiskStorage) Initialize() error {
	if err := os.MkdirAll(s.ApplicationProofLocation, s.storagePermission); err != nil {
		return err
	}

	if err := os.MkdirAll(s.ServicePhotoLocation, s.storagePermission); err != nil {
		return err
	}

	if err := os.MkdirAll(s.SystemComplaintLocation, s.storagePermission); err != nil {
		return err
	}

	return nil
}

func (s *DiskStorage) SaveFile(path string, file []byte) error {
	if err := os.WriteFile(path, file, 0777); err != nil {
		return err
	}

	return nil
}

func (s *DiskStorage) SaveServicePhoto(file []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storageDir := s.SystemComplaintLocation
	path := filepath.Join(storageDir, filename)

	if err := s.SaveFile(path, file); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return filename, nil
}

func (s *DiskStorage) SaveApplicationProof(file []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storageDir := s.ApplicationProofLocation
	path := filepath.Join(storageDir, filename)

	if err := s.SaveFile(path, file); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return filename, nil
}

func (s *DiskStorage) SaveSystemComplaint(file []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storageDir := s.SystemComplaintLocation
	path := filepath.Join(storageDir, filename)

	if err := s.SaveFile(path, file); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return filename, nil
}
