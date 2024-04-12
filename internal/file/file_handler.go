package filehandler

import (
	"mime/multipart"
	"nearbyassist/internal/storage"

	"github.com/google/uuid"
)

type FileHandler struct {
	storage storage.Storage
}

func NewFileHandler(storage storage.Storage) *FileHandler {
	return &FileHandler{
		storage: storage,
	}
}

func (f *FileHandler) SaveServicePhoto(file *multipart.FileHeader) (string, error) {
	uuid := uuid.New().String()

	filename, err := f.storage.SaveServicePhoto(uuid, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (f *FileHandler) SaveApplicationProof(file *multipart.FileHeader) (string, error) {
	uuid := uuid.New().String()

	filename, err := f.storage.SaveApplicationProof(uuid, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}
