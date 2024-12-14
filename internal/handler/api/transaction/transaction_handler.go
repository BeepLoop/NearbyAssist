package transaction

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	transaction_service "nearbyassist/internal/service/transaction"
	user_service "nearbyassist/internal/service/user"
	"nearbyassist/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type transactionHandler struct {
	transactionService *transaction_service.Service
	userService        *user_service.Service
}

func NewHandler(transactionService *transaction_service.Service, useService *user_service.Service) *transactionHandler {
	return &transactionHandler{transactionService: transactionService, userService: useService}
}

func (h *transactionHandler) CreateTransaction(c echo.Context) error {
	req := new(request.NewTransactionPayload)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	transactionId, err := h.transactionService.CreateTransaction(req)
	if err != nil {
		if strings.Contains(err.Error(), "You already have an ongoing or pending transaction for this service") {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "You already have an ongoing or pending transaction for this service",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transaction": transactionId,
	})
}

func (h *transactionHandler) GetUserTransactionList(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]
	filter := c.QueryParam("filter")

	var transactions []*models.TransactionModel
	if filter == "" || filter == "all" {
		if result, err := h.transactionService.GetUserTransactionList(bearerToken); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting user transactions",
				Error:   err.Error(),
			})
		} else {
			transactions = result
		}
	} else if filter == "sent" {
		if result, err := h.transactionService.GetTransactionUserSent(bearerToken); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting user transactions",
				Error:   err.Error(),
			})
		} else {
			transactions = result
		}
	} else if filter == "received" {
		if result, err := h.transactionService.GetTransactionUserReceived(bearerToken); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error getting user transactions",
				Error:   err.Error(),
			})
		} else {
			transactions = result
		}
	}

	response := []struct {
		Id       string              `json:"id"`
		Cost     float64             `json:"cost"`
		VendorId string              `json:"vendorId"`
		ClientId string              `json:"clientId"`
		Status   string              `json:"status"`
		Service  models.ServiceModel `json:"service"`
		Extras   []models.ExtraModel `json:"extras"`
	}{}

	for _, transaction := range transactions {
		cost, err := strconv.ParseFloat(transaction.Cost, 64)
		if err != nil {
			return err
		}

		response = append(response, struct {
			Id       string              `json:"id"`
			Cost     float64             `json:"cost"`
			VendorId string              `json:"vendorId"`
			ClientId string              `json:"clientId"`
			Status   string              `json:"status"`
			Service  models.ServiceModel `json:"service"`
			Extras   []models.ExtraModel `json:"extras"`
		}{
			Id:       transaction.Id,
			Cost:     cost,
			VendorId: transaction.VendorId,
			ClientId: transaction.ClientId,
			Status:   string(transaction.Status),
			Service:  transaction.Service,
			Extras:   transaction.Extras,
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": response,
	})
}

func (h *transactionHandler) GetOngoingTransactions(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	transactions, err := h.transactionService.GetOngoingTransactions(bearerToken)
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

func (h *transactionHandler) GetTransactionHistory(c echo.Context) error {
	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	transactions, err := h.transactionService.GetTransactionHistory(bearerToken)
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

func (h *transactionHandler) CompleteTransaction(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Transaction ID is required",
			Error:   "Transaction ID is required",
		})
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	if err := h.transactionService.CompleteTransaction(bearerToken, transactionId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error marking transaction as complete",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
