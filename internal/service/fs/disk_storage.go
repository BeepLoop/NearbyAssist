package fs

import (
	"errors"
	"nearbyassist/internal/config"
	"nearbyassist/internal/service/auth"
	"os"
	"path/filepath"
)

type DiskStorage struct {
	categoryDirectory map[Category]string
	hash              auth.Hash
}

func NewDiskStorage(config *config.Config, hash auth.Hash) *DiskStorage {

	return &DiskStorage{
		hash: hash,
		categoryDirectory: map[Category]string{
			ID_BACK:               config.ID_BACK_DIR,
			ID_FRONT:              config.ID_FRONT_DIR,
			FACE:                  config.FACE_IMG_DIR,
			APPLICATION_PROOF_DIR: config.APPLICATION_PROOF_DIR,
			SERVICE_PHOTO_DIR:     config.SERVICE_PHOTO_DIR,
			SYS_COMPLAINT_DIR:     config.SYS_COMPLAINT_DIR,
		},
	}
}

func (s *DiskStorage) generateFilename(file []byte) (string, error) {
	name, err := s.hash.Generate(file)
	if err != nil {
		return "", err
	}

	ext, err := GetFiletype(file)
	if err != nil {
		return "", err
	}

	filename := name + "." + string(ext)
	return filename, nil
}

func (s *DiskStorage) SaveFile(file File) (string, error) {
	directory, exists := s.categoryDirectory[file.Category]
	if !exists {
		return "", errors.New("category directory not found")
	}

	filename, err := s.generateFilename(file.Data)
	if err != nil {
		return "", err
	}

	workDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(filepath.Join(workDir, directory), os.ModePerm); err != nil {
		return "", err
	}

	url := filepath.Join(directory, filename)
	fullpath := filepath.Join(workDir, url)

	if err := os.WriteFile(fullpath, file.Data, 0666); err != nil {
		return "", nil
	}

	return url, nil
}

func (s *DiskStorage) GetFile(path string) ([]byte, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(filepath.Join(workDir, path))
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
