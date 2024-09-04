package handlers

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"nearbyassist/internal/server"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type analyticsHandler struct {
	server *server.Server
}

func NewAnalyticsHandler(server *server.Server) *analyticsHandler {
	return &analyticsHandler{
		server: server,
	}
}

func (h *analyticsHandler) HandleAnalytics(c echo.Context) error {
	analytics := new(response.Analytics)

	// Count user
	user := models.NewUserModel(h.server.IdGen, h.server.DB)
	if user == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to initialize model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if count, err := user.Count(models.USER_STATUS_ALL); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve user count",
			Error:   err.Error(),
		})
	} else {
		analytics.User = count
	}

	if count, err := user.Count(models.USER_STATUS_VERIFIED); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve verified user count",
			Error:   err.Error(),
		})
	} else {
		analytics.VerifiedUser = count
	}

	// Count vendor
	vendor := models.NewVendorModel(h.server.IdGen, h.server.DB)
	if vendor == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to initialize model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if count, err := vendor.Count(models.VENDOR_STATUS_ALL); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve vendor count",
			Error:   err.Error(),
		})
	} else {
		analytics.Vendor = count
	}

	// Count applications
	application := models.NewApplicationModel(h.server.IdGen, h.server.DB)
	if application == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to initialize model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if count, err := application.Count(models.APPLICATION_STATUS_PENDING); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve application count",
			Error:   err.Error(),
		})
	} else {
		analytics.PendingApplication = count
	}

	// Count complaints
	complaint := models.NewComplaintModel(h.server.IdGen, h.server.DB)
	if complaint == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to initialize model",
			Error:   models.MODEL_INIT_ERROR,
		})
	}

	if count, err := complaint.Count(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve complaints count",
			Error:   err.Error(),
		})
	} else {
		analytics.Complaint = count
	}

	return c.JSON(200, utils.Mapper{
		"analytics": analytics,
	})
}
