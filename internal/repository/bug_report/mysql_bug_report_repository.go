package bug_report_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlBugReportRepository struct {
	db *sqlx.DB
}

func NewMysqlBugReportRepository(db *sqlx.DB) *MysqlBugReportRepository {
	return &MysqlBugReportRepository{
		db: db,
	}
}

func (s *MysqlBugReportRepository) Create(data *models.BugReportModel) (string, error) {
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
            BugReportImage (id, reportId, url)
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

func (s *MysqlBugReportRepository) GetAll(limit, offset int) ([]*models.BugReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getReportsQuery := "SELECT * FROM BugReport ORDER BY createdAt DESC LIMIT ? OFFSET ?"
	reports := make([]*models.BugReportModel, 0)
	if err := s.db.SelectContext(ctx, &reports, getReportsQuery, limit, offset); err != nil {
		return nil, err
	}

	getImagesQuery := "SELECT url FROM BugReportImage WHERE reportId = ?"
	for _, report := range reports {
		images := make([]string, 0)
		if err := s.db.SelectContext(ctx, &images, getImagesQuery, report.Id); err != nil {
			return nil, err
		}

		report.Images = images
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reports, nil
}

func (s *MysqlBugReportRepository) FindById(id string) (*models.BugReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getReportQuery := "SELECT * FROM ReportedUser WHERE id = ?"

	bugReport := new(models.BugReportModel)
	if err := s.db.GetContext(ctx, &bugReport, getReportQuery, id); err != nil {
		return nil, err
	}

	if bugReport.Id == "" {
		return nil, errors.New("not found")
	}

	getImagesQuery := "SELECT url FROM BugReportImage WHERE reportId = ?"
	images := make([]string, 0)
	if err := s.db.SelectContext(ctx, &images, getImagesQuery, bugReport.Id); err != nil {
		return nil, err
	}
	bugReport.Images = images

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return bugReport, nil
}
