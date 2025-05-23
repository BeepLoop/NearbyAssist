package resource

import (
	"fmt"
	"nearbyassist/internal/models"
	resource_service "nearbyassist/internal/service/resource"
	"net/http"
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

func (h *resourceHandler) GetPublicFile(c echo.Context) error {
	path := c.Param("path")

	file, err := h.resourceService.GetRawFile(path)
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "File not found",
				Error:   "File not found",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error occurred retrieving the requested file",
			Error:   err.Error(),
		})
	}

	contentType, err := h.resourceService.GetPathContentType(path)
	if err != nil || contentType == "" {
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
		fmt.Println(err.Error())
		if strings.Contains(err.Error(), "invalid resource URL") {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
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

		if strings.Contains(err.Error(), "no such file") {
			return echo.NewHTTPError(http.StatusNotFound, models.Error{
				Message: "File not found",
				Error:   "File not found",
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error occurred retrieving the requested file",
			Error:   err.Error(),
		})
	}

	contentType, err := h.resourceService.GetPathContentType(params.Get("path"))
	if err != nil || contentType == "" {
		contentType = "application/octet-stream"
	}

	return c.Blob(http.StatusOK, contentType, file)
}
