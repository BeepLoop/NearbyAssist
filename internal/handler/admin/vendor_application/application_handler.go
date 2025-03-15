package application

import (
	application_service "nearbyassist/internal/service/application"
	resource_service "nearbyassist/internal/service/resource"
)

type applicationHandler struct {
	applicationService *application_service.Service
	resourceService    *resource_service.Service
}

func NewHandler(applicationService *application_service.Service, resourceService *resource_service.Service) *applicationHandler {
	return &applicationHandler{
		applicationService: applicationService,
		resourceService:    resourceService,
	}
}
