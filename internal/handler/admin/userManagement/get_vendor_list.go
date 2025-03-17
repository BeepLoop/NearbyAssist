package userManagement

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) GetVendorList(c echo.Context) error {
	params := c.QueryParams()

	results := make([]*models.VendorModel, 0)

	// If query exists, search is performed
	if params.Has("query") && params.Get("query") != "" {
		query := params.Get("query")

		user, err := h.vendorService.GetVendor(query)
		if err != nil {
			page := pages.UserList(make([]models.UserModel, 0))
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
			page := pages.UserList(make([]models.UserModel, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = accounts
	}

	data := make([]models.VendorModel, 0)
	for _, account := range results {
		data = append(data, models.VendorModel{
			Model: models.Model{
				Id:        account.Id,
				CreatedAt: utils.FormatDate(account.CreatedAt),
			},
			VendorId:    account.VendorId,
			Rating:      account.Rating,
			Restricted:  account.Restricted,
			Name:        account.Name,
			Email:       account.Email,
			PhoneString: account.PhoneString,
			ImageUrl:    account.ImageUrl,
			Socials:     account.Socials,
			Expertise:   account.Expertise,
		})
	}

	page := pages.VendorList(data)
	return page.Render(context.Background(), c.Response().Writer)
}
