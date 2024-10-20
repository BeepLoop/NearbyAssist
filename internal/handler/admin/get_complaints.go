package admin

import (
	"context"
	admin_repo "nearbyassist/internal/repository/admin"
	admin_service "nearbyassist/internal/service/admin"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *adminHandler) GetComplaints(c echo.Context) error {
	adminService := admin_service.NewService(
		admin_repo.NewMysqlAdminRepository(h.db),
		h.encrypt,
		h.hash,
	)

	complaints, err := adminService.Complaints()
	if err != nil {
		//
	}

	page := pages.Complaints(complaints)
	return page.Render(context.Background(), c.Response().Writer)
}
