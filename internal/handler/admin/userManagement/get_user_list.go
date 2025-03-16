package userManagement

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_LIMIT  = 10
	DEFAULT_OFFSET = 0
)

func (h *userManagementHandler) GetUserList(c echo.Context) error {
	params := c.QueryParams()

	results := make([]*models.UserModel, 0)

	// If query exists, search is performed
	if params.Has("query") && params.Get("query") != "" {
		query := params.Get("query")

		user, err := h.managementService.FindUserByEmail(query)
		if err != nil {
			page := pages.UserList(make([]pages.UserAccounts, 0))
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

		accounts, err := h.managementService.GetUsers(limit, offset)
		if err != nil {
			page := pages.UserList(make([]pages.UserAccounts, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = accounts
	}

	data := make([]pages.UserAccounts, 0)
	for _, account := range results {
		data = append(data, pages.UserAccounts{
			Id:         account.Id,
			Name:       account.Name,
			Email:      account.Email,
			ProfileURL: account.ImageUrl,
			Verified:   account.Verified,
			CreatedAt:  utils.FormatDate(account.CreatedAt),
		})
	}

	page := pages.UserList(data)
	return page.Render(context.Background(), c.Response().Writer)
}
