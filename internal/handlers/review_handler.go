package handlers

import (
	"nearbyassist/internal/models"
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
	model := models.NewReviewModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	reviewId, err := h.server.DB.CreateReview(model)
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
