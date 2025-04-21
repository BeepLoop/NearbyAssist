package userManagement

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management/user"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_LIMIT  = 10
	DEFAULT_OFFSET = 0
)

func (h *userManagementHandler) GetUserList(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	params := c.QueryParams()

	results := make([]*models.UserModel, 0)

	// If query exists, search is performed
	if params.Has("query") && params.Get("query") != "" {
		query := params.Get("query")

		user, err := h.userService.FindByEmail(query)
		if err != nil {
			page := pages.UserList(*admin, make([]models.UserModel, 0))
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
			page := pages.UserList(*admin, make([]models.UserModel, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		results = accounts
	}

	data := slices.AppendSeq(
		make([]models.UserModel, 0),
		utils.Map(results, func(user *models.UserModel) models.UserModel {
			return models.UserModel{
				Model:      models.Model{Id: user.Id, CreatedAt: utils.FormatDate(user.CreatedAt)},
				Name:       user.Name,
				Email:      user.Email,
				ImageUrl:   user.ImageUrl,
				Verified:   user.Verified,
				VerifiedAt: user.VerifiedAt,
			}
		}),
	)

	page := pages.UserList(*admin, data)
	return page.Render(context.Background(), c.Response().Writer)
}
