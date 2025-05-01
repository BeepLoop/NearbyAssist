package message_repo

import "nearbyassist/internal/models"

type MessageRepository interface {
	Create(data *models.MessageModel) (string, error)
	FindById(messageId string) (*models.MessageModel, error)
	GetMessages(user1, user2 string) ([]*models.MessageModel, error)
	MarkSeen(messageId string) error
	GetConversations(userId string) ([]*models.ConversationModel, error)
}
