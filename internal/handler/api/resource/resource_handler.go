package resource

import (
	"mime"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/fs"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type resourceHandler struct {
	fs fs.FileStorage
}

func NewHandler(fs fs.FileStorage) *resourceHandler {
	return &resourceHandler{
		fs: fs,
	}
}

func (h *resourceHandler) GetFile(c echo.Context) error {
	path := c.Param("path")

	// Retrieve the file
	bytes, err := h.fs.GetFile(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: "Failed to get file",
			Error:   err.Error(),
		})
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, bytes)
}
