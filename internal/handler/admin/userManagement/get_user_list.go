package userManagement

import (
	"context"
	"fmt"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management/user"
	"net/http"
	"strconv"
	"strings"

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

	limit, _ := strconv.Atoi(params.Get("limit"))
	if limit == 0 {
		limit = DEFAULT_LIMIT
	}

	offset, _ := strconv.Atoi(params.Get("offset"))
	if offset == 0 {
		offset = DEFAULT_OFFSET
	}

	accounts := make([]dto.User, 0)

	if params.Has("query") && params.Get("query") != "" {
		if strings.HasPrefix(params.Get("query"), "email:") {
			query := params.Get("query")[len("email:"):]

			user, err := h.userService.FindByEmail(query)
			if err != nil {
				fmt.Println("error find by email: ", err.Error())
				page := pages.UserList(*admin, make([]dto.User, 0))
				return page.Render(context.Background(), c.Response().Writer)
			}

			accounts = append(accounts, *user)
		} else if strings.HasPrefix(params.Get("query"), "id:") {
			query := params.Get("query")[len("id:"):]

			user, err := h.userService.FindById(query)
			if err != nil {
				fmt.Println("error find by id: ", err.Error())
				page := pages.UserList(*admin, make([]dto.User, 0))
				return page.Render(context.Background(), c.Response().Writer)
			}

			accounts = append(accounts, *user)
		}
	} else {
		res, err := h.userService.GetAllBasicUsers(limit, offset)
		if err != nil {
			fmt.Println("error get users: ", err.Error())
			page := pages.UserList(*admin, make([]dto.User, 0))
			return page.Render(context.Background(), c.Response().Writer)
		}
		accounts = res
	}

	page := pages.UserList(*admin, accounts)
	return page.Render(context.Background(), c.Response().Writer)
}
