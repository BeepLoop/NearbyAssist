package report_user_repo

import (
	"context"
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

func (s *MysqlReportUserRepository) Create(data *models.UserReportModel) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	input := struct {
		ReporterUserId string  `db:"reporterUserId"`
		ReportedUserId string  `db:"reportedUserId"`
		Category       string  `db:"category"`
		BookingId      *string `db:"bookingId"`
		Reason         string  `db:"reason"`
		Detail         string  `db:"detail"`
	}{
		ReporterUserId: data.ReporterUserId,
		ReportedUserId: data.ReportedUserId,
		Category:       string(data.Category),
		BookingId:      utils.StringOrNil(data.BookingIdInput),
		Reason:         data.Reason,
		Detail:         data.Detail,
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	insertQuery := `
        INSERT INTO 
            UserReport (reporterUserId, reportedUserId, category, bookingId, reason, detail)
        VALUES
            (:reporterUserId, :reportedUserId, :category, :bookingId, :reason, :detail)
    `
	res, err := tx.NamedExecContext(ctx, insertQuery, input)
	if err != nil {
		return err
	}

	insertId, err := res.LastInsertId()
	if err != nil {
		return err
	}

	insertImage := `
        INSERT INTO 
            UserReportImage (id, reportId, url)
        VALUES
            (?, ?, ?)
    `
	for _, url := range data.Images {
		imageId := utils.GenerateId()
		if _, err := tx.ExecContext(ctx, insertImage, imageId, insertId, url); err != nil {
			return nil
		}
	}

	if err := tx.Commit(); err != nil {
		return nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlReportUserRepository) GetAllWithStatus(status string, limit, offset int) ([]*models.UserReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            *
        FROM
            UserReport
        WHERE
            status = ?
        ORDER BY
            createdAt DESC
        LIMIT ? OFFSET ?
    `

	reports := make([]*models.UserReportModel, 0)
	if err := s.db.SelectContext(ctx, &reports, query, status, limit, offset); err != nil {
		return nil, err
	}

	for _, report := range reports {
		if res, err := s.getImages(report.Id); err != nil {
			return nil, err
		} else {
			report.Images = res
		}

		if reported, err := s.getUser(report.ReportedUserId); err != nil {
			return nil, err
		} else {
			report.Reported = *reported
		}

		if reporter, err := s.getUser(report.ReporterUserId); err != nil {
			return nil, err
		} else {
			report.Reporter = *reporter
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reports, nil
}

func (s *MysqlReportUserRepository) FindById(id int) (*models.UserReportModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getUserQuery := "SELECT * FROM UserReport WHERE id = ?"

	report := new(models.UserReportModel)
	if err := s.db.GetContext(ctx, report, getUserQuery, id); err != nil {
		return nil, err
	}

	if res, err := s.getImages(report.Id); err != nil {
		return nil, err
	} else {
		report.Images = res
	}

	if reported, err := s.getUser(report.ReportedUserId); err != nil {
		return nil, err
	} else {
		report.Reported = *reported
	}

	if reporter, err := s.getUser(report.ReporterUserId); err != nil {
		return nil, err
	} else {
		report.Reporter = *reporter
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

		if res, err := s.getImages(report.Id); err != nil {
			return nil, err
		} else {
			report.Images = res
		}

		if reported, err := s.getUser(report.ReportedUserId); err != nil {
			return nil, err
		} else {
			report.Reported = *reported
		}

		if reporter, err := s.getUser(report.ReporterUserId); err != nil {
			return nil, err
		} else {
			report.Reporter = *reporter
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

		if res, err := s.getImages(report.Id); err != nil {
			return nil, err
		} else {
			report.Images = res
		}

		if reported, err := s.getUser(report.ReportedUserId); err != nil {
			return nil, err
		} else {
			report.Reported = *reported
		}

		if reporter, err := s.getUser(report.ReporterUserId); err != nil {
			return nil, err
		} else {
			report.Reporter = *reporter
		}
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reports, nil
}

func (s *MysqlReportUserRepository) UpdateStatus(reportId int, status string) error {
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

func (s *MysqlReportUserRepository) Close(reportId int, action, adminId, note string) error {
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

func (s *MysqlReportUserRepository) getImages(reportId int) ([]string, error) {
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

func (s *MysqlReportUserRepository) getUser(userId string) (*models.UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            id, name, email
        FROM
            User
        WHERE
            id = ?
    `

	user := new(models.UserModel)
	if err := s.db.GetContext(ctx, user, query, userId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return user, nil
}
