package transaction

import (
	transaction_query "nearbyassist/internal/db/query/transaction"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetVendorTransactionHistory(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "missing user ID",
		})
	}
	id, err := strconv.Atoi(userId)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "user ID must be a number",
		})
	}

	history, err := transaction_query.GetVendorHistory(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, history)
}
