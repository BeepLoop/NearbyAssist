package accountmanagement

import (
	admin_service "nearbyassist/internal/service/admin"
	"nearbyassist/internal/utils"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *accountManagementHandler) SuspendAccount(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	if err := h.adminService.Suspend(admin.Id, c.Param("adminId")); err != nil {
		if strings.Contains(err.Error(), admin_service.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "You are not authorized to do this actions"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=unauthorized_action")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Account suspension failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=account_suspension_failed")
			}

		}
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts")
	}

	if err := utils.SetFlashMessage(c, "success", "Suspended account"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?success=suspended_account")
	}

	return c.Redirect(http.StatusSeeOther, "/admina/account-management/accounts")
}

func (h *accountManagementHandler) UnsuspendAccount(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	if err := h.adminService.Unsuspend(admin.Id, c.Param("adminId")); err != nil {
		if strings.Contains(err.Error(), admin_service.ERR_UNAUTHORIZED) {
			if err := utils.SetFlashMessage(c, "error", "You are not authorized to do this actions"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=unauthorized_action")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Removing suspension failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?error=unsuspend_failed")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts")
	}

	if err := utils.SetFlashMessage(c, "success", "Account unsuspended"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/accounts?success=unsuspended_account")
	}

	return c.Redirect(http.StatusSeeOther, "/admina/account-management/accounts")
}
