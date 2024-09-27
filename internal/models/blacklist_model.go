package models

// import (
// 	"context"
// 	"nearbyassist/internal/id_generator"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type BlacklistModel struct {
	Model
	UpdateableModel
	Token string `json:"token" db:"token"`
}

// func NewBlacklistModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *BlacklistModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &BlacklistModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func NewBlacklistModelWithToken(token string, idGenerator id_generator.IdGenerator, conn *sqlx.DB) *BlacklistModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &BlacklistModel{
// 		Model: Model{Id: id, Conn: conn},
// 		Token: token,
// 	}
// }
//
// func (b *BlacklistModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "INSERT INTO Blacklist (id, token) VALUES (:id, :token)"
// 	if _, err := b.Conn.NamedExecContext(ctx, query, b); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return b.Id, nil
// }
//
// func (b *BlacklistModel) FindByToken(token string) (*BlacklistModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, token FROM Blacklist WHERE token = ?"
//
// 	if err := b.Conn.GetContext(ctx, b, query, b.Token); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return b, nil
// }
