package resource

import (
	"mime"
	"nearbyassist/internal/models"
	resource_service "nearbyassist/internal/service/resource"
	"net/http"
	"path/filepath"
	"strings"

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

	file, err := h.resourceService.GetRawFile(path)
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

func (h *resourceHandler) GetPrivateFile(c echo.Context) error {
	params := c.QueryParams()
	if !params.Has("path") || !params.Has("expiry") || !params.Has("signature") {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid resource URL",
			Error:   "Invalid resource URL",
		})
	}

	file, err := h.resourceService.GetPrivateFile(params.Get("path"), params.Get("signature"), params.Get("expiry"))
	if err != nil {
		if strings.Contains(err.Error(), "unvalid resource URL") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "Invalid resource URL",
				Error:   err.Error(),
			})
		}

		if strings.Contains(err.Error(), "unauthorized access") {
			return echo.NewHTTPError(http.StatusUnauthorized, models.Error{
				Message: "Unauthorized file access",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error occurred retrieving the requested file",
			Error:   err.Error(),
		})
	}

	contentType := mime.TypeByExtension(filepath.Ext(params.Get("path")))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, file)
}
