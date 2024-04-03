package service_vendor

import (
	vendor_query "nearbyassist/internal/db/query/service_vendor"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func VendorApplication(c echo.Context) error {
	applicantData := new(types.VendorApplication)
	err := c.Bind(applicantData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "unable to process data provided")
	}

	err = c.Validate(applicantData)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	applicationId, err := vendor_query.VendorApplication(*applicantData)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"applicationId": applicationId,
	})
}
