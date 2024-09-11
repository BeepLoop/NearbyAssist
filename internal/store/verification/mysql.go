package verification

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlVerificationStore struct {
	db *sqlx.DB
}

func NewMysqlVerificationStore(db *sqlx.DB) *MysqlVerificationStore {
	return &MysqlVerificationStore{
		db: db,
	}
}

func (s *MysqlVerificationStore) Create(data *models.IdentityVerificationModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := `
        INSERT INTO IdentityVerification (id, name, address, idType, idNumber, frontId, backId, face)
        VALUES ( :id, :name, :address, :idType, :idNumber, :frontId, :backId, :face)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}
