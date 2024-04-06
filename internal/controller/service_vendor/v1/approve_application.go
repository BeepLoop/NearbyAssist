package service_vendor

import (
	"nearbyassist/internal/db/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func ApproveApplication(c echo.Context) error {
	applicationId := c.Param("applicationId")
	id, err := strconv.Atoi(applicationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "application ID must be a number")
	}

	model := models.NewApplicationModel()

	err = model.Approve(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Application approved successfully",
		"applicationId": id,
	})
}
