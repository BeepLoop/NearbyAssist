package management

import (
	management_service "nearbyassist/internal/service/management"
	resource_service "nearbyassist/internal/service/resource"
)

type managementHandler struct {
	managementService *management_service.Service
	resourceService   *resource_service.Service
}

func NewHandler(managementService *management_service.Service, resourceService *resource_service.Service) *managementHandler {
	return &managementHandler{
		managementService: managementService,
		resourceService:   resourceService,
	}
}
