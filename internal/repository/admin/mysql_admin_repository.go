package admin_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

	if id, err := gonanoid.New(); err != nil {
		return err
	} else {
		data.Id = id
	}

	query := "INSERT INTO Admin (id, username, password, usernameHash, role) VALUES (:id, :username, :password, :usernameHash, :role)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlAdminRepository) Login(data *models.SessionModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := gonanoid.New(); err != nil {
		return err
	} else {
		data.Id = id
	}

	query := "INSERT INTO Session (id, refreshToken) VALUES (:id, :refreshToken)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlAdminRepository) Logout(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	updateSession := "UPDATE Session SET status = 'offline' WHERE refreshToken = ? AND status = 'online'"
	if _, err := s.db.ExecContext(ctx, updateSession, refreshToken); err != nil {
		return err
	}

	id, err := gonanoid.New()
	if err != nil {
		return err
	}

	blacklistToken := `INSERT INTO Blacklist (id, token) VALUES (?, ?)`
	if _, err := tx.ExecContext(ctx, blacklistToken, id, refreshToken); err != nil {
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

func (s *MysqlAdminRepository) FindById(id string) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	admin := new(models.AdminModel)

	query := "SELECT id, username, password, role FROM Admin WHERE id = ?"
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

	query := "SELECT id, username, password, role FROM Admin WHERE usernameHash = ?"
	if err := s.db.GetContext(ctx, admin, query, hash); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}

func (s *MysqlAdminRepository) DoesRefreshTokenExists(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Session WHERE refreshToken = ? AND status = 'online'"
	if err := s.db.GetContext(ctx, &count, query, refreshToken); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	if count > 0 {
		return nil
	}

	return errors.New("refreshToken not found")
}

func (s *MysqlAdminRepository) IsRefreshTokenBlacklisted(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Blacklist WHERE token = ?"
	if err := s.db.GetContext(ctx, &count, query, refreshToken); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	if count > 0 {
		return nil
	}

	return errors.New("refreshToken blacklisted")
}

func (s *MysqlAdminRepository) UserCount(filter UserStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM User"

	switch filter {
	case USER_STATUS_VERIFIED:
		query += " WHERE verified = 1"
	case USER_STATUS_UNVERIFIED:
		query += " WHERE verified = 0"
	case USER_STATUS_ALL:
	}

	count := 0
	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlAdminRepository) VendorCount(filter VendorStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM Vendor"

	switch filter {
	case VENDOR_STATUS_RESTRICTED:
		query += " WHERE restricted = 1"
	case VENDOR_STATUS_UNRESTRICTED:
		query += " WHERE restricted = 0"
	case VENDOR_STATUS_ALL:
	}

	count := 0
	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlAdminRepository) ApplicationCount(filter ApplicationStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM Application"

	switch filter {
	case APPLICATION_STATUS_PENDING:
		query += " WHERE status = 'pending'"
	case APPLICATION_STATUS_APPROVED:
		query += " WHERE status = 'approved'"
	case APPLICATION_STATUS_REJECTED:
		query += " WHERE status = 'rejected'"
	}

	count := 0
	err := s.db.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *MysqlAdminRepository) ComplaintCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(id) FROM Complaint"

	count := 0
	if err := s.db.GetContext(ctx, &count, query); err != nil {
		return 0, nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}
