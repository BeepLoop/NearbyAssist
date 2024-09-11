package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/store/application"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ApplicationService struct {
	store application.ApplicationStore
}

func NewApplicationService(store application.ApplicationStore) *ApplicationService {
	return &ApplicationService{
		store: store,
	}
}

func (s *ApplicationService) Create(c echo.Context) error {
	req := new(request.NewApplicationPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	application := new(models.ApplicationModel)
	application.ApplicantId = req.ApplicantId
	application.Job = req.Job

	applicationId, err := s.store.Create(application)
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
