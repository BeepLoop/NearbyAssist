package accountmanagement

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/account_management/accounts"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *accountManagementHandler) AccountList(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	accounts, err := h.adminService.GetAccounts(c.QueryParam("role"))
	if err != nil {
		page := pages.AccountList(*admin, make([]dto.Admin, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.AccountList(*admin, accounts, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
