package management

import (
	"context"
	"nearbyassist/internal/utils"
	"nearbyassist/views/pages/account_management"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_LIMIT  = 10
	DEFAULT_OFFSET = 0
)

func (h *managementHandler) GetAccountManagement(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = DEFAULT_LIMIT
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = DEFAULT_OFFSET
	}

	accounts, err := h.managementService.GetUsers(limit, offset)
	if err != nil {
		page := pages.AccountManagement(make([]pages.UserAccounts, 0))
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]pages.UserAccounts, 0)
	for _, account := range accounts {
		data = append(data, pages.UserAccounts{
			Id:         account.Id,
			Name:       account.Name,
			ProfileURL: account.ImageUrl,
			CreatedAt:  utils.FormatDate(account.CreatedAt),
		})
	}

	page := pages.AccountManagement(data)
	return page.Render(context.Background(), c.Response().Writer)
}
