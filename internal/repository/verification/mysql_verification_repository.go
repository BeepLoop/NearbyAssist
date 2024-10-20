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
