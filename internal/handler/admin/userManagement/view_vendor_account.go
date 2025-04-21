package userManagement

import (
	"context"
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management/seller"
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
		fmt.Println(err.Error())
		page := pages.VendorAccount(*admin, dto.VendorAccountDetail{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.VendorAccount(*admin, *data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
