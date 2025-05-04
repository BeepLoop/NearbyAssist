package invitation_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type mysqlRepository struct {
	db *sqlx.DB
}

func NewMysqlRepository(db *sqlx.DB) *mysqlRepository {
	return &mysqlRepository{
		db: db,
	}
}

func (s *mysqlRepository) Create(invitation *models.InvitationModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createInvite := `
        INSERT INTO
            Invitation (username, email, code, password, usernameHash, emailHash, expiredAt)
        VALUES
            (:username, :email, :code, :password, :usernameHash, :emailHash, :expiredAt)
        ON DUPLICATE KEY UPDATE expiredAt = :expiredAt
    `

	invitation.Id = utils.GenerateId()
	if _, err := s.db.NamedExecContext(ctx, createInvite, invitation); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return invitation.Id, nil
}

func (s *mysqlRepository) FindById(inviteId string) (*models.InvitationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM Invitation WHERE id = ?"

	invite := new(models.InvitationModel)
	if err := s.db.GetContext(ctx, invite, query, inviteId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return invite, nil
}

func (s *mysqlRepository) FindByCode(code string) (*models.InvitationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM Invitation WHERE code = ?"

	invite := new(models.InvitationModel)
	if err := s.db.GetContext(ctx, invite, query, code); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return invite, nil
}

func (s *mysqlRepository) IsExpired(inviteId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT CASE
            WHEN (SELECT 1 FROM Invitation WHERE id = ? AND expiredAt < CURRENT_TIMESTAMP())
            THEN 1
            ELSE 0
        END AS is_expired
    `
	isExpired := false
	if err := s.db.GetContext(ctx, &isExpired, query, inviteId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isExpired, nil
}

func (s *mysqlRepository) Accept(inviteId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	createAdmin := `
        INSERT INTO
            Admin(id, username, email, password, usernameHash, emailHash, mustChangePassword)
        SELECT
            ?, username, email, password, usernameHash, emailHash, 1
        FROM
            Invitation
        WHERE
            id = ?
    `
	adminId := utils.GenerateUserId()
	if _, err := tx.ExecContext(ctx, createAdmin, adminId, inviteId); err != nil {
		return err
	}

	removeInvite := "DELETE FROM Invitation WHERE id = ?"
	if _, err := tx.ExecContext(ctx, removeInvite, inviteId); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
