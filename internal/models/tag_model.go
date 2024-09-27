package models

// import (
// 	"context"
// 	"nearbyassist/internal/id_generator"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type TagModel struct {
	Model
	UpdateableModel
	Title string `json:"title" db:"title"`
}

//
// func NewTagModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *TagModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &TagModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func (t *TagModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "INSERT INTO Tag (id, title) VALUES (:id, :title)"
// 	if _, err := t.Conn.NamedExecContext(ctx, query, t); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return t.Id, nil
// }
//
// func (t *TagModel) FindAll() ([]TagModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, title FROM Tag"
//
// 	tags := make([]TagModel, 0)
// 	if err := t.Conn.SelectContext(ctx, &tags, query); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return tags, nil
// }
//
// func (t *TagModel) Delete() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "DELETE FROM Tag WHERE id = ?"
//
// 	if _, err := t.Conn.ExecContext(ctx, query, t.Id); err != nil {
// 		return err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return context.DeadlineExceeded
// 	}
//
// 	return nil
// }
