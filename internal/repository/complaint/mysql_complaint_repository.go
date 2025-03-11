package complaint_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlComplaintRepository struct {
	db *sqlx.DB
}

func NewMysqlComplaintRepository(db *sqlx.DB) *MysqlComplaintRepository {
	return &MysqlComplaintRepository{
		db: db,
	}
}

func (s *MysqlComplaintRepository) CreateBugReport(data *models.BugReportModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	insertComplaint := `
        INSERT INTO 
            BugReport (id, title, detail)
        VALUES
            (:id, :title, :detail)
    `
	if _, err := tx.NamedExecContext(ctx, insertComplaint, data); err != nil {
		return "", err
	}

	insertImage := `
        INSERT INTO 
            BugReportImage (id, complaintId, url)
        VALUES
            (?, ?, ?)
    `
	for _, url := range data.Images {
		imageId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, insertImage, imageId, data.Id, url); err != nil {
			return "", nil
		}
	}

	if err := tx.Commit(); err != nil {
		return "", nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlComplaintRepository) GetAll(limit, offset int) ([]*models.ComplaintModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "SELECT * FROM BugReport ORDER BY createdAt DESC LIMIT ? OFFSET ?"

	complaints := make([]*models.ComplaintModel, 0)
	if err := s.db.SelectContext(ctx, &complaints, query, limit, offset); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return complaints, nil
}

func (s *MysqlComplaintRepository) FindById(id string) (*models.ComplaintModel, error) {
	// TODO: Implement this method
	return nil, nil
}
