package models

//
// import (
// 	"context"
// 	"nearbyassist/internal/id_generator"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type ReviewModel struct {
	Model
	UpdateableModel
	ServiceId string `json:"serviceId" db:"serviceId" validate:"required"`
	Rating    int    `json:"rating" db:"rating" validate:"required"`

	// Additional fields for creating a review
	TransactionId string `json:"transactionId" db:"transactionId" validate:"required"`
}

//
// func NewReviewModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *ReviewModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &ReviewModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func (r *ReviewModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	tx, err := r.Conn.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	insertReview := `
//         INSERT INTO
//             Review (id, serviceId, rating)
//         VALUES
//             (:id, :serviceId, :rating)
//     `
//
// 	if _, err := tx.NamedExecContext(ctx, insertReview, r); err != nil {
// 		if err := tx.Rollback(); err != nil {
// 			return "", err
// 		}
//
// 		return "", err
// 	}
//
// 	updateReviewedFlag := `
//         UPDATE
//             Transaction
//         SET
//             isReviewed = 1
//         WHERE
//             id = ?
//     `
// 	_, err = tx.ExecContext(ctx, updateReviewedFlag, r.TransactionId)
// 	if err != nil {
// 		if err := tx.Rollback(); err != nil {
// 			return "", err
// 		}
//
// 		return "", err
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		if err := tx.Rollback(); err != nil {
// 			return "", err
// 		}
//
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return r.Id, nil
// }
//
// func (r *ReviewModel) FindById(id string) (*ReviewModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, serviceId, rating FROM Review WHERE id = ?"
//
// 	if err := r.Conn.GetContext(ctx, r, query, id); err != nil {
// 		return nil, err
// 	}
//
// 	return r, nil
// }
//
// func (r *ReviewModel) FindAllByServiceId(id string) ([]ReviewModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT id, serviceId, rating FROM Review WHERE serviceId = ?"
//
// 	reviews := make([]ReviewModel, 0)
// 	if err := r.Conn.SelectContext(ctx, &reviews, query, id); err != nil {
// 		return nil, err
// 	}
//
// 	return reviews, nil
// }
