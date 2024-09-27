package models

// import (
// 	"context"
// 	"nearbyassist/internal/id_generator"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type ServiceTagModel struct {
	Model
	UpdateableModel
	ServiceId string `json:"serviceId" db:"serviceId"`
	TagId     string `json:"tagId" db:"tagId"`

	// Additional fields for joining tag title
	Title string `json:"title" db:"title"`
}

// func NewServiceTagModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *ServiceTagModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &ServiceTagModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func (s *ServiceTagModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	insertTag := `
//         INSERT INTO
//             ServiceTag (id, serviceId, tagId)
//         VALUES
//             (:id, :serviceId, (SELECT id FROM Tag WHERE title = :title))
//     `
// 	if _, err := s.Conn.NamedExecContext(ctx, insertTag, s); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return s.Id, nil
// }
//
// func (s *ServiceTagModel) Delete() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "DELETE FROM ServiceTag WHERE id = ?"
// 	if _, err := s.Conn.ExecContext(ctx, query, s, s.Id); err != nil {
// 		return err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return context.DeadlineExceeded
// 	}
//
// 	return nil
// }
