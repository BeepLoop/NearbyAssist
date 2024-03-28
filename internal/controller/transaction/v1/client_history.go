package transaction

import (
	transaction_query "nearbyassist/internal/db/query/transaction"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetClientTransactionHistory(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "missing user ID")
	}

	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	history, err := transaction_query.GetClientHistory(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, history)
}
