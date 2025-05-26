package application

import (
	"nearbyassist/internal/models"
	application_service "nearbyassist/internal/service/application"
	"nearbyassist/internal/service/sse"
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
	expertiseId := c.FormValue("expertiseId")
	if expertiseId == "" {
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

	bearerToken := utils.BearerTokenFromHeader(c)

	applicationId, err := h.applicationService.CreateApplication(bearerToken, expertiseId, files)
	if err != nil {
		if strings.Contains(err.Error(), application_service.ERR_DUPLICATE_EXPERTISE) {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Already have this expertise",
				Error:   err.Error(),
			})
		}

		if strings.Contains(err.Error(), application_service.ERR_DUPLICATE_APPLICATION) {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "You already have pending application for this expertise",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating application",
			Error:   err.Error(),
		})
	}

	sse.New().IncreaseApplication()

	return c.JSON(http.StatusCreated, utils.Mapper{
		"application": applicationId,
	})
}
