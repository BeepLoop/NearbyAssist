package userManagement

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/user_management/user"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *userManagementHandler) GetUser(c echo.Context) error {
	flash, _, _ := utils.RetrieveFlashMessage(c)
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	data, err := h.managementService.GetUserAccountDetail(c.Param("userId"))
	if err != nil {
		page := pages.UserAccount(*admin, dto.UserAccountDetail{}, flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	page := pages.UserAccount(*admin, *data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
