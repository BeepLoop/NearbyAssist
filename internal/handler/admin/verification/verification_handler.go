package verification

import verification_service "nearbyassist/internal/service/verification"

type verificationHandler struct {
	verificationService *verification_service.Service
}

func NewHandler(verificationService *verification_service.Service) *verificationHandler {
	return &verificationHandler{verificationService: verificationService}
}
