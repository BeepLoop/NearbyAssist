package application_repo

import (
	"context"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

	if id, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := `
        INSERT INTO
            Application (id, applicantId, expertiseId, supportingDocumentUrl, policeClearanceUrl)
        VALUES
            (:id, :applicantId, :expertiseId, :supportingDocumentUrl, :policeClearanceUrl)
    `
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlApplicationRepository) FindApplication(id string) (*models.ApplicationModel, error) {
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

	if id, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

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

	if id, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

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

	query := "SELECT * FROM Application WHERE status = ?"
	if err := s.db.SelectContext(ctx, &applications, query, status); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return applications, nil
}
