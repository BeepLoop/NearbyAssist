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

	storageDir := s.SystemComplaintLocation
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

func (s *DiskStorage) SaveSystemComplaint(uuid string, file *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	mimeType := strings.Split(file.Header["Content-Type"][0], "/")[1]
	filename := fmt.Sprintf("%s.%s", uuid, mimeType)

	storageDir := s.SystemComplaintLocation
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
