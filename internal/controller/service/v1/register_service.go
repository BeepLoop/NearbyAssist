package service

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterService(c echo.Context) error {
	model := models.NewServiceModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	insertId, err := model.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"serviceId": insertId,
	})
}
