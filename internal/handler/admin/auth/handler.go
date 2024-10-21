package auth

import (
	admin_service "nearbyassist/internal/service/admin"
)

type authHandler struct {
	adminService *admin_service.Service
}

func NewHandler(adminService *admin_service.Service) *authHandler {
	return &authHandler{adminService: adminService}
}
