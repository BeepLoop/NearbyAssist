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
	FrontIdLocation          string
	BackIdLocation           string
	FaceLocation             string
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

	if err := os.MkdirAll(s.FrontIdLocation, s.storagePermission); err != nil {
		return err
	}

	if err := os.MkdirAll(s.BackIdLocation, s.storagePermission); err != nil {
		return err
	}

	if err := os.MkdirAll(s.FaceLocation, s.storagePermission); err != nil {
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

	url := s.SystemComplaintLocation + "/" + filename
	return url, nil
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

	url := s.ApplicationProofLocation + "/" + filename
	return url, nil
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

	url := s.SystemComplaintLocation + "/" + filename
	return url, nil
}

func (s *DiskStorage) SaveFrontId(file []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storageDir := s.FrontIdLocation
	path := filepath.Join(storageDir, filename)

	if err := s.SaveFile(path, file); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	url := s.FrontIdLocation + "/" + filename
	return url, nil
}

func (s *DiskStorage) SaveBackId(file []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storageDir := s.BackIdLocation
	path := filepath.Join(storageDir, filename)

	if err := s.SaveFile(path, file); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	url := s.BackIdLocation + "/" + filename
	return url, nil
}

func (s *DiskStorage) SaveFace(file []byte, filename string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	storageDir := s.FaceLocation
	path := filepath.Join(storageDir, filename)

	if err := s.SaveFile(path, file); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	url := s.FaceLocation + "/" + filename
	return url, nil
}
