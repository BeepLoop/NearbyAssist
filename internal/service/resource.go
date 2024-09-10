package service

import (
	"mime"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/auth"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type ResourceService struct {
	encryptor auth.Encryption
}

func NewResourceService(encryptor auth.Encryption) *ResourceService {
	return &ResourceService{
		encryptor: encryptor,
	}
}

func (s *ResourceService) GetFile(c echo.Context) error {
	path := c.Param("path")

	wd, err := os.Getwd()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting working directory",
			Error:   err.Error(),
		})
	}

	bytes, err := os.ReadFile(filepath.Join(wd, path))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error reading file",
			Error:   err.Error(),
		})
	}

	decrypted, err := s.encryptor.DecryptFile(bytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: auth.DECRYPTION_ERR,
			Error:   err.Error(),
		})
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, decrypted)
}
