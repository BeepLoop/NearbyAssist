package application

import application_service "nearbyassist/internal/service/application"

type applicationHandler struct {
	applicationService *application_service.Service
}

func NewHandler(applicationService *application_service.Service) *applicationHandler {
	return &applicationHandler{applicationService: applicationService}
}
