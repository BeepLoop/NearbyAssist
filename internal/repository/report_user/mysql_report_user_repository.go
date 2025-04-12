package report_user_repo

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlReportUserRepository struct {
	db *sqlx.DB
}

func NewMysqlReportUserRepository(db *sqlx.DB) *MysqlReportUserRepository {
	return &MysqlReportUserRepository{
		db: db,
	}
}

func (s *MysqlReportUserRepository) Create(data *models.UserReportModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	insertQuery := `
        INSERT INTO 
            UserReport (id, reporterUserId, reportedUserId, category, bookingId, reason, detail)
        VALUES
            (?, ?, ?, ?, ?, ?, ?)
    `
	if _, err := tx.ExecContext(
		ctx,
		insertQuery,
		data.Id,
		data.ReporterUserId,
		data.ReportedUserId,
		data.Category,
		utils.StringOrNil(data.BookingIdInput),
		data.Reason,
		data.Detail,
	); err != nil {
		return "", err
	}

	insertImage := `
        INSERT INTO 
            UserReportImage (id, reportId, url)
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

func (s *MysqlReportUserRepository) GetAllWithStatus(status string, limit, offset int) ([]*models.UserReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`
        SELECT
            *
        FROM
            UserReport
        WHERE
            status = '%s'
        ORDER BY
            createdAt DESC
        LIMIT %d OFFSET %d
    `, status, limit, offset)

	reports := make([]*models.UserReportModel, 0)
	if err := s.db.SelectContext(ctx, &reports, query); err != nil {
		return nil, err
	}

	for _, report := range reports {
		if res, err := s.GetImages(report.Id); err != nil {
			return nil, err
		} else {
			report.Images = res
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reports, nil
}

func (s *MysqlReportUserRepository) FindById(id string) (*models.UserReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getUserQuery := `
        SELECT
            *
        FROM
            UserReport
        WHERE
            id = ?
    `

	report := new(models.UserReportModel)
	if err := s.db.GetContext(ctx, report, getUserQuery, id); err != nil {
		return nil, err
	}

	if report.Id == "" {
		return nil, errors.New("not found")
	}

	if res, err := s.GetImages(report.Id); err != nil {
		return nil, err
	} else {
		report.Images = res
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return report, nil
}

func (s *MysqlReportUserRepository) GetAllReportedIs(userId string) ([]*models.UserReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getReportsQuery := `
        SELECT
            *
        FROM
            UserReport
        WHERE
            reportedUserId = ?
        ORDER BY
            createdAt DESC
    `

	reports := make([]*models.UserReportModel, 0)
	if err := s.db.SelectContext(ctx, &reports, getReportsQuery, userId); err != nil {
		return nil, err
	}

	for _, report := range reports {
		report.Images = make([]string, 0)

		if res, err := s.GetImages(report.Id); err != nil {
			return nil, err
		} else {
			report.Images = res
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reports, nil
}

func (s *MysqlReportUserRepository) GetAllReportedBy(userId string) ([]*models.UserReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getReportsQuery := `
        SELECT
            *
        FROM
            UserReport
        WHERE
            reporterUserId = ?
        ORDER BY
            createdAt DESC
    `

	reports := make([]*models.UserReportModel, 0)
	if err := s.db.SelectContext(ctx, &reports, getReportsQuery, userId); err != nil {
		return nil, err
	}

	for _, report := range reports {
		report.Images = make([]string, 0)

		if res, err := s.GetImages(report.Id); err != nil {
			return nil, err
		} else {
			report.Images = res
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reports, nil
}

func (s *MysqlReportUserRepository) GetImages(reportId string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            url
        FROM
            UserReportImage
        WHERE
            reportId = ?
    `

	images := make([]string, 0)
	if err := s.db.SelectContext(ctx, &images, query, reportId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return images, nil
}

func (s *MysqlReportUserRepository) UpdateStatus(reportId, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE UserReport SET status = ? WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, status, reportId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlReportUserRepository) Close(reportId, action, adminId, note string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            UserReport
        SET
            status = ?, adminId = ?, adminNote = ?
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, action, adminId, note, reportId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
