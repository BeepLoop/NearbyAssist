package admin

import (
	"context"
	admin_repo "nearbyassist/internal/repository/admin"
	admin_service "nearbyassist/internal/service/admin"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *adminHandler) GetDashboard(c echo.Context) error {
	adminService := admin_service.NewService(
		admin_repo.NewMysqlAdminRepository(h.db),
		h.encrypt,
		h.hash,
	)

	analytics, err := adminService.Analytics()
	if err != nil {
		//
	}

	page := pages.Dashboard(analytics)
	return page.Render(context.Background(), c.Response().Writer)
}
