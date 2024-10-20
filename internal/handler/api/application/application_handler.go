package application

import (
	"nearbyassist/internal/models"
	application_service "nearbyassist/internal/service/application"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type applicationHandler struct {
	applicationService *application_service.Service
}

func NewHandler(applicationService *application_service.Service) *applicationHandler {
	return &applicationHandler{applicationService: applicationService}
}

func (h *applicationHandler) CreateApplication(c echo.Context) error {
	job := c.FormValue("job")
	if job == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Missing required fields",
			Error:   "Missing required fields",
		})
	}

	files, err := utils.FormParser(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error reading files from form",
			Error:   err.Error(),
		})
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	applicationId, err := h.applicationService.CreateApplication(bearerToken, job, files)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating application",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"application": applicationId,
	})
}
