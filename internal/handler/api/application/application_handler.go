package application

import (
	"nearbyassist/internal/models"
	application_service "nearbyassist/internal/service/application"
	user_service "nearbyassist/internal/service/user"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type applicationHandler struct {
	applicationService *application_service.Service
	userService        *user_service.Service
}

func NewHandler(applicationService *application_service.Service, userService *user_service.Service) *applicationHandler {
	return &applicationHandler{applicationService: applicationService, userService: userService}
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
		if strings.Contains(err.Error(), "Duplicate entry") {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Application already exists",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating application",
			Error:   err.Error(),
		})
	}

	user, err := h.userService.GetUser(bearerToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving user information",
			Error:   err.Error(),
		})
	}

	// TODO: Handle notifying the user of the application status
	_ = user

	return c.JSON(http.StatusCreated, utils.Mapper{
		"application": applicationId,
	})
}
