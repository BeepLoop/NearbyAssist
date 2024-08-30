package handlers

import (
	"nearbyassist/internal/encryption"
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

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
	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if transaction == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := c.Bind(transaction); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	if err := c.Validate(transaction); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate that the date is valid
	if err := utils.ValidateDateRange(transaction.Start, transaction.End); err != nil {
		if err.Error() == utils.DATE_PARSE_ERR {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if _, err := transaction.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while creating transaction")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "transaction created successfully",
		"transactionId": transaction.Id,
	})
}

func (h *transactionHandler) HandleCount(c echo.Context) error {
	param := c.QueryParam("filter")
	var filter models.TransactionStatusFilter
	switch param {
	case "ongoing":
		filter = models.TRANSACTION_STATUS_ONGOING
	case "done":
		filter = models.TRANSACTION_STATUS_DONE
	case "cancelled":
		filter = models.TRANSACTION_STATUS_CANCELLED
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid parameter found")
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if transaction == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	count, err := transaction.Count(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *transactionHandler) HandleGetMyTransactions(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	transactions, err := user.GetTransactions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while fetching transactions")
	}

	for _, transaction := range transactions {
		if _, err := transaction.DecryptVendorName(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
		}

		if _, err := transaction.DecryptClientName(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, encryption.DECRYPTION_ERR)
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": transactions,
	})
}

func (h *transactionHandler) HandleGetSpecificTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Transaction ID must be a number")
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if transaction == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := transaction.FindById(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Transaction not found")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transaction": transaction,
	})
}

func (h *transactionHandler) HandleOngoingTransaction(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	transactions, err := user.GetOngoingTransactions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": transactions,
	})
}

func (h *transactionHandler) HandleGetMyHistory(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	transactions, err := user.GetTransactionHistory()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"history": transactions,
	})
}

func (h *transactionHandler) HandleCompleteTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "transaction ID must be a number")
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if transaction == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := transaction.FindById(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Transaction not found")
	} else {
		if transaction.ClientId != id {
			return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to mark this transaction complete")
		}
	}

	if err := transaction.MarkComplete(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while marking transaction complete")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "transaction marked as complete",
		"transactionId": transactionId,
	})
}
