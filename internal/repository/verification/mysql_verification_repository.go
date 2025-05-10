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

func (s *MysqlVerificationRepository) CreateLink(userId string) (string, error) {
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

func (s *MysqlVerificationRepository) UpdateUserInfo(input *models.IdentityVerificationModel, newTimestamp bool) error {
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
	if _, err := tx.ExecContext(ctx, updateUser, input.User.Name, input.User.Phone, input.User.Id); err != nil {
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
		input.User.Address.Address,
		input.User.Address.Latitude,
		input.User.Address.Longitude,
		input.User.Address.Id,
	); err != nil {
		return err
	}

	// Update Identification
	createIdentification := `
        INSERT INTO
            Identification (id, type, referenceNumber, frontImageUrl, backImageUrl, selfieImageUrl)
        VALUES
            (:id, :type, :referenceNumber, :frontImageUrl, :backImageUrl, :selfieImageUrl)
    `
	input.User.Identification.Id = utils.GenerateId()
	if _, err := tx.NamedExecContext(ctx, createIdentification, input.User.Identification); err != nil {
		return err
	}

	createIdentificationRelation := `
        INSERT INTO
            UserIdentification (userId, identificationId)
        VALUES
            (?, ?)
    `
	if _, err := tx.ExecContext(ctx, createIdentificationRelation, input.User.Id, input.User.Identification.Id); err != nil {
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
        ORDER BY
            createdAt DESC
        LIMIT 1
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
            u.verifiedAt = CURRENT_TIMESTAMP()
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
