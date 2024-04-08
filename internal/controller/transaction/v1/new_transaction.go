package transaction

import (
	"nearbyassist/internal/db/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func NewTransaction(c echo.Context) error {
	model := models.NewTransactionModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transactionId, err := model.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "transaction created successfully",
		"transactionId": transactionId,
	})
}
