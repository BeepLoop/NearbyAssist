package vendor

import (
	"context"
	"nearbyassist/internal/models"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	vendor := new(models.VendorModel)
	query := "SELECT id, vendorId, rating, job, restricted FROM Vendor WHERE vendorId = ?"
	if err := s.db.GetContext(ctx, vendor, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return vendor, nil
}
