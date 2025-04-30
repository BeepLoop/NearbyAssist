package invitation

import (
	"nearbyassist/internal/service/invite_service"
)

type handler struct {
	inviteService *invite_service.Service
}

func NewHandler(inviteService *invite_service.Service) *handler {
	return &handler{
		inviteService: inviteService,
	}
}
