package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/request"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
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

	return c.JSON(http.StatusOK, utils.Mapper{
		"applicationCount": count,
	})
}

func (h *applicationHandler) HandleNewApplication(c echo.Context) error {
	req := &request.NewApplication{}
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "unable to process data provided")
	}

	authHeader := c.Request().Header.Get("Authorization")
	if userId, err := utils.GetUserIdFromJWT(h.server.Auth, authHeader); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	} else {
		req.ApplicantId = userId
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "missing required fields")
	}

	if vendor, _ := h.server.DB.FindVendorById(req.ApplicantId); vendor != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "applicant is already a vendor")
	}

	applicationId, err := h.server.DB.CreateApplication(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
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

	return c.JSON(http.StatusOK, utils.Mapper{
		"applications": applications,
	})
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

	return c.JSON(http.StatusOK, utils.Mapper{
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

	if err := h.server.DB.RejectApplication(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "Application rejected successfully",
		"applicationId": id,
	})
}
