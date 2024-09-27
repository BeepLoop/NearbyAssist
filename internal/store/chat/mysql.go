package chat

import (
	"nearbyassist/internal/models"

	"github.com/jmoiron/sqlx"
)

type MysqlChatStore struct {
	db *sqlx.DB
}

func NewMysqlChatStore(db *sqlx.DB) *MysqlChatStore {
	return &MysqlChatStore{
		db: db,
	}
}
func (s *MysqlChatStore) Create(data *models.MessageModel) (string, error) {
	return "", nil
}

func (s *MysqlChatStore) GetMessages(user1, user2 string) ([]*models.MessageModel, error) {
	return nil, nil
}

func (s *MysqlChatStore) GetConversations(userId string) ([]*models.UserModel, error) {
	return nil, nil
}
