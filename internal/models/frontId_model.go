package models

import (
	"context"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type FrontIdModel struct {
	Model
	UpdateableModel
	Url string `json:"url" db:"url"`
}

func NewFrontIdModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *FrontIdModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &FrontIdModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func NewFrontIdModelWithImageUrl(url string, idGenerator id_generator.IdGenerator, conn *sqlx.DB) *FrontIdModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &FrontIdModel{
		Model: Model{Id: id, Conn: conn},
		Url:   url,
	}
}

func (f *FrontIdModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO FrontId (id, url) VALUES (:id, :url)"

	if _, err := f.Conn.NamedExecContext(ctx, query, f); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return f.Id, nil
}
