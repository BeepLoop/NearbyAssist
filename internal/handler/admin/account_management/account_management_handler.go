package accountmanagement

import (
	"context"
	"nearbyassist/internal/models"
	admin_service "nearbyassist/internal/service/admin"
	passwordreset_service "nearbyassist/internal/service/password_reset"
	"nearbyassist/internal/utils"
	pages "nearbyassist/views/pages/account_management"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type accountManagementHandler struct {
	adminService         *admin_service.Service
	passwordResetService *passwordreset_service.Service
}

func NewHandler(adminService *admin_service.Service, passwordResetService *passwordreset_service.Service) *accountManagementHandler {
	return &accountManagementHandler{
		adminService:         adminService,
		passwordResetService: passwordResetService,
	}
}

func (h *accountManagementHandler) GetAccounts(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	flash, _, _ := utils.RetrieveFlashMessage(c)

	filter := c.QueryParam("filter")

	accounts, err := h.adminService.GetAll(filter)
	if err != nil {
		page := pages.AccountList(*admin, make([]models.AdminModel, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.AdminModel, 0)
	for _, account := range accounts {
		data = append(data, models.AdminModel{
			Model: models.Model{
				Id:        account.Id,
				CreatedAt: utils.FormatDate(account.CreatedAt),
			},
			Username:           account.Username,
			Email:              account.Email,
			Role:               account.Role,
			MustChangePassword: account.MustChangePassword,
		})
	}

	page := pages.AccountList(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *accountManagementHandler) ResetRequests(c echo.Context) error {
	admin, err := utils.GetAdminFromSession(c)
	if err != nil {
		return c.Redirect(http.StatusSeeOther, "/auth/login")
	}

	flash, _, _ := utils.RetrieveFlashMessage(c)

	requests, err := h.passwordResetService.GetResetRequests()
	if err != nil {
		c.Logger().Warnf("Error getting reset requests: %s", err.Error())
		page := pages.ResetRequests(*admin, make([]models.PasswordResetRequestModel, 0), flash)
		return page.Render(context.Background(), c.Response().Writer)
	}

	data := make([]models.PasswordResetRequestModel, 0)
	for _, request := range requests {
		data = append(data, models.PasswordResetRequestModel{
			Model: models.Model{
				Id:        request.Id,
				CreatedAt: utils.FormatDate(request.CreatedAt),
			},
			AdminId:  request.AdminId,
			Username: request.Username,
		})
	}

	page := pages.ResetRequests(*admin, data, flash)
	return page.Render(context.Background(), c.Response().Writer)
}

func (h *accountManagementHandler) FufillResetRequest(c echo.Context) error {
	requestId := c.FormValue("requestId")
	password := c.FormValue("password")
	confirmationUsername := c.FormValue("confirmationUsername")
	confirmationPassword := c.FormValue("confirmationPassword")

	// Prevent resetting own account
	request, err := h.passwordResetService.GetResetRequest(requestId)
	if err != nil {
		if err := utils.SetFlashMessage(c, "error", "request not found"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=not_found_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
	}

	admin, _ := utils.GetAdminFromSession(c)
	if admin.Id == request.AdminId {
		if err := utils.SetFlashMessage(c, "error", "Resetting own password not allowed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=reset_failed")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
	}

	// Perform reset
	if err := h.passwordResetService.FulfillResetPassword(requestId, password, confirmationUsername, confirmationPassword); err != nil {
		if strings.Contains(err.Error(), "Invalid credentials") {
			if err := utils.SetFlashMessage(c, "error", "Invalid confirmation credentials"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=credential_error")
			}
		} else if strings.Contains(err.Error(), "insecure password") {
			if err := utils.SetFlashMessage(c, "error", "New password not secure enough"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=password_rule_error")
			}
		} else {
			if err := utils.SetFlashMessage(c, "error", "Reset password failed"); err != nil {
				return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=reset_failed")
			}
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
	}

	if err := utils.SetFlashMessage(c, "success", "Request success"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?success=password_change_success")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
}

func (h *accountManagementHandler) RejectResetRequest(c echo.Context) error {
	requestId := c.FormValue("requestId")
	reason := c.FormValue("reason")

	if err := h.passwordResetService.RejectResetPassword(requestId, reason); err != nil {
		if err := utils.SetFlashMessage(c, "error", "Rejection failed"); err != nil {
			return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?error=rejection_error")
		}

		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
	}

	if err := utils.SetFlashMessage(c, "success", "Reject success"); err != nil {
		return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset?success=reject_success")
	}

	return c.Redirect(http.StatusSeeOther, "/admin/account-management/reset")
}
