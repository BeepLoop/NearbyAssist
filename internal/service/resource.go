package handler

import (
	"mime"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/service/fs"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type ResourceService struct {
	encryptor auth.Encryption
	fs        fs.FileStorage
}

func NewResourceService(encryptor auth.Encryption, fs fs.FileStorage) *ResourceService {
	return &ResourceService{
		encryptor: encryptor,
		fs:        fs,
	}
}

func (s *ResourceService) GetFile(c echo.Context) error {
	path := c.Param("path")

	// Retrieve the file
	bytes, err := s.fs.GetFile(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: "Failed to get file",
			Error:   err.Error(),
		})
	}

	// Decrypt file
	decrypted, err := s.encryptor.DecryptFile(bytes)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.Error{
			Message: "Failed to decrypt file",
			Error:   err.Error(),
		})
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, decrypted)
}
