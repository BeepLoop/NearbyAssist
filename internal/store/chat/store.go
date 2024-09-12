package chat

import "nearbyassist/internal/models"

type ChatStore interface {
	GetMessages(user1, user2 string) ([]*models.MessageModel, error)
	GetConversations(userId string) ([]*models.UserModel, error)
}
