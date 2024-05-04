package handlers

import (
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
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

func (h *complaintHandler) HandleBaseRoute(c echo.Context) error {
	return c.JSON(http.StatusOK, utils.Mapper{
		"message": "Complaints base route",
	})
}

func (h *complaintHandler) HandleCount(c echo.Context) error {
	count, err := h.server.DB.CountComplaint()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"complaintsCount": count,
	})
}

func (h *complaintHandler) HandleNewComplaint(c echo.Context) error {
	req := &request.NewComplaint{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	complaintId, err := h.server.DB.FileComplaint(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"message":     "complaint created successfully",
		"complaintId": complaintId,
	})
}
