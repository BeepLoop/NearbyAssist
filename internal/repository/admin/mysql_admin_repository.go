package admin_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlAdminRepository struct {
	db *sqlx.DB
}

func NewMysqlAdminRepository(db *sqlx.DB) *MysqlAdminRepository {
	return &MysqlAdminRepository{db: db}
}

func (s *MysqlAdminRepository) Create(data *models.AdminModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateUserId()

	query := `
        INSERT INTO
            Admin (id, username, password, usernameHash)
        VALUES
            (:id, :username, :password, :usernameHash)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlAdminRepository) FindById(id string) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	admin := new(models.AdminModel)

	query := "SELECT id, username, password, mustChangePassword FROM Admin WHERE id = ?"
	if err := s.db.GetContext(ctx, admin, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}

func (s *MysqlAdminRepository) FindByUsernameHash(hash string) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	admin := new(models.AdminModel)

	query := "SELECT id, username, password, mustChangePassword FROM Admin WHERE usernameHash = ?"
	if err := s.db.GetContext(ctx, admin, query, hash); err != nil {
		return nil, err
	}

	if admin.Id == "" {
		return nil, errors.New("not found")
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}

func (s *MysqlAdminRepository) ShouldChangePassword(id string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT mustChangePassword FROM Admin WHERE id = ?"
	shouldChange := false
	if err := s.db.GetContext(ctx, &shouldChange, query, id); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return shouldChange, nil
}

func (s *MysqlAdminRepository) RequestPasswordReset(data *models.PasswordResetRequestModel) error {
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
				return err
			}

			return nil
		}

		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
