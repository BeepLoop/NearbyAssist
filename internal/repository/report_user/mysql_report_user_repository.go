package report_user_repo

import (
	"context"
	"errors"
	"nearbyassist/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

	if id, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	insertQuery := `
        INSERT INTO 
            ReportedUser (id, userId, reason, detail)
        VALUES
            (:id, :userId, :reason, :detail)
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
		imageId, err := gonanoid.New()
		if err != nil {
			return "", err
		}

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

	getUsersQuery := "SELECT * FROM ReportedUser ORDER BY createdAt DESC LIMIT ? OFFSET ?"
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

	getUserQuery := "SELECT * FROM ReportedUser WHERE id = ?"

	reportedUser := new(models.ReportedUserModel)
	if err := s.db.GetContext(ctx, &reportedUser, getUserQuery, id); err != nil {
		return nil, err
	}

	if reportedUser.Id == "" {
		return nil, errors.New("not found")
	}

	getImagesQuery := "SELECT url FROM ReportedUserImage WHERE reportId = ?"
	images := make([]string, 0)
	if err := s.db.SelectContext(ctx, &images, getImagesQuery, reportedUser.Id); err != nil {
		return nil, err
	}
	reportedUser.Images = images

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reportedUser, nil
}

func (s *MysqlReportUserRepository) FindByUserId(userId string) (*models.ReportedUserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	getUserQuery := "SELECT * FROM ReportedUser WHERE userId = ?"

	reportedUser := new(models.ReportedUserModel)
	if err := s.db.GetContext(ctx, &reportedUser, getUserQuery, userId); err != nil {
		return nil, err
	}

	if reportedUser.Id == "" {
		return nil, errors.New("not found")
	}

	getImagesQuery := "SELECT url FROM ReportedUserImage WHERE reportId = ?"
	images := make([]string, 0)
	if err := s.db.SelectContext(ctx, &images, getImagesQuery, reportedUser.Id); err != nil {
		return nil, err
	}
	reportedUser.Images = images

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return reportedUser, nil
}
