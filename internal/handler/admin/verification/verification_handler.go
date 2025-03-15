package verification

import (
	resource_service "nearbyassist/internal/service/resource"
	verification_service "nearbyassist/internal/service/verification"
)

type verificationHandler struct {
	verificationService *verification_service.Service
	resourceService     *resource_service.Service
}

func NewHandler(verificationService *verification_service.Service, resourceService *resource_service.Service) *verificationHandler {
	return &verificationHandler{
		verificationService: verificationService,
		resourceService:     resourceService,
	}
}
