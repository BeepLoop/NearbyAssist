package transaction

import (
	transaction_query "nearbyassist/internal/db/query/transaction"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func OngoingClientTransactions(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing owner ID",
		})
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "owner ID must be a number",
		})
	}

	transactions, err := transaction_query.ClientOngoingTransactions(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, transactions)
}
