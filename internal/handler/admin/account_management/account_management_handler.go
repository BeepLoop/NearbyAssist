package accountmanagement

import (
	admin_service "nearbyassist/internal/service/admin"
	passwordreset_service "nearbyassist/internal/service/password_reset"
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
