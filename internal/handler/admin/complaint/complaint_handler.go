package complaint

import complaint_service "nearbyassist/internal/service/complaint"

type complaintHandler struct {
	complaintService *complaint_service.Service
}

func NewHandler(complaintService *complaint_service.Service) *complaintHandler {
	return &complaintHandler{complaintService: complaintService}
}
