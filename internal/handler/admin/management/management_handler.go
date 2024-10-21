package management

import management_service "nearbyassist/internal/service/management"

type managementHandler struct {
	managementService *management_service.Service
}

func NewHandler(managementService *management_service.Service) *managementHandler {
	return &managementHandler{managementService: managementService}
}
