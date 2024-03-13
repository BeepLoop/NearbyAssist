package complaint

import (
	complaint_query "nearbyassist/internal/db/query/complaint"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CreateComplaint(c echo.Context) error {
	complaint := new(types.Complaint)
	err := c.Bind(complaint)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	err = c.Validate(complaint)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid fields",
		})
	}

	err = complaint_query.CreateComplaint(*complaint)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]string{
		"message": "complaint created successfully",
	})
}
