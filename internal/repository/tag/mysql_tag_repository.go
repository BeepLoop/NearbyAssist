package tag_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlTagRepository struct {
	db *sqlx.DB
}

func NewMysqlTagRepository(db *sqlx.DB) *MysqlTagRepository {
	return &MysqlTagRepository{
		db: db,
	}
}

func (s *MysqlTagRepository) Create(data *models.TagModel) error {
	return nil
}

func (s *MysqlTagRepository) FindById(id string) (*models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tag := new(models.TagModel)
	query := "SELECT id, title from Tag WHERE id = ?"
	if err := s.db.GetContext(ctx, tag, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tag, nil
}

func (s *MysqlTagRepository) FindAll() ([]*models.TagModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tags := make([]*models.TagModel, 0)
	query := "SELECT id, title from Tag"
	if err := s.db.SelectContext(ctx, &tags, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return tags, nil
}
