package transaction

import (
	transaction_query "nearbyassist/internal/db/query/transaction"
	"nearbyassist/internal/types"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewTransaction(c echo.Context) error {
	transaction := new(types.NewTransaction)
	err := c.Bind(&transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	err = c.Validate(transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transactionId, err := transaction_query.NewTransaction(*transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactionId": transactionId,
	})
}
