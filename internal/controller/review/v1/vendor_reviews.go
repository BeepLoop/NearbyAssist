package review

import (
	review_query "nearbyassist/internal/db/query/review"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func VendorReview(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing vendor ID")
	}

	id, err := strconv.Atoi(vendorId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "vendor ID must be a number")
	}

	reviews, err := review_query.VendorReview(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, reviews)
}
