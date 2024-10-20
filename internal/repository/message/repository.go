package message_repo

import "nearbyassist/internal/models"

type MessageRepository interface {
	Create(data *models.MessageModel) (string, error)
	GetMessages(user1, user2 string) ([]*models.MessageModel, error)
	GetConversations(userId string) ([]*models.UserModel, error)
}
