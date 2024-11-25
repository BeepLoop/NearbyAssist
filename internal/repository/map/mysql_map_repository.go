package map_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlMapRepository struct {
	db *sqlx.DB
}

func NewMysqlMapRepository(db *sqlx.DB) *MysqlMapRepository {
	return &MysqlMapRepository{db: db}
}

func (s *MysqlMapRepository) GetAllByTag(tag string) ([]*models.ServiceModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	services := make([]*models.ServiceModel, 0)

	query := `
        SELECT
            s.latitude,
            s.longitude
        FROM 
            ServiceTag st
            JOIN Service s ON s.id = st.serviceId
            JOIN Tag t ON t.id = st.tagId
        WHERE
            t.title = ?
    `

	if err := s.db.SelectContext(ctx, &services, query, tag); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return services, nil
}
