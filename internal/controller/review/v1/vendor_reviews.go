package review

import (
	review_query "nearbyassist/internal/db/query/review"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func VendorReview(c echo.Context) error {
	vendorId := c.Param("vendorId")
	if vendorId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "vendorId required",
		})
	}
	id := strings.ReplaceAll(vendorId, "/", "")

	reviews, err := review_query.VendorReview(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, reviews)
}
