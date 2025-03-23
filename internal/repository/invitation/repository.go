package invitation_repo

import "nearbyassist/internal/models"

type Repository interface {
	Create(data *models.InvitationModel) (string, error)
	FindById(inviteId string) (*models.InvitationModel, error)
	FindByCode(code string) (*models.InvitationModel, error)
	IsExpired(inviteId string) (bool, error)
	Accept(inviteId, defaultPassword string) error
}
