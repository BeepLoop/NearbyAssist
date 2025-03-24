package userManagement

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) GetVendorList(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	params := c.QueryParams()

	results := make([]*models.VendorModel, 0)

	// If query exists, search is performed
	if params.Has("query") && params.Get("query") != "" {
		query := params.Get("query")

		user, err := h.vendorService.FindByEmail(query)
		if err != nil {
			page := pages.VendorList(*admin, make([]models.VendorModel, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}

		results = append(results, user)
	} else {
		limit, _ := strconv.Atoi(params.Get("limit"))
		if limit == 0 {
			limit = DEFAULT_LIMIT
		}

		offset, _ := strconv.Atoi(params.Get("offset"))
		if offset == 0 {
			offset = DEFAULT_OFFSET
		}

		accounts, err := h.vendorService.GetAll(limit, offset)
		if err != nil {
			page := pages.VendorList(*admin, make([]models.VendorModel, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = accounts
	}

	data := make([]models.VendorModel, 0)
	for _, account := range results {
		data = append(data, models.VendorModel{
			VendorId:    account.VendorId,
			Rating:      account.Rating,
			JoinedAt:    utils.FormatDate(account.JoinedAt),
			Restricted:  account.Restricted,
			Name:        account.Name,
			Email:       account.Email,
			PhoneString: account.PhoneString,
			ImageUrl:    account.ImageUrl,
			Socials:     account.Socials,
			Expertise:   account.Expertise,
		})
	}

	page := pages.VendorList(*admin, data)
	return page.Render(context.Background(), c.Response().Writer)
}
