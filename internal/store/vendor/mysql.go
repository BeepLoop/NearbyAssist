package vendor

import (
	"nearbyassist/internal/models"

	"github.com/jmoiron/sqlx"
)

type MysqlVendorStore struct {
	db *sqlx.DB
}

func NewMysqlVendorStore(db *sqlx.DB) *MysqlVendorStore {
	return &MysqlVendorStore{
		db: db,
	}
}

func (s *MysqlVendorStore) FindById(id string) (*models.VendorModel, error) {
	return nil, nil
}
