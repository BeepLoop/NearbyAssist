package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

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
	review := models.NewReviewModel(h.server.IdGen, h.server.DB)
	if review == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := c.Bind(review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if _, err := transaction.FindById(review.TransactionId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Transaction ID not found")
	}

	// Validate that the transaction is marked as done
	if transaction.IsReviewable() == false {
		return echo.NewHTTPError(http.StatusForbidden, "Transaction does not meet requirements for review")
	}

	// Validate that user is the client of the given transaction
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	if id, ok := claims["userId"].(string); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "User ID not found in JWT")
	} else {
		if id != transaction.ClientId {
			return echo.NewHTTPError(http.StatusForbidden, "You are not allowed to review this transaction")
		}
	}

	if _, err := review.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while creating review")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":  "Review created successfully!",
		"reviewId": review.Id,
	})
}

func (h *reviewHandler) HandleGetReview(c echo.Context) error {
	reviewId := c.Param("reviewId")
	if reviewId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "review ID must be a number")
	}

	review := models.NewReviewModel(h.server.IdGen, h.server.DB)
	if review == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if _, err := review.FindById(reviewId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Review not found")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"review": review,
	})
}

func (h *reviewHandler) HandleServiceReview(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "service ID must be a number")
	}

	service := models.NewServiceModelWithId(serviceId, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	reviews, err := service.GetReviews()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "service not found")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"reviews": reviews,
	})
}
