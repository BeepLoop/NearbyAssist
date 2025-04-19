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
	return &MysqlVerificationRepository{db: db}
}

func (s *MysqlVerificationRepository) Create(userId string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            IdentityVerification (id, userId)
        VALUES 
            (?, ?)
    `
	id := utils.GenerateId()
	if _, err := s.db.ExecContext(ctx, query, id, userId); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return id, nil
}

func (s *MysqlVerificationRepository) Update(requestId string, updated *models.IdentityVerificationModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	// Update User name and phone
	updateUser := `
        UPDATE
            User
        SET
            Name = ?, Phone = ?
        WHERE
            id = ?
    `
	if _, err := tx.ExecContext(ctx, updateUser, updated.User.Name, updated.User.Phone, updated.User.Id); err != nil {
		return err
	}

	// Update Address
	updateAddress := `
        UPDATE
            Address
        SET
            address = ?, latitude = ?, longitude = ?
        WHERE
            id = ?
    `
	if _, err := tx.ExecContext(
		ctx,
		updateAddress,
		updated.User.Address.Address,
		updated.User.Address.Latitude,
		updated.User.Address.Longitude,
		updated.User.Address.Id,
	); err != nil {
		return err
	}

	// Update Identification
	updateIdentification := `
        UPDATE
            Identification
        SET
            type = ?,
            referenceNumber = ?,
            frontImageUrl = ?,
            backImageUrl = ?,
            selfieImageUrl = ?
        WHERE
            id = ?
    `
	if _, err := tx.ExecContext(
		ctx,
		updateIdentification,
		updated.User.Identification.Type,
		updated.User.Identification.ReferenceNumber,
		updated.User.Identification.FrontImageUrl,
		updated.User.Identification.BackImageUrl,
		updated.User.Identification.SelfieImageUrl,
		updated.User.Identification.Id,
	); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlVerificationRepository) FindById(id string) (*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	getRequest := `
        SELECT
            *
        FROM
            IdentityVerification
        WHERE
            id = ?
    `
	request := new(models.IdentityVerificationModel)
	if err := s.db.GetContext(ctx, request, getRequest, id); err != nil {
		return nil, err
	}

	getUser := `
        SELECT
            u.*
        FROM
            User u
            JOIN IdentityVerification iv ON iv.userId = u.id
        WHERE
            iv.id = ?
    `
	if err := s.db.GetContext(ctx, &request.User, getUser, id); err != nil {
		return nil, err
	}

	if err := s.getUserAddress(&request.User.Address, request.User.Id); err != nil {
		return nil, err
	}

	if err := s.getUserIdentification(&request.User.Identification, request.User.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return request, nil
}

func (s *MysqlVerificationRepository) FindByUserId(userId string) (*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getRequest := `
        SELECT
            *
        FROM
            IdentityVerification
        WHERE
            userId = ?
    `
	request := new(models.IdentityVerificationModel)
	if err := s.db.GetContext(ctx, request, getRequest, userId); err != nil {
		return nil, err
	}

	getUser := `
        SELECT
            u.*
        FROM
            User u
            JOIN IdentityVerification iv ON iv.userId = u.id
        WHERE
            iv.id = ?
    `
	if err := s.db.GetContext(ctx, &request.User, getUser, request.Id); err != nil {
		return nil, err
	}

	if err := s.getUserAddress(&request.User.Address, request.User.Id); err != nil {
		return nil, err
	}

	if err := s.getUserIdentification(&request.User.Identification, request.User.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return request, nil
}

func (s *MysqlVerificationRepository) GetAll(status string) ([]*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	ids := make([]string, 0)
	getIds := "SELECT id FROM IdentityVerification WHERE status = ?"
	if err := s.db.SelectContext(ctx, &ids, getIds, status); err != nil {
		return nil, err
	}

	requests := make([]*models.IdentityVerificationModel, 0)
	for _, id := range ids {
		request, err := s.FindById(id)
		if err != nil {
			return nil, err
		}

		requests = append(requests, request)
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return requests, nil
}

func (s *MysqlVerificationRepository) getUserAddress(dest *models.AddressModel, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getAddress := `
        SELECT
            a.*
        FROM
            Address a
            JOIN UserAddress ua ON ua.addressId = a.id
        WHERE
            ua.userId = ?
    `
	if err := s.db.GetContext(ctx, dest, getAddress, userId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlVerificationRepository) getUserIdentification(dest *models.IdentificationModel, userId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getIdentification := `
        SELECT
            i.*
        FROM
            Identification i
            JOIN UserIdentification ui ON ui.identificationId = i.id
        WHERE
            ui.userId = ?
    `
	if err := s.db.GetContext(ctx, dest, getIdentification, userId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
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
            u.verified = 1,
            u.verifiedAt = ?
        WHERE
            iv.id = ?
    `
	timestamp := utils.CurrentTimeStamp()
	if _, err := tx.ExecContext(ctx, updateUser, timestamp, id); err != nil {
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
            rejectionNote = ?
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
