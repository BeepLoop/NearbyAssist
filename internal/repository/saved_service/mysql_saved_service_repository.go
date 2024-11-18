package saved_service_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MysqlSavedServiceRepository struct {
	db *sqlx.DB
}

func NewMysqlSavedServiceRepository(db *sqlx.DB) *MysqlSavedServiceRepository {
	return &MysqlSavedServiceRepository{db: db}
}

func (s *MysqlSavedServiceRepository) SaveService(data *models.SavedServiceModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := gonanoid.New(); err != nil {
		return err
	} else {
		data.Id = id
	}

	query := "INSERT INTO SavedService (id, userId, serviceId) VALUES (:id, :userId, :serviceId)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlSavedServiceRepository) UnsaveService(data *models.SavedServiceModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "DELETE FROM SavedService WHERE userId = :userId AND serviceId = :serviceId"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlSavedServiceRepository) FindByUserId(userId string) ([]*models.SavedServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	services := make([]*models.SavedServiceModel, 0)
	query := "SELECT * FROM SavedService WHERE userId = ?"
	if err := s.db.SelectContext(ctx, &services, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}
