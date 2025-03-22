package auth

import (
	adminauth_service "nearbyassist/internal/service/admin_auth"
)

type authHandler struct {
	authService *adminauth_service.Service
}

func NewHandler(adminService *adminauth_service.Service) *authHandler {
	return &authHandler{
		authService: adminService,
	}
}
