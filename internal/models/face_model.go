package models

import (
	"context"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type FaceModel struct {
	Model
	UpdateableModel
	Url string `json:"url" db:"url"`
}

func NewFaceModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *FaceModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &FaceModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func NewFaceModelWithImageUrl(url string, idGenerator id_generator.IdGenerator, conn *sqlx.DB) *FaceModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &FaceModel{
		Model: Model{Id: id, Conn: conn},
		Url:   url,
	}
}

func (f *FaceModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Face (id, url) VALUES (:id, :url)"

	if _, err := f.Conn.NamedExecContext(ctx, query, f); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return f.Id, nil
}
