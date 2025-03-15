package management

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/account_management"

	"github.com/labstack/echo/v4"
)

func (h *managementHandler) ViewUserAccount(c echo.Context) error {
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
		signedURLs := make([]string, 0)
		for _, image := range service.Images {
			signedURL, err := h.resourceService.SignURLWithDefaultDuration(image)
			if err != nil {
				fmt.Println("Error generating signed url: ", err.Error())
				continue
			}

			signedURLs = append(signedURLs, signedURL)
		}

		service.Images = signedURLs
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
