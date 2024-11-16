package dashboard

import (
	"context"
	"nearbyassist/internal/response"
	dashboard_service "nearbyassist/internal/service/dashboard"
	"nearbyassist/views/pages/dashboard"

	"github.com/labstack/echo/v4"
)

type dashboardHandler struct {
	dashboardService *dashboard_service.Service
}

func NewHandler(dashboardService *dashboard_service.Service) *dashboardHandler {
	return &dashboardHandler{dashboardService: dashboardService}
}

func (h *dashboardHandler) GetDashboard(c echo.Context) error {
	analytics, err := h.dashboardService.GetAnalytics()
	if err != nil {
		page := pages.Dashboard(response.Analytics{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.Dashboard(*analytics)
	return page.Render(context.Background(), c.Response().Writer)
}
