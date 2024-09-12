package handler

import (
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"nearbyassist/internal/store/analytics"
	"nearbyassist/internal/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AnalyticsService struct {
	store analytics.AnalyticsStore
}

func NewAnalyticsService(store analytics.AnalyticsStore) *AnalyticsService {
	return &AnalyticsService{
		store: store,
	}
}

func (s *AnalyticsService) Analytics(c echo.Context) error {
	analyticsData := new(response.Analytics)

	if count, err := s.store.UserCount(analytics.USER_STATUS_ALL); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve user count",
			Error:   err.Error(),
		})
	} else {
		analyticsData.User = count
	}

	if count, err := s.store.UserCount(analytics.USER_STATUS_VERIFIED); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve verified user count",
			Error:   err.Error(),
		})
	} else {
		analyticsData.VerifiedUser = count
	}

	if count, err := s.store.VendorCount(analytics.VENDOR_STATUS_ALL); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve vendor count",
			Error:   err.Error(),
		})
	} else {
		analyticsData.Vendor = count
	}

	if count, err := s.store.ApplicationCount(analytics.APPLICATION_STATUS_PENDING); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve application count",
			Error:   err.Error(),
		})
	} else {
		analyticsData.PendingApplication = count
	}

	if count, err := s.store.ComplaintCount(); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, models.Error{
			Message: "Failed to retrieve complaints count",
			Error:   err.Error(),
		})
	} else {
		analyticsData.Complaint = count
	}

	return c.JSON(http.StatusOK, utils.Mapper{
		"analytics": analyticsData,
	})
}
