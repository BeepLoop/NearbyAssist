package handlers

import (
	"nearbyassist/internal/db/models"
	"nearbyassist/internal/server"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type transactionHandler struct {
	server *server.Server
}

func NewTransactionHandler(server *server.Server) *transactionHandler {
	return &transactionHandler{
		server: server,
	}
}

func (h *transactionHandler) HandleNewTransaction(c echo.Context) error {
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

func (h *transactionHandler) HandleCount(c echo.Context) error {
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

func (h *transactionHandler) HandleClientOngoingTransaction(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	model := models.NewTransactionModel()

	transactions, err := model.GetClientOngoing(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *transactionHandler) HandleVendorOngoingTransaction(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	model := models.NewTransactionModel()

	transactions, err := model.GetVendorOngoing(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *transactionHandler) HandleClientHistory(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	model := models.NewTransactionModel()

	history, err := model.VendorHistory(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, history)
}

func (h *transactionHandler) HandleVendorHistory(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	model := models.NewTransactionModel()

	history, err := model.VendorHistory(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, history)
}
