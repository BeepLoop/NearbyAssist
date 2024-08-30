package models

import (
	"context"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type BackIdModel struct {
	Model
	UpdateableModel
	Url string `json:"url" db:"url"`
}

func NewBackIdModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *BackIdModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &BackIdModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func NewBackIdModelWithImageUrl(url string, idGenerator id_generator.IdGenerator, conn *sqlx.DB) *BackIdModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &BackIdModel{
		Model: Model{Id: id, Conn: conn},
		Url:   url,
	}
}

func (b *BackIdModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO BackId (id, url) VALUES (:id, :url)"

	if _, err := b.Conn.NamedExecContext(ctx, query, b); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return b.Id, nil
}
