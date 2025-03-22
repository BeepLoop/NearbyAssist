package policeclearance_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type mysqlImplementation struct {
	db *sqlx.DB
}

func NewMysqlImplementation(db *sqlx.DB) *mysqlImplementation {
	return &mysqlImplementation{
		db: db,
	}
}

func (s *mysqlImplementation) Create(url string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id := utils.GenerateId()

	query := "INSERT INTO PoliceClearance(id, url) VALUES (?, ?)"
	if _, err := s.db.ExecContext(ctx, query, id, url); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return id, nil
}

func (s *mysqlImplementation) FindById(id string) (*models.PoliceClearanceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM PoliceClearance WHERE id = ?"

	policeClearance := new(models.PoliceClearanceModel)
	if err := s.db.GetContext(ctx, policeClearance, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return policeClearance, nil
}
