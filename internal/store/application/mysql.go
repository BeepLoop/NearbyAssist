package application

import (
	"nearbyassist/internal/models"

	"github.com/jmoiron/sqlx"
)

type MysqlApplicationStore struct {
	db *sqlx.DB
}

func NewMysqlApplicationStore(db *sqlx.DB) *MysqlApplicationStore {
	return &MysqlApplicationStore{
		db: db,
	}
}

func (s *MysqlApplicationStore) Create(data *models.ApplicationModel) (string, error) {
	return "", nil
}
