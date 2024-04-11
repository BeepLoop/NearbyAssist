package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"net/http"

	"github.com/labstack/echo/v4"
)

type complaintHandler struct {
	server *server.Server
}

func NewComplaintServer(server *server.Server) *complaintHandler {
	return &complaintHandler{
		server: server,
	}
}

func (h *complaintHandler) HandleCount(c echo.Context) error {
	count, err := h.server.DB.CountComplaint()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"complaintsCount": count,
	})
}

func (h *complaintHandler) HandleNewComplaint(c echo.Context) error {
	model := models.NewComplaintModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	complaintId, err := h.server.DB.FileComplaint(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":     "complaint created successfully",
		"complaintId": complaintId,
	})
}
