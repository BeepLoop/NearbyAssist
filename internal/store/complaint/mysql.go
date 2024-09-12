package complaint

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlComplaintStore struct {
	db *sqlx.DB
}

func NewMysqlComplaintStore(db *sqlx.DB) *MysqlComplaintStore {
	return &MysqlComplaintStore{
		db: db,
	}
}

func (s *MysqlComplaintStore) CreateSystemComplaint(data *models.SystemComplaintModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	insertComplaint := `
        INSERT INTO 
            SystemComplaint (id, title, detail)
        VALUES
            (:id, :title, :detail)
    `
	if _, err := tx.NamedExecContext(ctx, insertComplaint, data); err != nil {
		return "", err
	}

	insertImage := `
        INSERT INTO 
            SystemComplaintImage (id, complaintId, url)
        VALUES
            (?, ?, ?)
    `
	for _, url := range data.Images {
		imageId, err := store.GenerateNanoId()
		if err != nil {
			return "", err
		}

		if _, err := tx.ExecContext(ctx, insertImage, imageId, data.Id, url); err != nil {
			return "", nil
		}
	}

	if err := tx.Commit(); err != nil {
		return "", nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}
