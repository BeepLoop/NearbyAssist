package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/transaction"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TransactionService struct {
	store     transaction.TransactionStore
	jwt       auth.Authenticator
	encryptor auth.Encryption
}

func NewTransactionService(store transaction.TransactionStore, jwt auth.Authenticator, encryptor auth.Encryption) *TransactionService {
	return &TransactionService{
		store:     store,
		jwt:       jwt,
		encryptor: encryptor,
	}
}

func (s *TransactionService) Create(c echo.Context) error {
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

	// Validate that the date is valid
	if err := utils.ValidateDateRange(req.Start, req.End); err != nil {
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

	transaction := new(models.TransactionModel)
	transaction.ClientId = req.ClientId
	transaction.VendorId = req.VendorId
	transaction.ServiceId = req.ServiceId
	transaction.Start = req.Start
	transaction.End = req.End

	transactionId, err := s.store.Create(transaction)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to create transaction",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transaction": transactionId,
	})
}

func (s *TransactionService) GetMyTransactions(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	transactions, err := s.store.GetMyTransactions(userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting transactions",
			Error:   err.Error(),
		})
	}

	for _, transaction := range transactions {
		if plain, err := s.encryptor.DecryptString(transaction.Vendor); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error decrypting vendor name",
				Error:   auth.DECRYPTION_ERR,
			})
		} else {
			transaction.Vendor = plain
		}

		if plain, err := s.encryptor.DecryptString(transaction.Client); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
				Message: "Error decrypting client name",
				Error:   auth.DECRYPTION_ERR,
			})
		} else {
			transaction.Client = plain
		}
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": transactions,
	})
}

func (s *TransactionService) GetAll(c echo.Context) error {
	transactions, err := s.store.GetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting transactions",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"transactions": transactions,
	})
}

func (s *TransactionService) GetOngoing(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	transactions, err := s.store.GetOngoing(userId)
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

func (s *TransactionService) GetHistory(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	transactions, err := s.store.GetHistory(userId)
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

func (s *TransactionService) Complete(c echo.Context) error {
	transactionId := c.Param("transactionId")
	if transactionId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Transaction ID is required",
			Error:   "Transaction ID is required",
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	if transaction, err := s.store.FindById(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	} else {
		if transaction.ClientId != userId {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Unauthorized access",
				Error:   "Unauthorized access",
			})
		}
	}

	if err := s.store.MarkComplete(transactionId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error marking transaction as complete",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
