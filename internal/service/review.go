package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/service/auth"
	"nearbyassist/internal/store/review"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type ReviewService struct {
	store review.ReviewStore
	jwt   auth.Authenticator
}

func NewReviewService(store review.ReviewStore, jwt auth.Authenticator) *ReviewService {
	return &ReviewService{
		store: store,
		jwt:   jwt,
	}
}

func (s *ReviewService) Create(c echo.Context) error {
	req := new(request.NewReviewPayload)
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

	review, err := s.store.FindById(req.TransactionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	}

	if err := s.store.IsServiceReviewable(review.ServiceId); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service is not reviewable",
			Error:   err.Error(),
		})
	}

	token := c.Request().Header.Get("Authorization")[len("Bearer "):]
	userId, err := utils.GetUserIdFromToken(token, s.jwt.GetClaims)

	transaction, err := s.store.GetTransactionById(req.TransactionId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Transaction not found",
			Error:   err.Error(),
		})
	}

	if transaction.ClientId != userId {
		return echo.NewHTTPError(http.StatusForbidden, models.Error{
			Message: "User is not the client of the transaction",
			Error:   "User is not the client of the transaction",
		})
	}

	newReview := new(models.ReviewModel)
	newReview.TransactionId = req.TransactionId
	newReview.ServiceId = req.ServiceId
	newReview.Rating = req.Rating

	reviewId, err := s.store.Create(newReview)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating review",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"review": reviewId,
	})
}

func (s *ReviewService) GetById(c echo.Context) error {
	reviewId := c.Param("reviewId")
	if reviewId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Review ID is required",
			Error:   "Review ID is required",
		})
	}

	review, err := s.store.FindById(reviewId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Review not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"review": review,
	})
}

func (s *ReviewService) GetByService(c echo.Context) error {
	serviceId := c.Param("serviceId")
	if serviceId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Service ID is required",
			Error:   "Service ID is required",
		})
	}

	reviews, err := s.store.GetReviewsByService(serviceId)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, models.Error{
			Message: "Reviews not found",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"reviews": reviews,
	})
}
