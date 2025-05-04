package application_repo

import (
	"context"
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
	createQuery := `
        INSERT INTO
            Application (id, applicantId, expertiseId, supportingDocument, policeClearance)
        VALUES
            (:id, :applicantId, :expertiseId, :supportingDocument, :policeClearance)
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
            a.createdAt,
            a.applicantId,
            a.expertiseId,
            a.supportingDocument,
            a.policeClearance,
            i.url AS supportingDocumentUrl,
            p.url AS policeClearanceUrl,
            a.status,
            u.name AS applicantName,
            e.title AS expertise
        FROM 
            Application a
            JOIN User u ON u.id = a.applicantId
            JOIN Expertise e ON e.id = a.expertiseId
            JOIN SupportingImage i ON i.id = a.supportingDocument
            JOIN PoliceClearance p ON p.id = a.policeClearance
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

func (s *MysqlApplicationRepository) NewSupportingImage(data *models.SupportingImageModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	query := `
        INSERT INTO
            SupportingImage (id, url)
        VALUES
            (:id, :url)
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
            PoliceClearance (id, url)
        VALUES
            (:id, :url)
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
            a.supportingDocument,
            a.policeClearance,
            i.url AS supportingDocumentUrl,
            p.url AS policeClearanceUrl,
            a.status,
            u.name AS applicantName
        FROM
            Application a
            JOIN User u ON u.id = a.applicantId
            JOIN SupportingImage i ON i.id = a.supportingDocument
            JOIN PoliceClearance p ON p.id = a.policeClearance
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
            status = 'approved',
            updatedAt = CURRENT_TIMESTAMP()
        WHERE
            id = ?
    `
	if _, err := tx.ExecContext(ctx, updateApplicationStatusQuery, applicationId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	addUserExpertiseQuery := `
        INSERT INTO
            UserExpertise (userId, expertiseId, supportingImage)
        SELECT
            applicantId, expertiseId, supportingDocument
        FROM 
            Application
        WHERE
            id = ?
    `
	if _, err := tx.ExecContext(ctx, addUserExpertiseQuery, applicationId); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	createVendorQuery := `
        INSERT IGNORE
            Vendor (vendorId)
        SELECT DISTINCT
            applicantId
        FROM
            Application
        WHERE
            id = ?
    `
	if _, err := tx.ExecContext(ctx, createVendorQuery, applicationId); err != nil {
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
            rejectionReason = ?
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

func (s *MysqlApplicationRepository) HasPendingApplication(userId, expertiseId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT EXISTS
            (SELECT 1 FROM Application WHERE applicantId = ? AND expertiseId = ? AND status = 'pending')
        AS has_duplicate_pending
    `
	hasPending := false
	if err := s.db.GetContext(ctx, &hasPending, query, userId, expertiseId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return hasPending, nil
}
