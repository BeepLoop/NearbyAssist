package management

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/account_management"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) ViewUserAccount(c echo.Context) error {
	userId := c.Param("userId")
	if userId == "" {
		page := pages.ViewUserAccount(models.UserAccountPageData{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	accountData, err := h.managementService.GetSingleUser(userId)
	if err != nil {
		page := pages.ViewUserAccount(models.UserAccountPageData{})
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := models.UserAccountPageData{
		Id:         accountData.Id,
		ProfileURL: accountData.ProfileURL,
		Name:       accountData.Name,
		Email:      accountData.Email,
		Address:    accountData.Address,
		CreatedAt:  utils.FormatDate(accountData.CreatedAt),
		Expertise:  accountData.Expertise,
		Services:   accountData.Services,
	}

	page := pages.ViewUserAccount(data)
	return page.Render(context.Background(), c.Response().Writer)
}
