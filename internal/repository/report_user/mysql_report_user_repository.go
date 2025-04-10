package report_user_repo

import (
	"context"
	"errors"
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

func (s *MysqlReportUserRepository) Create(data *models.ReportedUserModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	data.Id = utils.GenerateId()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	insertQuery := `
        INSERT INTO 
            ReportedUser (id, reportedBy, userId, reason, detail)
        VALUES
            (:id, :reportedBy, :userId, :reason, :detail)
    `
	if _, err := tx.NamedExecContext(ctx, insertQuery, data); err != nil {
		return "", err
	}

	insertImage := `
        INSERT INTO 
            ReportedUserImage (id, reportId, url)
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

func (s *MysqlReportUserRepository) GetAll(limit, offset int) ([]*models.ReportedUserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getUsersQuery := `
        SELECT
            *
        FROM
            ReportedUser
        WHERE
            completedAt IS NULL
        ORDER BY
            createdAt DESC
        LIMIT ? OFFSET ?
    `
	reportedUsers := make([]*models.ReportedUserModel, 0)
	if err := s.db.SelectContext(ctx, &reportedUsers, getUsersQuery, limit, offset); err != nil {
		return nil, err
	}

	getImagesQuery := "SELECT url FROM ReportedUserImage WHERE reportId = ?"
	for _, user := range reportedUsers {
		images := make([]string, 0)
		if err := s.db.SelectContext(ctx, &images, getImagesQuery, user.Id); err != nil {
			return nil, err
		}

		user.Images = images
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reportedUsers, nil
}

func (s *MysqlReportUserRepository) FindById(id string) (*models.ReportedUserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getUserQuery := `
        SELECT
            ru.id,
            ru.reportedBy,
            ru.userId,
            ru.reason,
            ru.detail,
            ru.createdAt,
            ru.completedAt
        FROM
            ReportedUser ru
            JOIN User u ON u.id = ru.userId
        WHERE
            ru.id = ?
    `

	report := new(models.ReportedUserModel)
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

func (s *MysqlReportUserRepository) FindAllWithUserId(userId string) ([]*models.ReportedUserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getReportsQuery := `
        SELECT
            ru.id,
            ru.reportedBy,
            ru.userId,
            ru.reason,
            ru.detail,
            ru.createdAt,
            ru.completedAt
        FROM
            ReportedUser ru
        WHERE
            ru.userId = ?
        ORDER BY
            ru.createdAt DESC
    `

	reports := make([]*models.ReportedUserModel, 0)
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
            ReportedUserImage
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

func (s *MysqlReportUserRepository) CloseReport(reportId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE ReportedUser SET completedAt = CURRENT_TIMESTAMP() WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, reportId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
