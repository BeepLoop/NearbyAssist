package expertise_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MysqlExpertiseRepository struct {
	db *sqlx.DB
}

func NewMysqlExpertiseRepository(db *sqlx.DB) *MysqlExpertiseRepository {
	return &MysqlExpertiseRepository{
		db: db,
	}
}

func (s *MysqlExpertiseRepository) Create(data *models.ExpertiseModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if generatedId, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = generatedId
	}

	createQuery := "INSERT INTO Expertise (id, title) VALUES (:id, :title)"
	if _, err := s.db.NamedExecContext(ctx, createQuery, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlExpertiseRepository) CreateTag(expertiseId string, data *models.TagModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if generatedId, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = generatedId
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	createTagQuery := "INSERT INTO Tag (id, title) VALUES (:id, :title)"
	if _, err := tx.NamedExecContext(ctx, createTagQuery, data); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	createRelationshipQuery := "INSERT INTO ExpertiseTag (expertiseId, tagId) VALUES (?, ?)"
	if _, err := tx.ExecContext(ctx, createRelationshipQuery, expertiseId, data.Id); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlExpertiseRepository) GetAll() ([]*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getExpertiseQuery := "SELECT * FROM Expertise"
	expertise := make([]*models.ExpertiseModel, 0)
	if err := s.db.SelectContext(ctx, &expertise, getExpertiseQuery); err != nil {
		return nil, err
	}

	getExpertiseTagsQuery := `
        SELECT
            t.id,
            t.title,
            t.createdAt
        FROM
            ExpertiseTag et
            JOIN Tag t ON t.id = et.tagId
        WHERE
            et.expertiseId = ?
    `

	for _, entry := range expertise {
		tags := make([]*models.TagModel, 0)
		if err := s.db.SelectContext(ctx, &tags, getExpertiseTagsQuery, entry.Id); err != nil {
			return nil, err
		}

		entry.Tags = tags
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertise, nil
}

func (s *MysqlExpertiseRepository) FindById(id string) (*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getExpertiseQuery := "SELECT * FROM Expertise WHERE id = ?"
	expertise := new(models.ExpertiseModel)
	if err := s.db.GetContext(ctx, &expertise, getExpertiseQuery, id); err != nil {
		return nil, err
	}

	if expertise.Id == "" {
		return nil, errors.New("not found")
	}

	getTagsQuery := `
        SELECT
            t.id,
            t.title,
            t.createdAt
        FROM
            ExpertiseTag et
            JOIN Tag t ON t.id = et.tagId
        WHERE
            et.expertiseId = ?
    `

	tags := make([]*models.TagModel, 0)
	if err := s.db.SelectContext(ctx, &tags, getTagsQuery, expertise.Id); err != nil {
		return nil, err
	}
	expertise.Tags = tags

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertise, nil
}
