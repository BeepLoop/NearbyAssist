package admin_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
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

	query := `
        INSERT INTO
            Admin (id, username, password, usernameHash)
        VALUES
            (:id, :username, :password, :usernameHash)
    `
	data.Id = utils.GenerateUserId()
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

	query := `
        SELECT
            id, username, email, password, mustChangePassword, suspended, role, createdAt, updatedAt
        FROM
            Admin
        WHERE
            id = ?
    `
	admin := new(models.AdminModel)
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

	query := "SELECT id FROM Admin WHERE usernameHash = ?"
	var id string
	if err := s.db.GetContext(ctx, &id, query, hash); err != nil {
		return nil, err
	}

	admin, err := s.FindById(id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}

func (s *MysqlAdminRepository) GetAll() ([]*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            id, username, email, mustChangePassword, suspended, role, createdAt, updatedAt
        FROM
            Admin
        ORDER BY
            createdAt DESC
    `

	accounts := make([]*models.AdminModel, 0)
	if err := s.db.SelectContext(ctx, &accounts, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return accounts, nil
}

func (s *MysqlAdminRepository) GetAllWithRole(role string) ([]*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            id, username, email, mustChangePassword, suspended, role, createdAt, updatedAt
        FROM
            Admin
        WHERE
            role = ?
        ORDER BY
            createdAt DESC
    `

	accounts := make([]*models.AdminModel, 0)
	if err := s.db.SelectContext(ctx, &accounts, query, role); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return accounts, nil
}

func (s *MysqlAdminRepository) DoesUsernameExists(usernamehash string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            CASE
                WHEN EXISTS (SELECT 1 FROM Admin WHERE usernameHash = ?)
                THEN 1
                ELSE 0
            END AS account_exists;
    `

	doesExist := false
	if err := s.db.GetContext(ctx, &doesExist, query, usernamehash); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return doesExist, nil

}

func (s *MysqlAdminRepository) Suspend(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Admin
        SET
            suspended = 1
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlAdminRepository) Unsuspend(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Admin
        SET
            suspended = 0
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
