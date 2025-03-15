package complaint

import (
	complaint_service "nearbyassist/internal/service/complaint"
	resource_service "nearbyassist/internal/service/resource"
)

const (
	DEFAULT_LIMIT  = 10
	DEFAULT_OFFSET = 0
)

type complaintHandler struct {
	complaintService *complaint_service.Service
	resourceService  *resource_service.Service
}

func NewHandler(complaintService *complaint_service.Service, resourceService *resource_service.Service) *complaintHandler {
	return &complaintHandler{
		complaintService: complaintService,
		resourceService:  resourceService,
	}
}
