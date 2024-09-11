package tag

import (
	"nearbyassist/internal/models"

	"github.com/jmoiron/sqlx"
)

type MysqlTagStore struct {
	db *sqlx.DB
}

func NewMysqlTagStore(db *sqlx.DB) *MysqlTagStore {
	return &MysqlTagStore{
		db: db,
	}
}

func (s *MysqlTagStore) Create(data *models.TagModel) error {
	return nil
}

func (s *MysqlTagStore) FindById(id string) (*models.TagModel, error) {
	return nil, nil
}

func (s *MysqlTagStore) FindAll() ([]*models.TagModel, error) {
	return nil, nil
}
