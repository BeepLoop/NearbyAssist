package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

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

func (h *applicationHandler) HandleNewApplication(c echo.Context) error {
	application := models.NewApplicationModel(h.server.IdGen, h.server.DB)
	if application == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := c.Bind(application); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid data")
	}

	if err := c.Validate(application); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing required fields")
	}

	user := models.NewUserModelWithId(application.ApplicantId, h.server.DB)
	if user.IsVerified() == false {
		return echo.NewHTTPError(http.StatusForbidden, "User not verified")
	}

	if _, err := application.Create(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred creating application")
	}

	return c.JSON(http.StatusCreated, utils.Mapper{
		"applicationId": application.Id,
	})
}

func (h *applicationHandler) HandleCount(c echo.Context) error {
	param := c.QueryParam("filter")
	var filter models.ApplicationStatusFilter
	switch param {
	case "all":
		filter = models.APPLICATION_STATUS_ALL
	case "pending":
		filter = models.APPLICATION_STATUS_PENDING
	case "approved":
		filter = models.APPLICATION_STATUS_REJECTED
	case "rejected":
		filter = models.APPLICATION_STATUS_REJECTED
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, "Invalid parameter found")
	}

	application := models.NewApplicationModel(h.server.IdGen, h.server.DB)
	if application == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	count, err := application.Count(filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"count": count,
	})
}

func (h *applicationHandler) HandleGetAllApplications(c echo.Context) error {
	application := models.NewApplicationModel(h.server.IdGen, h.server.DB)
	if application == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

    params := utils.ParseQuery(c.QueryString())
	applications, err := application.FindAll(params)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while retrieving applications")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"applications": applications,
	})
}

func (h *applicationHandler) HandleApprove(c echo.Context) error {
	applicationId := c.Param("applicationId")
	if applicationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "application ID must be a number")
	}

	application := models.NewApplicationModel(h.server.IdGen, h.server.DB)
	if application == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := application.Approve(applicationId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while approving application")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "Application approved successfully",
		"applicationId": applicationId,
	})
}

func (h *applicationHandler) HandleReject(c echo.Context) error {
	applicationId := c.Param("applicationId")
	if applicationId == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "application ID must be a number")
	}

	application := models.NewApplicationModel(h.server.IdGen, h.server.DB)
	if application == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.MODEL_INIT_ERROR)
	}

	if err := application.Reject(applicationId); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Error occurred while rejecting application")
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"message":       "Application rejected successfully",
		"applicationId": applicationId,
	})
}
