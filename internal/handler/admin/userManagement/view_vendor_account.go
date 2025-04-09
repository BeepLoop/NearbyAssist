package userManagement

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) GetVendor(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	data, err := h.managementService.GetVendorAccountDetail(c.Param("userId"))
	if err != nil {
		page := pages.VendorAccountDetail(*admin, dto.VendorAccountDetail{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.VendorAccountDetail(*admin, *data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
