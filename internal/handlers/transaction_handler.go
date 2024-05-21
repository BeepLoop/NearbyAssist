package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
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

func (h *transactionHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Transaction base route",
	})
}

func (h *transactionHandler) HandleGetTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	id, err := strconv.Atoi(transactionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "transaction ID must be a number")
	}

	transaction, err := h.server.DB.FindTransactionById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transaction": transaction,
	})
}

func (h *transactionHandler) HandleCount(c echo.Context) error {
	status := models.TransactionStatus(c.QueryParam("status"))

	count, err := h.server.DB.CountTransaction(status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactionCount": count,
	})
}

func (h *transactionHandler) HandleNewTransaction(c echo.Context) error {
	req := &request.NewTransaction{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	authHeader := c.Request().Header.Get("Authorization")
	if userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		req.ClientId = userId
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// TODO: Validate that the date is valid
	if err := utils.ValidateDateRange(req.Start, req.End); err != nil {
		if err.Error() == utils.DATE_PARSE_ERROR {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate that the vendor entered exists
	if _, err := h.server.DB.FindVendorById(req.VendorId); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Vendor not found")
	}

	// Validate that the service exists
	if _, err := h.server.DB.FindServiceById(req.ServiceId); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Service not found")
	}

	transactionId, err := h.server.DB.CreateTransaction(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "transaction created successfully",
		"transactionId": transactionId,
	})
}

func (h *transactionHandler) HandleOngoingTransaction(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	param := c.QueryParam("filter")
	filter := models.TransactionFilter(param)

	transactions, err := h.server.DB.FindAllOngoingTransaction(userId, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": transactions,
	})
}

func (h *transactionHandler) HandleHistory(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	param := c.QueryParam("filter")
	filter := models.TransactionFilter(param)

	history, err := h.server.DB.GetTransactionHistory(userId, filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"history": history,
	})
}

func (h *transactionHandler) HandleCompleteTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	id, err := strconv.Atoi(transactionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "transaction ID must be a number")
	}

	authHeader := c.Request().Header.Get("Authorization")
	userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if transaction, err := h.server.DB.FindTransactionById(id); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "transaction not found")
	} else {
		if transaction.ClientId != userId {
			return echo.NewHTTPError(http.StatusForbidden, "you're not the client of this transaction")
		}

		if transaction.Status == models.TRANSACTION_STATUS_DONE {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "transaction already marked as completed")
		}
	}

	if err := h.server.DB.CompleteTransaction(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "transaction marked as complete",
		"transactionId": transactionId,
	})
}
