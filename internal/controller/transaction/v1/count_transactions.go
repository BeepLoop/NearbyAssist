package transaction

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func CountTransactions(c echo.Context) error {
	status := c.QueryParam("status")

	model := models.NewTransactionModel()

	count, err := model.Count(status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactionCount": count,
	})
}
