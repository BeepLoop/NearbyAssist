package complaint

import (
	complaint_query "nearbyassist/internal/db/query/complaint"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountComplaints(c echo.Context) error {
	complaints, err := complaint_query.CountComplaints()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"complaints": complaints,
	})
}
