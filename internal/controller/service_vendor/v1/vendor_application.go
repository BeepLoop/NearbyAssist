package service_vendor

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func VendorApplication(c echo.Context) error {
	model := models.NewApplicationModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "unable to process data provided")
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	applicationId, err := model.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"applicationId": applicationId,
	})
}
