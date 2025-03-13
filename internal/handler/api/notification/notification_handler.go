package notification

import (
	"nearbyassist/internal/models"
	notification_service "nearbyassist/internal/service/notification"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type notificationHandler struct {
	service *notification_service.Service
}

func NewHandler(service *notification_service.Service) *notificationHandler {
	return &notificationHandler{service: service}
}

func (h *notificationHandler) GetNotifications(c echo.Context) error {
	status := "all"
	statusParam := c.QueryParam("status")
	if statusParam == "unread" {
		status = statusParam
	}

	bearerToken := c.Request().Header.Get("Authorization")[len("Bearer "):]

	notifications, err := h.service.GetNotifications(bearerToken, status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error retrieving notifications",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"notifications": notifications,
	})
}

func (h *notificationHandler) ReadNotification(c echo.Context) error {
	notificationId := c.Param("notificationId")
	if notificationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, models.Error{
			Message: "notificationId must be valid",
			Error:   "notificationId must be valid",
		})
	}

	if err := h.service.ReadNotification(notificationId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Error occurred marking notification as read",
			Error:   err.Error(),
		})
	}

	return c.JSON(http.StatusNoContent, nil)
}
