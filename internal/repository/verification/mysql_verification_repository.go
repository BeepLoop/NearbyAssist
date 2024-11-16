package verification_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MysqlVerificationRepository struct {
	db *sqlx.DB
}

func NewMysqlVerificationRepository(db *sqlx.DB) *MysqlVerificationRepository {
	return &MysqlVerificationRepository{
		db: db,
	}
}

func (s *MysqlVerificationRepository) Create(data *models.IdentityVerificationModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := `
        INSERT INTO IdentityVerification 
            (id, user, name, address, idType, idNumber, frontIdImageUrl, backIdImageUrl, faceImageUrl)
        VALUES 
            ( :id, :userId, :name, :address, :idType, :idNumber, :frontIdImageUrl, :backIdImageUrl, :faceImageUrl)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlVerificationRepository) GetAll() ([]*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	requests := make([]*models.IdentityVerificationModel, 0)
	query := `
        SELECT 
            id,
            user AS userId,
            name,
            address,
            idType,
            idNumber,
            frontIdImageUrl,
            backIdImageUrl,
            faceImageUrl,
            status,
            createdAt
        FROM 
            IdentityVerification
    `
	if err := s.db.SelectContext(ctx, &requests, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return requests, nil
}
