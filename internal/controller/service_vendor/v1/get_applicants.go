package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetApplicants(c echo.Context) error {
	filter := c.QueryParam("filter")

	applications := make([]types.Application, 0)
	switch filter {
	case "pending":
		result, err := vendor_query.GetPendingApplicants()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		applications = result

	case "approved":
		result, err := vendor_query.GetApprovedApplicants()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		applications = result

	case "rejected":
		result, err := vendor_query.GetRejectedApplicants()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		applications = result

	default:
		result, err := vendor_query.GetAllApplicants()
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		applications = result
	}

	return c.JSON(http.StatusOK, applications)
}
