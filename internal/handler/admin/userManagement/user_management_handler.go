package userManagement

import (
	management_service "nearbyassist/internal/service/management"
	resource_service "nearbyassist/internal/service/resource"
	vendor_service "nearbyassist/internal/service/vendor"
)

type userManagementHandler struct {
	managementService *management_service.Service
	vendorService     *vendor_service.Service
	resourceService   *resource_service.Service
}

func NewHandler(managementService *management_service.Service, vendorService *vendor_service.Service, resourceService *resource_service.Service) *userManagementHandler {
	return &userManagementHandler{
		managementService: managementService,
		vendorService:     vendorService,
		resourceService:   resourceService,
	}
}
