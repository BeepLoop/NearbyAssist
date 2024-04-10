package models

import (
	"context"
	"nearbyassist/internal/db"
	"time"
)

type NamedTransactionModel struct {
	Model
	UpdateableModel
	Vendor  string `json:"vendor" db:"vendor"`
	Client  string `json:"client" db:"client"`
	Service string `json:"service" db:"service"`
	Status  string `json:"status" db:"status"`
}

type TransactionModel struct {
	Model
	UpdateableModel
	VendorId  int    `json:"vendorId" db:"vendorId" validate:"required"`
	ClientId  int    `json:"clientId" db:"clientId" validate:"required"`
	ServiceId int    `json:"serviceId" db:"serviceId" validate:"required"`
	Status    string `json:"status" db:"status"`
}

func NewTransactionModel(db *db.DB) *TransactionModel {
	return &TransactionModel{
		Model: Model{Db: db},
	}
}

func (t *TransactionModel) Create() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            Transaction (vendorId, clientId, serviceId)
        VALUES
            (:vendorId, :clientId, :serviceId)
    `

	res, err := t.Db.Conn.NamedExecContext(ctx, query, t)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}

func (t *TransactionModel) Update(id int) error {
	return nil
}

func (t *TransactionModel) Delete(id int) error {
	return nil
}

func (t *TransactionModel) Count(status string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT 
            COUNT(*)
        FROM
            Transaction
    `

	switch status {
	case "ongoing":
		query += " WHERE status = 'ongoing'"
	case "done":
		query += " WHERE status = 'done'"
	case "cancelled":
		query += " WHERE status = 'cancelled'"
	}

	count := 0
	err := t.Db.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (t *TransactionModel) GetClientOngoing(clientId int) ([]NamedTransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            uVendor.name as vendor,
            uClient.name as client,
            s.title as service,
            t.status
        FROM 
            Transaction t
        LEFT JOIN User uVendor ON uVendor.id = t.vendorId
        LEFT JOIN User uClient ON uClient.id = t.clientId
        LEFT JOIN Service s ON s.id = t.serviceId
        WHERE
            status = 'ongoing' AND t.clientId= ?
    `

	transactions := make([]NamedTransactionModel, 0)
	err := t.Db.Conn.SelectContext(ctx, &transactions, query, clientId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (t *TransactionModel) GetVendorOngoing(clientId int) ([]NamedTransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            uVendor.name as vendor,
            uClient.name as client,
            s.title as service,
            t.status
        FROM 
            Transaction t
        LEFT JOIN User uVendor ON uVendor.id = t.vendorId
        LEFT JOIN User uClient ON uClient.id = t.clientId
        LEFT JOIN Service s ON s.id = t.serviceId
        WHERE
            status = 'ongoing' AND t.vendorId= ?
    `

	transactions := make([]NamedTransactionModel, 0)
	err := t.Db.Conn.SelectContext(ctx, &transactions, query, clientId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (t *TransactionModel) ClientHistory(clientId int) ([]NamedTransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            uVendor.name as vendor,
            uClient.name as client,
            s.title as service,
            t.status,
            t.createdAt,
            t.updatedAt
        FROM 
            Transaction t
        LEFT JOIN User uVendor ON uVendor.id = t.vendorId
        LEFT JOIN User uClient ON uClient.id = t.clientId
        LEFT JOIN Service s ON s.id = t.serviceId
        WHERE
            status = 'done' OR status = 'cancelled' AND t.clientId = ?
    `

	transactions := make([]NamedTransactionModel, 0)
	err := t.Db.Conn.SelectContext(ctx, &transactions, query, clientId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (t *TransactionModel) VendorHistory(vendorId int) ([]NamedTransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            uVendor.name as vendor,
            uClient.name as client,
            s.title as service,
            t.status
        FROM 
            Transaction t
        LEFT JOIN User uVendor ON uVendor.id = t.vendorId
        LEFT JOIN User uClient ON uClient.id = t.clientId
        LEFT JOIN Service s ON s.id = t.serviceId
        WHERE
            status = 'done' OR status = 'cancelled' AND t.vendorId = ?
    `

	transactions := make([]NamedTransactionModel, 0)
	err := t.Db.Conn.SelectContext(ctx, &transactions, query, vendorId)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}
