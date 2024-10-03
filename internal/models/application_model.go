package models

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"nearbyassist/internal/id_generator"
// 	"strconv"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type ApplicationStatusFilter string

const (
	APPLICATION_STATUS_ALL      ApplicationStatusFilter = "all"
	APPLICATION_STATUS_PENDING  ApplicationStatusFilter = "pending"
	APPLICATION_STATUS_APPROVED ApplicationStatusFilter = "approved"
	APPLICATION_STATUS_REJECTED ApplicationStatusFilter = "rejected"
)

type ApplicationModel struct {
	Model
	UpdateableModel
	GeoSpatialModel
	ApplicantId string                  `db:"applicantId" validate:"required"`
	Job         string                  `db:"job" validate:"required"`
	Status      ApplicationStatusFilter `db:"status"`
}

//
// func NewApplicationModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *ApplicationModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &ApplicationModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func (a *ApplicationModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO
//             Application (id, applicantId, job, latitude, longitude)
//         VALUES
//             (:id, :applicantId, :job, :latitude, :longitude)
//     `
//
// 	if _, err := a.Conn.NamedExecContext(ctx, query, a); err != nil {
// 		return "", err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return "", context.DeadlineExceeded
// 	}
//
// 	return a.Id, nil
// }
//
// func (a *ApplicationModel) Count(filter ApplicationStatusFilter) (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT COUNT(*) FROM Application"
//
// 	switch filter {
// 	case APPLICATION_STATUS_PENDING:
// 		query += " WHERE status = 'pending'"
// 	case APPLICATION_STATUS_APPROVED:
// 		query += " WHERE status = 'approved'"
// 	case APPLICATION_STATUS_REJECTED:
// 		query += " WHERE status = 'rejected'"
// 	}
//
// 	count := 0
// 	err := a.Conn.GetContext(ctx, &count, query)
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
// func (a *ApplicationModel) FindById(id string) (*ApplicationModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         SELECT
//             id, applicantId, job, status, latitude, longitude
//         FROM
//             Application
//         WHERE
//             id = ?
//     `
//
// 	if err := a.Conn.GetContext(ctx, a, query, id); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return a, nil
// }
//
// func (a *ApplicationModel) FindAll(params map[string]string) ([]ApplicationModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	// Initial query
// 	query := "SELECT id, applicantId, status, createdAt FROM Application"
//
// 	// Apply filters
// 	if filter, ok := params["filter"]; ok {
// 		switch filter {
// 		case "all":
// 			query += " WHERE status = 'approved' OR status = 'rejected' OR status = 'pending'"
// 		case "approved":
// 			query += " WHERE status = 'approved'"
// 		case "rejected":
// 			query += " WHERE status = 'rejected'"
// 		case "pending":
// 			query += " WHERE status = 'pending'"
// 		default:
// 			return nil, errors.New("Invalid parameter found")
// 		}
// 	} else {
// 		query += " WHERE status = 'approved' OR status = 'rejected' OR status = 'pending'"
// 	}
//
// 	// Order by id and createdAt (deterministic order)
// 	query += " ORDER BY id, createdAt"
//
// 	// Paginate using limit and offset
// 	if page, ok := params["page"]; ok {
// 		pageNumber, err := strconv.Atoi(page)
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		pageSize := DEFAULT_LIMIT
// 		if limit, ok := params["limit"]; ok {
// 			if size, err := strconv.Atoi(limit); err != nil {
// 				return nil, err
// 			} else {
// 				pageSize = size
// 			}
// 		}
//
// 		offset := (pageNumber - 1) * pageSize
// 		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
// 	} else {
// 		query += fmt.Sprintf(" LIMIT %d OFFSET %d", DEFAULT_LIMIT, DEFAULT_OFFSET)
// 	}
//
// 	// Execute query
// 	applications := make([]ApplicationModel, 0)
// 	if err := a.Conn.SelectContext(ctx, &applications, query); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return applications, nil
// }
//
// func (a *ApplicationModel) Approve(id string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	tx, err := a.Conn.BeginTxx(ctx, nil)
// 	if err != nil {
// 		return err
// 	}
//
// 	updateStatus := "UPDATE Application SET status = 'approved' WHERE id = ?"
//
// 	if _, err := tx.ExecContext(ctx, updateStatus, id); err != nil {
// 		rollbackErr := tx.Rollback()
// 		if rollbackErr != nil {
// 			return errors.New("Failed to approve application and rollback transaction")
// 		}
//
// 		return err
// 	}
//
// 	promoteVendor := fmt.Sprintf(`
//         INSERT INTO
//             Vendor (vendorId, job)
//         VALUES
//             (
//                 (SELECT applicantId FROM Application WHERE id = %s),
//                 (SELECT job FROM Application WHERE id = %s)
//             )
//     `, id, id)
//
// 	if _, err := tx.ExecContext(ctx, promoteVendor); err != nil {
// 		rollbackErr := tx.Rollback()
// 		if rollbackErr != nil {
// 			return errors.New("Failed to promote applicant to vendor and rollback transaction")
// 		}
//
// 		return err
// 	}
//
// 	if err := tx.Commit(); err != nil {
// 		rollbackErr := tx.Rollback()
// 		if rollbackErr != nil {
// 			return errors.New("Failed to commit transaction and rollback")
// 		}
//
// 		return err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return context.DeadlineExceeded
// 	}
//
// 	return nil
// }
//
// func (a *ApplicationModel) Reject(id string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "UPDATE Application SET status = 'rejected' WHERE id = ?"
//
// 	if _, err := a.Conn.ExecContext(ctx, query, id); err != nil {
// 		return err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return context.DeadlineExceeded
// 	}
//
// 	return nil
// }
