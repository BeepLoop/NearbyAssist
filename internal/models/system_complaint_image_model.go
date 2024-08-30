package models

import (
	"context"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type SystemComplaintData struct {
	ComplaintId string
	Url         string
}

type SystemComplaintImageModel struct {
	Model
	UpdateableModel
	ComplaintId string `json:"complaintId" db:"complaintId"`
	Url         string `json:"url" db:"url"`
}

func NewSystemComplaintImageModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *SystemComplaintImageModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &SystemComplaintImageModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func NewSystemComplaintImageWithData(data SystemComplaintData, idGenerator id_generator.IdGenerator, conn *sqlx.DB) *SystemComplaintImageModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &SystemComplaintImageModel{
		Model:       Model{Id: id, Conn: conn},
		ComplaintId: data.ComplaintId,
		Url:         data.Url,
	}
}

func (s *SystemComplaintImageModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            SystemComplaintImage (id, complaintId, url)
        VALUES
            (:id, :complaintId, :url)
    `

	if _, err := s.Conn.NamedExecContext(ctx, query, s); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return s.Id, nil
}
