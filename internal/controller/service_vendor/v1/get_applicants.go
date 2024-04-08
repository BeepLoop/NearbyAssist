package service_vendor

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetApplicants(c echo.Context) error {
	filter := c.QueryParam("filter")

    model := models.NewApplicationModel()

    applications, err := model.FindAll(filter)
    if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }

	return c.JSON(http.StatusOK, applications)
}
