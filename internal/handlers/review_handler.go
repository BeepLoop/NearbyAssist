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

type reviewHandler struct {
	server *server.Server
}

func NewReviewHandler(server *server.Server) *reviewHandler {
	return &reviewHandler{
		server: server,
	}
}

func (h *reviewHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Transaction base route",
	})
}

func (h *reviewHandler) HandleNewReview(c echo.Context) error {
	req := &request.NewReview{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Validate that transaction ID exists
	transaction, err := h.server.DB.FindTransactionById(req.TransactionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// Validate that the transaction is marked as done
	if transaction.Status != models.TRANSACTION_STATUS_DONE {
		return echo.NewHTTPError(http.StatusForbidden, "Transaction is not yet completed")
	}

	// Validate that user has not posted a review yet
	if transaction.IsReviewed != false {
		return echo.NewHTTPError(http.StatusForbidden, "Transaction has already been reviewed")
	}

	// Validate that user is the client of the given transaction
	authHeader := c.Request().Header.Get("Authorization")
	if userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		if transaction.ClientId != userId {
			return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to post a review for this transaction")
		}
	}

	reviewId, err := h.server.DB.CreateReview(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":  "Review created successfully!",
		"reviewId": reviewId,
	})
}

func (h *reviewHandler) HandleGetReview(c echo.Context) error {
	reviewId := c.Param("reviewId")
	id, err := strconv.Atoi(reviewId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "review ID must be a number")
	}

	review, err := h.server.DB.FindReviewById(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"review": review,
	})
}

func (h *reviewHandler) HandleServiceReview(c echo.Context) error {
	serviceId := c.Param("serviceId")
	id, err := strconv.Atoi(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	reviews, err := h.server.DB.FindAllReviewByService(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"reviews": reviews,
	})
}
