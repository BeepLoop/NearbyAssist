package application

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := `
        INSERT INTO
            Application (id, applicantId, job, latitude, longitude)
        VALUES
            (:id, :applicantId, :job, :latitude, :longitude)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}
