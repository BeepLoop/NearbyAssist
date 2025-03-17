package passwordreset_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlPasswordResetRepository struct {
	db *sqlx.DB
}

func NewMysqlPasswordResetRepository(db *sqlx.DB) *MysqlPasswordResetRepository {
	return &MysqlPasswordResetRepository{
		db: db,
	}
}

func (s *MysqlPasswordResetRepository) Create(data *models.PasswordResetRequestModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data.Id = utils.GenerateId()

	insertQuery := `
        INSERT INTO
            PasswordResetRequest (id, adminId)
        VALUES
            (:id, :adminId)
    `

	updateTimestampQuery := "UPDATE PasswordResetRequest SET createdAT = NOW() WHERE adminId = ?"

	_, err := s.db.NamedExecContext(ctx, insertQuery, data)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			if _, err := s.db.ExecContext(ctx, updateTimestampQuery, data.AdminId); err != nil {
				return "", err
			}

			return data.Id, nil
		}

		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlPasswordResetRepository) Delete(id string) error {
	return nil
}

func (s *MysqlPasswordResetRepository) GetAll() ([]*models.PasswordResetRequestModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            prr.id,
            prr.adminId,
            prr.createdAt,
            a.username
        FROM
            PasswordResetRequest prr
            JOIN Admin a ON a.id = prr.adminId
        ORDER By
            prr.createdAt DESC
    `

	requests := make([]*models.PasswordResetRequestModel, 0)
	if err := s.db.SelectContext(ctx, &requests, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return requests, nil
}

func (s *MysqlPasswordResetRepository) FindById(id string) (*models.PasswordResetRequestModel, error) {
	return nil, nil
}

func (s *MysqlPasswordResetRepository) FindByAdminId(id string) (*models.PasswordResetRequestModel, error) {
	return nil, nil
}
