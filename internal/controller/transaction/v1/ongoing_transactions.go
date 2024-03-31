package transaction

import (
	transaction_query "nearbyassist/internal/db/query/transaction"
	"net/http"

	"github.com/labstack/echo/v4"
)

func OngoingTransactions(c echo.Context) error {
	transactions, err := transaction_query.OngoingTransactions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"ongoing": transactions,
	})
}
