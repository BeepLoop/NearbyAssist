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
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if err := c.Bind(review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error binding request body",
			Error:   err.Error(),
		})
	}

	if err := c.Validate(review); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Error validating request body",
			Error:   err.Error(),
		})
	}

	transaction := models.NewTransactionModel(h.server.IdGen, h.server.DB)
	if _, err := transaction.FindById(review.TransactionId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	}

	// Validate that the transaction is marked as done
	if transaction.IsReviewable() == false {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "Transaction is not reviewable",
			Error:   "Transaction is not reviewable",
		})
	}

	// Validate that user is the client of the given transaction
	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	claims, err := h.server.Auth.GetClaims(token)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error getting claims",
			Error:   err.Error(),
		})
	}

	if id, ok := claims["userId"].(string); !ok {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "User ID not found in JWT",
			Error:   "User ID not found in JWT",
		})
	} else {
		if id != transaction.ClientId {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "User is not the client of the transaction",
				Error:   "User is not the client of the transaction",
			})
		}
	}

	if _, err := review.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating review",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":  "Review created successfully!",
		"reviewId": review.Id,
	})
}

func (h *reviewHandler) HandleGetReview(c echo.Context) error {
	reviewId := c.Param("reviewId")
	if reviewId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Review ID is required",
			Error:   "Review ID is required",
		})
	}

	review := models.NewReviewModel(h.server.IdGen, h.server.DB)
	if review == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if _, err := review.FindById(reviewId); err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Review not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"review": review,
	})
}

func (h *reviewHandler) HandleServiceReview(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID is required",
			Error:   "Service ID is required",
		})
	}

	service := models.NewServiceModelWithId(serviceId, h.server.DB)
	if service == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error initializing model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	reviews, err := service.GetReviews()
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Error getting reviews",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"reviews": reviews,
	})
}
