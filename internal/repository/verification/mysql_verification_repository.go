package verification_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
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

	data.Id = utils.GenerateId()

	query := `
        INSERT INTO
            IdentityVerification 
            (id,
            userId,
            name,
            address,
            phone,
            latitude,
            longitude,
            idType,
            idNumber,
            frontIdImageUrl,
            backIdImageUrl,
            faceImageUrl)
        VALUES 
            ( :id, :userId, :name, :address, :phone, :latitude, :longitude, :idType, :idNumber, :frontIdImageUrl, :backIdImageUrl, :faceImageUrl)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlVerificationRepository) GetAll(status string) ([]*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	requests := make([]*models.IdentityVerificationModel, 0)
	query := `
        SELECT 
            id,
            userId,
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
        WHERE
            status = ?
    `
	if err := s.db.SelectContext(ctx, &requests, query, status); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return requests, nil
}

func (s *MysqlVerificationRepository) FindById(id string) (*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	request := new(models.IdentityVerificationModel)
	query := `
        SELECT 
            id,
            userId,
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
        WHERE
            id = ?
    `
	if err := s.db.GetContext(ctx, request, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return request, nil
}

func (s *MysqlVerificationRepository) AcceptRequest(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	updateStatus := "UPDATE IdentityVerification SET status = 'approved' WHERE id = ?"
	if _, err := tx.ExecContext(ctx, updateStatus, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	updateUser := `
        UPDATE
            User u
        JOIN IdentityVerification iv ON u.id = iv.userId
        SET
            u.address = iv.address,
            u.phone = iv.phone,
            u.latitude = iv.latitude,
            u.longitude = iv.longitude,
            u.name = iv.name
        WHERE
            iv.id = ?
    `
	if _, err := tx.ExecContext(ctx, updateUser, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

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

func (s *MysqlVerificationRepository) RejectRequest(id, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        UPDATE
            IdentityVerification
        SET
            status = 'rejected',
            rejectionReason = ?
        WHERE
            id = ?
    `

	if _, err := s.db.ExecContext(ctx, query, reason, id); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
