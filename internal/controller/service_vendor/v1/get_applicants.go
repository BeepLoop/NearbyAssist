package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetApplicants(c echo.Context) error {
	applicants, err := vendor_query.GetApplicants()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, applicants)
}
