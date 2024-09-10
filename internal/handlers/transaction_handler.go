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
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if err := c.Bind(transaction); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(transaction); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	// Validate that the date is valid
	if err := utils.ValidateDateRange(transaction.Start, transaction.End); err != nil {
		if err.Error() == utils.DATE_PARSE_ERR {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Error parsing date",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid date range",
			Error:   err.Error(),
		})
	}

	if _, err := transaction.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating transaction",
			Error:   err.Error(),
		})
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
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Invalid filter",
			Error:   "Invalid filter",
		})
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if transaction == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	count, err := transaction.Count(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting transaction count",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *transactionHandler) HandleGetMyTransactions(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	transactions, err := user.GetTransactions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting transactions",
			Error:   err.Error(),
		})
	}

	for _, transaction := range transactions {
		if _, err := transaction.DecryptVendorName(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error decrypting vendor name",
				Error:   encryption.DECRYPTION_ERR,
			})
		}

		if _, err := transaction.DecryptClientName(h.server.Encrypt.DecryptString); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error decrypting client name",
				Error:   encryption.DECRYPTION_ERR,
			})
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": transactions,
	})
}

func (h *transactionHandler) HandleGetSpecificTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Transaction ID is required",
			Error:   "Transaction ID is required",
		})
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if transaction == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := transaction.FindById(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transaction": transaction,
	})
}

func (h *transactionHandler) HandleOngoingTransaction(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	transactions, err := user.GetOngoingTransactions()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting ongoing transactions",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": transactions,
	})
}

func (h *transactionHandler) HandleGetMyHistory(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	user := models.NewUserModelWithId(id, h.server.DB)
	transactions, err := user.GetTransactionHistory()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting transaction history",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"history": transactions,
	})
}

func (h *transactionHandler) HandleCompleteTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Transaction ID is required",
			Error:   "Transaction ID is required",
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	id, ok := claims["userId"].(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if transaction == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := transaction.FindById(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	} else {
		if transaction.ClientId != id {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Unauthorized access",
				Error:   "Unauthorized access",
			})
		}
	}

	if err := transaction.MarkComplete(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error marking transaction as complete",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "transaction marked as complete",
		"transactionId": transactionId,
	})
}
