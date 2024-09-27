package models

// import (
// 	"context"
// 	"nearbyassist/internal/id_generator"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type ComplaintModel struct {
	Model
	UpdateableModel
	Code    int    `json:"code" db:"code"`
	Title   string `json:"title" db:"title"`
	Content string `json:"content" db:"content"`
}

// func NewComplaintModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *ComplaintModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &ComplaintModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func (c *ComplaintModel) Count() (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT COUNT(*) FROM Complaint"
//
// 	count := 0
// 	if err := c.Conn.GetContext(ctx, &count, query); err != nil {
// 		return 0, nil
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return 0, context.DeadlineExceeded
// 	}
//
// 	return count, nil
// }
