package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountApplications(c echo.Context) error {
	filter := c.QueryParam("filter")

	var applicationCount int
	switch filter {
	case "approved":
		count, err := vendor_query.CountApprovedApplications()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		applicationCount = count
	case "rejected":
		count, err := vendor_query.CountRejectedApplications()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		applicationCount = count
	case "pending":
		count, err := vendor_query.CountPendingApplications()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		applicationCount = count
	default:
		count, err := vendor_query.CountAllApplications()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		applicationCount = count
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"applicationCount": applicationCount,
	})
}
