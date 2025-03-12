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

func (s *MysqlBugReportRepository) Create(data *models.BugReportModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	insertComplaint := `
        INSERT INTO 
            BugReport (title, detail)
        VALUES
            (:title, :detail)
    `
	res, err := tx.NamedExecContext(ctx, insertComplaint, data)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	bugId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	insertImage := `
        INSERT INTO 
            BugReportImage (id, reportId, url)
        VALUES
            (?, ?, ?)
    `
	for _, url := range data.Images {
		imageId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, insertImage, imageId, bugId, url); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
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

func (s *MysqlBugReportRepository) GetAll(limit, offset int) ([]*models.BugReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getReportsQuery := `
        SELECT
            *
        FROM
            BugReport
        WHERE
            completedAt IS NULL
        ORDER BY
            createdAt DESC
        LIMIT ? OFFSET ?
    `
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

func (s *MysqlBugReportRepository) CompleteBug(bugId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE BugReport SET completedAt = CURRENT_TIMESTAMP() WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, bugId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlBugReportRepository) FindById(id int) (*models.BugReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getReportQuery := "SELECT * FROM BugReport WHERE id = ?"

	bugReport := new(models.BugReportModel)
	if err := s.db.GetContext(ctx, &bugReport, getReportQuery, id); err != nil {
		return nil, err
	}

	if bugReport.Id == 0 {
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
