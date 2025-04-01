package userManagement

import (
	"context"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/service/cache"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) ViewVendorAccount(c echo.Context) error {
	params := c.QueryParams()

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	flash, _, _ := utils.RetrieveFlashMessage(c)

	userId := c.Param("userId")
	if userId == "" {
		page := pages.VendorAccountDetail(*admin, models.UserAccountPageData{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	var accountData *models.UserAccountPageData

	if params.Has("fresh") && params.Get("fresh") == "true" {
		res, err := h.managementService.GetSingleUser(userId)
		if err != nil {
			page := pages.VendorAccountDetail(*admin, models.UserAccountPageData{}, flash)
			return page.Render(context.Background(), c.Response().Writer)
		}

		accountData = res
	} else {
		inCache, exists := cache.NewGoCache().Get(c.Request().RequestURI)
		if exists {
			accountData = inCache.(*models.UserAccountPageData)
		} else {
			res, err := h.managementService.GetSingleUser(userId)
			if err != nil {
				page := pages.VendorAccountDetail(*admin, models.UserAccountPageData{}, flash)
				return page.Render(context.Background(), c.Response().Writer)
			}

			accountData = res
		}
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

	page := pages.VendorAccountDetail(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
