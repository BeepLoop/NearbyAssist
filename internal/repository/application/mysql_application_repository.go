package application_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlApplicationRepository struct {
	db *sqlx.DB
}

func NewMysqlApplicationRepository(db *sqlx.DB) *MysqlApplicationRepository {
	return &MysqlApplicationRepository{
		db: db,
	}
}

func (s *MysqlApplicationRepository) Create(data *models.ApplicationModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	duplicatePendingQuery := `
        SELECT COUNT(id)
        FROM Application
        WHERE applicantId = ? AND expertiseId = ? AND status = 'pending'
    `
	dupePendingResult := 0
	if err := s.db.GetContext(ctx, &dupePendingResult, duplicatePendingQuery, data.ApplicantId, data.ExpertiseId); err != nil {
		return "", err
	}
	if dupePendingResult != 0 {
		return "", errors.New("Duplicate entry")
	}

	alreadyApprovedExpertiseCheck := `
        SELECT COUNT(id)
        FROM Application
        WHERE applicantId = ? AND expertiseId = ? AND status = 'approved'
    `
	alreadyApprovedResult := 0
	if err := s.db.GetContext(ctx, &alreadyApprovedResult, alreadyApprovedExpertiseCheck, data.ApplicantId, data.ExpertiseId); err != nil {
		return "", err
	}
	if alreadyApprovedResult != 0 {
		return "", errors.New("Already approved")
	}

	createQuery := `
        INSERT INTO
            Application (id, applicantId, expertiseId, supportingDocumentUrl, policeClearanceUrl)
        VALUES
            (:id, :applicantId, :expertiseId, :supportingDocumentUrl, :policeClearanceUrl)
    `
	if _, err := s.db.NamedExecContext(ctx, createQuery, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlApplicationRepository) FindById(id string) (*models.ApplicationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	application := new(models.ApplicationModel)
	query := `
        SELECT 
            a.id,
            a.applicantId,
            a.expertiseId,
            a.supportingDocumentUrl,
            a.policeClearanceUrl,
            a.status,
            u.name AS applicantName,
            e.title AS expertise
        FROM 
            Application a
            JOIN User u ON u.id = a.applicantId
            JOIN Expertise e ON e.id = a.expertiseId
        WHERE
            a.id = ?
    `
	if err := s.db.GetContext(ctx, application, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return application, nil
}

func (s *MysqlApplicationRepository) NewProof(data *models.ApplicationProofModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	query := `
        INSERT INTO
            ApplicationProof (id, applicationId, applicantId, url)
        VALUES
            (:id, :applicationId, :applicantId, :url)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlApplicationRepository) NewPoliceClearance(data *models.PoliceClearanceModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	query := `
        INSERT INTO
            PoliceClearance (id, applicationId, applicantId, url)
        VALUES
            (:id, :applicationId, :applicantId, :url)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlApplicationRepository) GetAll(status string) ([]*models.ApplicationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	applications := make([]*models.ApplicationModel, 0)

	query := `
        SELECT
            a.id,
            a.applicantId,
            a.expertiseId,
            a.createdAt,
            a.supportingDocumentUrl,
            a.policeClearanceUrl,
            a.status,
            u.name AS applicantName
        FROM
            Application a
            JOIN User u ON u.id = a.applicantId
        WHERE
            a.status = ?
    `
	if err := s.db.SelectContext(ctx, &applications, query, status); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return applications, nil
}

func (s *MysqlApplicationRepository) AcceptRequest(applicationId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	updateApplicationStatusQuery := `
        UPDATE
            Application
        SET
            status = 'approved'
        WHERE
            id = ?
    `
	if _, err := tx.ExecContext(ctx, updateApplicationStatusQuery, applicationId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	checkIfAlreadyVendorQuery := `
        SELECT 
            COUNT(v.id)
        FROM 
            Vendor v
            JOIN Application a ON a.applicantId = v.vendorId
        WHERE 
            a.id = ?
    `
	alreadyVendorResult := 0
	if err := tx.GetContext(ctx, &alreadyVendorResult, checkIfAlreadyVendorQuery, applicationId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}
	if alreadyVendorResult == 0 {
		// Applicant is not yet a vendor
		vendorRowId := utils.GenerateId()

		makeUserVendorQuery := `
        INSERT INTO
            Vendor (id, vendorId)
        SELECT
            ?, u.id
        FROM
            User u 
            JOIN Application a ON a.applicantId = u.id
        WHERE
            a.id = ?
    `
		if _, err := tx.ExecContext(ctx, makeUserVendorQuery, vendorRowId, applicationId); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	addVendorExpertiseQuery := `
        INSERT INTO
            VendorExpertise (vendorId, expertiseId)
        SELECT
            v.vendorId, a.expertiseId
        FROM 
            Application a
            JOIN Vendor v ON v.vendorId = a.applicantId
        WHERE
            a.id = ?
    `
	if _, err := tx.ExecContext(ctx, addVendorExpertiseQuery, applicationId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

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

func (s *MysqlApplicationRepository) RejectRequest(applicationId, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Application
        SET 
            status = 'rejected',
            reason = ?
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, reason, applicationId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
