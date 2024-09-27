package models

// import (
// 	"context"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type VendorComplaintModel struct {
	Model
	UpdateableModel
	VendorId string `json:"vendorId" db:"vendorId"`
	Title    string `json:"title" db:"title"`
	Content  string `json:"content" db:"content"`
}

//
// func NewVendorComplaintModel(conn *sqlx.DB) *VendorComplaintModel {
// 	return &VendorComplaintModel{
// 		Model: Model{Conn: conn},
// 	}
// }
//
// func (v *VendorComplaintModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO
//             Complaint (id, vendorId, code, title, content)
//         VALUES
//             (:id, :vendorId, :code, :title, :content)
//     `
//
// 	if _, err := v.Conn.NamedExecContext(ctx, query, v); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return v.Id, nil
// }
//
// func (v *VendorComplaintModel) Count() (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT COUNT(*) FROM VendorComplaint"
//
// 	count := 0
// 	err := v.Conn.GetContext(ctx, &count, query)
// 	if err != nil {
// 		return 0, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return 0, context.DeadlineExceeded
// 	}
//
// 	return count, nil
// }
//
// func (v *VendorComplaintModel) FindById(id string) (*VendorComplaintModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         SELECT
//             id, vendorId, title, content
//         FROM
//             VendorComplaint
//         WHERE
//             id = ?
//     `
//
// 	if err := v.Conn.GetContext(ctx, v, query, id); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return v, nil
// }
