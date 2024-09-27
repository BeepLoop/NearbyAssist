package models

// import (
// 	"context"
// 	"nearbyassist/internal/id_generator"
// 	"time"
//
// 	"github.com/jmoiron/sqlx"
// )

type TransactionStatusFilter string

const (
	TRANSACTION_STATUS_ONGOING   TransactionStatusFilter = "ongoing"
	TRANSACTION_STATUS_DONE      TransactionStatusFilter = "done"
	TRANSACTION_STATUS_CANCELLED TransactionStatusFilter = "cancelled"
)

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId   string                  `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId   string                  `json:"clientId" db:"clientId" validate:"required"`
	ServiceId  string                  `json:"serviceId" db:"serviceId" validate:"required"`
	Start      string                  `json:"start" db:"start" validate:"required"`
	End        string                  `json:"end" db:"end" validate:"required"`
	Status     TransactionStatusFilter `json:"status" db:"status"`
	IsReviewed bool                    `json:"isReviewed" db:"isReviewed"`
	IsReported bool                    `json:"isReported" db:"isReported"`

	// Additional fields for joins
	Vendor string `json:"vendor" db:"vendor"` // Vendor name
	Client string `json:"client" db:"client"` // Client name
}

// func NewTransactionModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *TransactionModel {
// 	id, err := idGenerator.Generate()
// 	if err != nil {
// 		return nil
// 	}
//
// 	return &TransactionModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func NewTransactionModelWithId(id string, conn *sqlx.DB) *TransactionModel {
// 	return &TransactionModel{
// 		Model: Model{Id: id, Conn: conn},
// 	}
// }
//
// func (t *TransactionModel) DecryptVendorName(decryptFunc func(string) (string, error)) (*TransactionModel, error) {
// 	if decrypted, err := decryptFunc(t.Vendor); err != nil {
// 		return nil, err
// 	} else {
// 		t.Vendor = decrypted
// 	}
// 	return t, nil
// }
//
// func (t *TransactionModel) DecryptClientName(decryptFunc func(string) (string, error)) (*TransactionModel, error) {
// 	if decrypted, err := decryptFunc(t.Client); err != nil {
// 		return nil, err
// 	} else {
// 		t.Client = decrypted
// 	}
// 	return t, nil
// }
//
// func (t *TransactionModel) Create() (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := `
//         INSERT INTO
//             Transaction (id, vendorId, clientId, serviceId, start, end)
//         VALUES
//             (:id, :vendorId, :clientId, :serviceId, :start, :end)
//     `
//
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
// func (t *TransactionModel) Count(filter TransactionStatusFilter) (int, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT COUNT(*) FROM Transaction"
//
// 	switch filter {
// 	case TRANSACTION_STATUS_ONGOING:
// 		query += " WHERE status = 'ongoing'"
// 	case TRANSACTION_STATUS_DONE:
// 		query += " WHERE status = 'done'"
// 	case TRANSACTION_STATUS_CANCELLED:
// 		query += " WHERE status = 'cancelled'"
// 	}
//
// 	count := 0
// 	err := t.Conn.GetContext(ctx, &count, query)
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
// func (t *TransactionModel) FindById(id string) (*TransactionModel, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "SELECT * FROM Transaction WHERE id = ?"
//
// 	if err := t.Conn.GetContext(ctx, t, query, id); err != nil {
// 		return nil, err
// 	}
//
// 	if ctx.Err() == context.DeadlineExceeded {
// 		return nil, context.DeadlineExceeded
// 	}
//
// 	return t, nil
// }
//
// func (t *TransactionModel) MarkComplete(id string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
// 	defer cancel()
//
// 	query := "UPDATE Transaction SET status = 'done' WHERE id = ?"
// 	if _, err := t.Conn.ExecContext(ctx, query, id); err != nil {
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
// func (t *TransactionModel) IsReviewable() bool {
// 	return t.Status == TRANSACTION_STATUS_DONE && t.IsReviewed == false
// }
