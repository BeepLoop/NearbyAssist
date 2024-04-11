package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"nearbyassist/internal/config"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DiskStorage struct {
	ApplicationProofLocation string
	ServicePhotoLocation     string
	storagePermission        os.FileMode
}

func NewDiskStorage(conf *config.Config) *DiskStorage {
	return &DiskStorage{
		ApplicationProofLocation: conf.ApplicationProofLocation,
		ServicePhotoLocation:     conf.ServicePhotoLocation,
		storagePermission:        0777,
	}
}

func (s *DiskStorage) CreateDirectories() error {
	if err := s.createApplicationProofLocation(); err != nil {
		return err
	}

	if err := s.createServicePhotoLocation(); err != nil {
		return err
	}

	return nil
}

func (s *DiskStorage) createApplicationProofLocation() error {
	err := os.MkdirAll(s.ApplicationProofLocation, s.storagePermission)
	if err != nil {
		return err
	}

	return nil
}

func (s *DiskStorage) createServicePhotoLocation() error {
	err := os.MkdirAll(s.ServicePhotoLocation, s.storagePermission)
	if err != nil {
		return err
	}

	return nil
}

func (s *DiskStorage) SaveServicePhoto(uuid string, file *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(file.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%s.%s", uuid, mimeType)

	storageDir := s.ServicePhotoLocation
	path := filepath.Join(storageDir, filename)

	dist, err := os.Create(path)
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

func (s *DiskStorage) SaveApplicationProof(uuid string, file *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(file.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%s.%s", uuid, mimeType)

	storageDir := s.ApplicationProofLocation
	path := filepath.Join(storageDir, filename)

	dist, err := os.Create(path)
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
