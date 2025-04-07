package recommendation_handler

import (
	"nearbyassist/internal/models"
	recommendation_service "nearbyassist/internal/service/recommendation"
	resource_service "nearbyassist/internal/service/resource"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type handler struct {
	recommendationService *recommendation_service.Service
	resourceService       *resource_service.Service
}

func NewHandler(recommendationService *recommendation_service.Service, resourceService *resource_service.Service) *handler {
	return &handler{
		recommendationService: recommendationService,
		resourceService:       resourceService,
	}
}

func (h *handler) GetRecommendations(c echo.Context) error {
	params := c.QueryParams()

	limit := 10
	offset := 0

	if params.Has("limit") {
		if parsedLimit, err := strconv.Atoi(params.Get("limit")); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Invalid limit",
				Error:   err.Error(),
			})
		} else {
			limit = parsedLimit
		}
	}

	if params.Has("offset") {
		if parsedOffset, err := strconv.Atoi(params.Get("offset")); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, models.Error{
				Message: "Invalid offset",
				Error:   err.Error(),
			})
		} else {
			offset = parsedOffset
		}
	}

	recommendation, err := h.recommendationService.GetRecommendations(limit, offset)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving recommendations",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, recommendation)
}
