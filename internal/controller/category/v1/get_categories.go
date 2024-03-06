package category

import (
	query "nearbyassist/internal/db/query/category"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetCategories(c echo.Context) error {

	categories, err := query.GetCategories()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, categories)
}
