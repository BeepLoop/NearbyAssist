package expertise_repo

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
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

	data.Id = utils.GenerateId()

	createQuery := "INSERT INTO Expertise (id, title) VALUES (:id, :title)"
	if _, err := s.db.NamedExecContext(ctx, createQuery, data); err != nil {
		return "", err
	}

	if err := s.CreateTags(data.Id, data.Tags); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlExpertiseRepository) CreateTags(expertiseId string, tags []*models.TagModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	createTagQuery := "INSERT INTO Tag (id, title) VALUES (:id, :title)"
	createRelationshipQuery := "INSERT INTO ExpertiseTag (expertiseId, tagId) VALUES (?, ?)"
	for _, tag := range tags {
		tag.Id = utils.GenerateId()
		if _, err := tx.NamedExecContext(ctx, createTagQuery, tag); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, createRelationshipQuery, expertiseId, tag.Id); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlExpertiseRepository) GetAll() ([]*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getExpertiseQuery := "SELECT * FROM Expertise ORDER BY createdAt DESC"
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
	if err := s.db.GetContext(ctx, expertise, getExpertiseQuery, id); err != nil {
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

func (s *MysqlExpertiseRepository) FindByTitle(title string) (*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getExpertiseQuery := "SELECT id, title, createdAt, updatedAt FROM Expertise WHERE title = ?"
	expertise := new(models.ExpertiseModel)
	if err := s.db.GetContext(ctx, expertise, getExpertiseQuery, title); err != nil {
		fmt.Println("error search: ", err.Error())
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
