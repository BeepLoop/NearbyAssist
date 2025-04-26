package accountmanagement

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/account_management/accounts"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *accountManagementHandler) AccountList(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	flash, _, _ := utils.RetrieveFlashMessage(c)

	filter := c.QueryParam("filter")

	accounts, err := h.adminService.GetAll(filter)
	if err != nil {
		page := pages.AccountList(*admin, make([]models.AdminModel, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.AdminModel, 0)
	for _, account := range accounts {
		data = append(data, models.AdminModel{
			Model: models.Model{
				Id:        account.Id,
				CreatedAt: utils.FormatDate(account.CreatedAt),
			},
			Username:           account.Username,
			Email:              account.Email,
			Role:               account.Role,
			MustChangePassword: account.MustChangePassword,
		})
	}

	page := pages.AccountList(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
