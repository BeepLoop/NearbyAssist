package userManagement

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management/seller"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) GetVendorList(c echo.Context) error {
	params := c.QueryParams()

	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	results := make([]dto.Vendor, 0)

	// If 'query' exists, perform search
	if params.Has("query") && params.Get("query") != "" {
		user, err := h.vendorService.FindByEmail(params.Get("query"))
		if err != nil {
			page := pages.VendorList(*admin, make([]dto.Vendor, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}

		results = append(results, *user)
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
			page := pages.VendorList(*admin, make([]dto.Vendor, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = accounts
	}

	page := pages.VendorList(*admin, results)
	return page.Render(context.Background(), c.Response().Writer)
}
