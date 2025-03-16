package userManagement

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) ViewUserAccount(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)

	userId := c.Param("userId")
	if userId == "" {
		page := pages.ViewUserAccount(models.UserAccountPageData{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	accountData, err := h.managementService.GetSingleUser(userId)
	if err != nil {
		page := pages.ViewUserAccount(models.UserAccountPageData{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	for _, service := range accountData.Services {
		for _, image := range service.Images {
			signedURL, err := h.resourceService.SignURLWithDefaultDuration(image.Url)
			if err != nil {
				fmt.Println("Error generating signed url: ", err.Error())
				continue
			}

			image.Url = signedURL
		}
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
		Banned:     accountData.Banned,
		Restricted: accountData.Restricted,
		Stat:       accountData.Stat,
	}

	page := pages.ViewUserAccount(data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
