package handlers

import (
	"nearbyassist/internal/db/models"
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

	model := models.NewApplicationModel()

	count, err := model.Count(filter)
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

	applicationId, err := model.Create()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"applicationId": applicationId,
	})
}

func (h *applicationHandler) HandleGetApplicants(c echo.Context) error {
	filter := c.QueryParam("filter")

	model := models.NewApplicationModel()

	applications, err := model.FindAll(filter)
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

	model := models.NewApplicationModel()

	err = model.Approve(id)
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

	model := models.NewApplicationModel()

	err = model.Reject(id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":       "Application rejected successfully",
		"applicationId": id,
	})
}
