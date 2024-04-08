package complaint

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateComplaint(c echo.Context) error {
	model := models.NewComplaintModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	complaintId, err := model.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":     "complaint created successfully",
		"complaintId": complaintId,
	})
}
