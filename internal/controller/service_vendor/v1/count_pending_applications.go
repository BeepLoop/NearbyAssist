package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountPendingApplications(c echo.Context) error {
	applications, err := vendor_query.CountPendingApplications()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"applications": applications,
	})
}
