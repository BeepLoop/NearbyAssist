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

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlExpertiseRepository) GetAll() ([]*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getExpertiseQuery := "SELECT * FROM Expertise ORDER BY createdAt DESC"
	expertise := make([]*models.ExpertiseModel, 0)
	if err := s.db.SelectContext(ctx, &expertise, getExpertiseQuery); err != nil {
		return nil, err
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

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertise, nil
}

func (s *MysqlExpertiseRepository) FindByTitle(title string) (*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getExpertiseQuery := "SELECT id, title, createdAt FROM Expertise WHERE title = ?"
	expertise := new(models.ExpertiseModel)
	if err := s.db.GetContext(ctx, expertise, getExpertiseQuery, title); err != nil {
		fmt.Println("error search: ", err.Error())
		return nil, err
	}

	if expertise.Id == "" {
		return nil, errors.New("not found")
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertise, nil
}
