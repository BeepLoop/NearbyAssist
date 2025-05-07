package userManagement

import (
	"context"
	"nearbyassist/internal/dto"
	"nearbyassist/internal/service/activitylog"
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

	activity := activitylog.Input{
		AdminId:    admin.Id,
		Action:     activitylog.ACTION_VIEWED_USER,
		TargetType: "user",
		TargetId:   data.User.Id,
	}
	activitylog.MustGetInstance().CreateWithTarget(activity)

	page := pages.UserAccount(*admin, *data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}
