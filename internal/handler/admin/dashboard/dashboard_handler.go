package dashboard

import (
	"context"
	"nearbyassist/internal/dto"
	dashboard_service "nearbyassist/internal/service/dashboard"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/dashboard"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dashboardHandler struct {
	dashboardService *dashboard_service.Service
}

func NewHandler(dashboardService *dashboard_service.Service) *dashboardHandler {
	return &dashboardHandler{dashboardService: dashboardService}
}

func (h *dashboardHandler) GetDashboard(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	data, err := h.dashboardService.GetDashbaordData()
	if err != nil {
		page := pages.Dashboard(*admin, dto.Dashboard{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.Dashboard(*admin, *data)
	return page.Render(context.Background(), c.Response().Writer)
}
