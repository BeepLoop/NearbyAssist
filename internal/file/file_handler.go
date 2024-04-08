package filehandler

import (
	"mime/multipart"
	"nearbyassist/internal/db/models"

	"github.com/google/uuid"
)

type FileHandler struct {
	model models.FileModelInterface
}

func NewFileHandler(model models.FileModelInterface) *FileHandler {
	return &FileHandler{model: model}
}

func (f *FileHandler) SaveFile(file *multipart.FileHeader) (int, error) {
	uuid := uuid.New().String()

	filename, err := f.model.SaveToDisk(uuid, file)
	if err != nil {
		return 0, err
	}

	return f.model.SaveToDb(filename)
}
