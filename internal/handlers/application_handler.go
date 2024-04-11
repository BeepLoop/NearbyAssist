package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type applicationHandler struct {
	server *server.Server
}

func NewApplicationHandler(server *server.Server) *applicationHandler {
	return &applicationHandler{
		server: server,
	}
}

func (h *applicationHandler) HandleCount(c echo.Context) error {
	filter := c.QueryParam("filter")
	status := models.ApplicationStatus(filter)

	count, err := h.server.DB.CountApplication(status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"applicationCount": count,
	})
}

func (h *applicationHandler) HandleNewApplication(c echo.Context) error {
	model := models.NewApplicationModel()
	err := c.Bind(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "unable to process data provided")
	}

	err = c.Validate(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	applicationId, err := h.server.DB.CreateApplication(model)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"applicationId": applicationId,
	})
}

func (h *applicationHandler) HandleGetApplications(c echo.Context) error {
	filter := c.QueryParam("filter")
	status := models.ApplicationStatus(filter)

	applications, err := h.server.DB.FindAllApplication(status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, applications)
}

func (h *applicationHandler) HandleApprove(c echo.Context) error {
	applicationId := c.Param("applicationId")
	id, err := strconv.Atoi(applicationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "application ID must be a number")
	}

	err = h.server.DB.ApproveApplication(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Application approved successfully",
		"applicationId": id,
	})
}

func (h *applicationHandler) HandleReject(c echo.Context) error {
	applicationId := c.Param("applicationId")
	id, err := strconv.Atoi(applicationId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "application ID must be a number")
	}

	err = h.server.DB.RejectApplication(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Application rejected successfully",
		"applicationId": id,
	})
}
