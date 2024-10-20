package admin

import (
	"context"
	admin_repo "nearbyassist/internal/repository/admin"
	admin_service "nearbyassist/internal/service/admin"
	"nearbyassist/views/pages"

	"github.com/labstack/echo/v4"
)

func (h *adminHandler) GetIdentityVerification(c echo.Context) error {
	adminService := admin_service.NewService(
		admin_repo.NewMysqlAdminRepository(h.db),
		h.encrypt,
		h.hash,
	)

	requests, err := adminService.IdentityVerificationRequests()
	if err != nil {
		//
	}

	page := pages.IdentityVerification(requests)
	return page.Render(context.Background(), c.Response().Writer)
}
