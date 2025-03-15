package resource

import (
	"mime"
	"nearbyassist/internal/models"
	resource_service "nearbyassist/internal/service/resource"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

type resourceHandler struct {
	resourceService *resource_service.Service
}

func NewHandler(service *resource_service.Service) *resourceHandler {
	return &resourceHandler{
		resourceService: service,
	}
}

func (h *resourceHandler) RequestFile(c echo.Context) error {
	path := c.Param("path")

	file, err := h.resourceService.GetFile(path)
	if err != nil {
		return c.JSON(http.StatusBadRequest, models.Error{
			Message: "Failed to retrieve file",
			Error:   err.Error(),
		})
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, file)
}
