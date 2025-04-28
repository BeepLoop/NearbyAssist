package expertise

import (
	"nearbyassist/internal/models"
	expertise_service "nearbyassist/internal/service/expertise"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type handler struct {
	expertiseService *expertise_service.Service
}

func NewHandler(expertiseService *expertise_service.Service) *handler {
	return &handler{
		expertiseService: expertiseService,
	}
}

func (h *handler) AddUserExpertise(c echo.Context) error {
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
	if err := h.expertiseService.AddUserExpertise(bearerToken, expertiseId, files[0]); err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Adding expertise not allowed",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error adding expertise",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
