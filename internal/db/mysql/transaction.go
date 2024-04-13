package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) CountTransaction(status models.TransactionStatus) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM Transaction"

	switch status {
	case models.TRANSACTION_STATUS_ONGOING:
		query += " WHERE status = 'ongoing'"
	case models.TRANSACTION_STATUS_DONE:
		query += " WHERE status = 'done'"
	case models.TRANSACTION_STATUS_CANCELLED:
		query += " WHERE status = 'cancelled'"
	}

	count := -1
	err := m.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return count, nil
}

func (m *Mysql) CreateTransaction(transaction *models.TransactionModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO
            Transaction (vendorId, clientId, serviceId)
        VALUES
            (:vendorId, :clientId, :serviceId)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, transaction)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return -1, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) FindAllOngoingTransaction(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error) {
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
        WHERE status = 'ongoing' 
    `

	switch filter {
	case models.FILTER_CLIENT:
		query += "AND t.clientId = ?"
	case models.FILTER_VENDOR:
		query += "AND t.vendorId = ?"
	default:
		query += "AND t.clientId = ?"
	}

	transactions := make([]models.DetailedTransactionModel, 0)
	err := m.Conn.SelectContext(ctx, &transactions, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (m *Mysql) GetTransactionHistory(id int, filter models.TransactionFilter) ([]models.DetailedTransactionModel, error) {
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
        WHERE status = 'done' OR status = 'cancelled'
    `

	switch filter {
	case models.FILTER_CLIENT:
		query += " AND t.clientId = ?"
	case models.FILTER_VENDOR:
		query += " AND t.vendorId = ?"
	default:
		query += " AND t.clientId = ?"
	}

	transactions := make([]models.DetailedTransactionModel, 0)
	err := m.Conn.SelectContext(ctx, &transactions, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}
