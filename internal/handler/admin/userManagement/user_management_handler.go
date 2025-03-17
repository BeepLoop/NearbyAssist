package userManagement

import (
	resource_service "nearbyassist/internal/service/resource"
	user_service "nearbyassist/internal/service/user"
	"nearbyassist/internal/service/user_management_service"
	vendor_service "nearbyassist/internal/service/vendor"
)

type userManagementHandler struct {
	managementService *user_management_service.Service
	userService       *user_service.Service
	vendorService     *vendor_service.Service
	resourceService   *resource_service.Service
}

func NewHandler(managementService *user_management_service.Service, userService *user_service.Service, vendorService *vendor_service.Service, resourceService *resource_service.Service) *userManagementHandler {
	return &userManagementHandler{
		managementService: managementService,
		userService:       userService,
		vendorService:     vendorService,
		resourceService:   resourceService,
	}
}
