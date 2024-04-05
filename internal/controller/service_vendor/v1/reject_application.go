package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RejectApplication(c echo.Context) error {
	applicationId := c.Param("applicationId")
	id, err := strconv.Atoi(applicationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "application ID must be a number")
	}

	err = vendor_query.RejectApplication(id)
	if err != nil {
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, "reject application")
}
