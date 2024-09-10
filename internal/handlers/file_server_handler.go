package handlers

import (
	"mime"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type fileServerHandler struct {
	server *server.Server
}

func NewFileServerHandler(s *server.Server) *fileServerHandler {
	return &fileServerHandler{
		server: s,
	}
}

func (h *fileServerHandler) HandleFileServer(c echo.Context) error {
	path := c.Param("path")

	wd, err := os.Getwd()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting working directory",
			Error:   err.Error(),
		})
	}

	bytes, err := os.ReadFile(wd + path)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error reading file",
			Error:   err.Error(),
		})
	}

	decrypted, err := h.server.Encrypt.DecryptFile(bytes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error decrypting file",
			Error:   err.Error(),
		})
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, decrypted)
}
