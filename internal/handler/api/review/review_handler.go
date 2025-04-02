package review

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	review_service "nearbyassist/internal/service/review"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type reviewHandler struct {
	reviewService *review_service.Service
}

func NewHandler(reviewService *review_service.Service) *reviewHandler {
	return &reviewHandler{reviewService: reviewService}
}

func (h *reviewHandler) CreateReview(c echo.Context) error {
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

	bearerToken := utils.BearerTokenFromHeader(c)

	reviewId, err := h.reviewService.CreateReview(bearerToken, req)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return echo.NewHTTPError(http.StatusForbidden, models.Error{
				Message: "Posting review not allowed",
				Error:   err.Error(),
			})
		}

		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error creating review",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"review": reviewId,
	})
}

func (h *reviewHandler) GetReview(c echo.Context) error {
	reviewId := c.Param("reviewId")
	if reviewId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "Review ID is required",
			Error:   "Review ID is required",
		})
	}

	review, err := h.reviewService.GetReview(reviewId)
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
