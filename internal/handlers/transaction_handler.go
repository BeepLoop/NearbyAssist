package handlers

import (
	"nearbyassist/internal/models"
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

func (h *transactionHandler) HandleCount(c echo.Context) error {
	status := models.TransactionStatus(c.QueryParam("status"))

	count, err := h.server.DB.CountTransaction(status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"transactionCount": count,
	})
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

	transactionId, err := h.server.DB.CreateTransaction(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "transaction created successfully",
		"transactionId": transactionId,
	})
}

func (h *transactionHandler) HandleOngoingTransaction(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	param := c.QueryParam("filter")
	filter := models.TransactionFilter(param)

	transactions, err := h.server.DB.FindAllOngoingTransaction(id, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, transactions)
}

func (h *transactionHandler) HandleHistory(c echo.Context) error {
	userId := c.Param("userId")
	id, err := strconv.Atoi(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user ID must be a number")
	}

	param := c.QueryParam("filter")
	filter := models.TransactionFilter(param)

	history, err := h.server.DB.GetTransactionHistory(id, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, history)
}
