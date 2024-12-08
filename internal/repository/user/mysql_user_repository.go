package user_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MysqlUserRepository struct {
	db *sqlx.DB
}

func NewMysqlUserRepository(db *sqlx.DB) *MysqlUserRepository {
	return &MysqlUserRepository{db: db}
}

func (s *MysqlUserRepository) CreateUser(user *models.UserModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	id, err := gonanoid.New()
	if err != nil {
		return "", err
	}
	user.Id = id

	query := "INSERT INTO User (id, name, email, imageUrl, emailHash) VALUES (:id, :name, :email, :imageUrl, :emailHash)"
	if _, err := s.db.NamedExecContext(ctx, query, user); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return id, nil
}

func (s *MysqlUserRepository) Login(data *models.SessionModel) error {
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

func (s *MysqlUserRepository) Logout(refreshToken string) error {
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

func (s *MysqlUserRepository) FindById(id string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := new(models.UserModel)
	query := "SELECT id, name, email, imageUrl, verified, address, latitude, longitude FROM User WHERE id = ?"
	err := s.db.GetContext(ctx, user, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (s *MysqlUserRepository) FindByEmailHash(emailHash string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	user := new(models.UserModel)

	query := "SELECT id, name, email, imageUrl, verified, address, latitude, longitude FROM User WHERE emailHash = ?"
	if err := s.db.GetContext(ctx, user, query, emailHash); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}

func (s *MysqlUserRepository) FindSessionByToken(refreshToken string) (*models.SessionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	session := new(models.SessionModel)
	query := "SELECT id, status, refreshToken FROM Session WHERE refreshToken = ? AND status = 'online'"
	if err := s.db.GetContext(ctx, session, query, refreshToken); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return session, nil
}

func (s *MysqlUserRepository) IsRefreshTokenBlacklisted(refreshToken string) error {
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

func (s *MysqlUserRepository) IsVendor(userId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	count := 0
	query := "SELECT COUNT(id) FROM Vendor WHERE vendorId = ?"
	if err := s.db.GetContext(ctx, &count, query, userId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

func (s *MysqlUserRepository) GetExpertise(userId string) ([]*models.ExpertiseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get vendor expertise
	expertiseQuery := `
        SELECT 
            e.id,
            e.title
        FROM 
            Expertise e
            JOIN VendorExpertise ve ON ve.expertiseId = e.id
        WHERE
            ve.vendorId = ?
    `

	expertises := make([]*models.ExpertiseModel, 0)
	if err := s.db.SelectContext(ctx, &expertises, expertiseQuery, userId); err != nil {
		return nil, err
	}

	// Get all tags for each expertise
	tagQuery := `
        SELECT
            t.id,
            t.title
        FROM 
            Tag t 
            JOIN ExpertiseTag et ON et.tagId = t.id
        WHERE
            et.expertiseId = ?
    `

	for _, expertise := range expertises {
		tags := make([]*models.TagModel, 0)
		if err := s.db.SelectContext(ctx, &tags, tagQuery, expertise.Id); err != nil {
			return nil, err
		}

		expertise.Tags = tags
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return expertises, nil
}
