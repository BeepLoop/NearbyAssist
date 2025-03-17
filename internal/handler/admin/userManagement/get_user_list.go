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

		user, err := h.userService.FindByEmail(query)
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

		accounts, err := h.userService.GetAllBasicUsers(limit, offset)
		if err != nil {
			page := pages.UserList(make([]models.UserModel, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = accounts
	}

	data := make([]models.UserModel, 0)
	for _, account := range results {
		if account.VerifiedAt.Valid {
			account.VerifiedAt.String = utils.FormatDate(account.VerifiedAt.String)
		}

		data = append(data, models.UserModel{
			Model: models.Model{
				Id:        account.Id,
				CreatedAt: utils.FormatDate(account.CreatedAt),
			},
			Name:       account.Name,
			Email:      account.Email,
			ImageUrl:   account.ImageUrl,
			Verified:   account.Verified,
			VerifiedAt: account.VerifiedAt,
		})
	}

	page := pages.UserList(data)
	return page.Render(context.Background(), c.Response().Writer)
}
