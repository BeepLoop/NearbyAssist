package transaction

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlTransactionStore struct {
	db *sqlx.DB
}

func NewMysqlTransactionStore(db *sqlx.DB) *MysqlTransactionStore {
	return &MysqlTransactionStore{
		db: db,
	}
}

func (s *MysqlTransactionStore) Create(data *models.TransactionModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := `
        INSERT INTO
            Transaction (id, vendorId, clientId, serviceId, start, end)
        VALUES
            (:id, :vendorId, :clientId, :serviceId, :start, :end)
    `

	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlTransactionStore) FindById(id string) (*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transaction := new(models.TransactionModel)

	query := "SELECT * FROM Transaction WHERE id = ?"
	if err := s.db.GetContext(ctx, transaction, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transaction, nil
}

func (s *MysqlTransactionStore) GetAll() ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactions := make([]*models.TransactionModel, 0)

	query := "SELECT * FROM Transaction"
	if err := s.db.SelectContext(ctx, &transactions, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionStore) GetMyTransactions(id string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactions := make([]*models.TransactionModel, 0)

	query := "SELECT * FROM Transaction WHERE clientId = ? OR vendorId = ?"
	if err := s.db.SelectContext(ctx, &transactions, query, id, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionStore) GetOngoing(id string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            uVendor.name as vendor,
            uClient.name as client,
            t.status
        FROM
            Transaction t
            LEFT JOIN User uVendor ON uVendor.id = t.vendorId
            LEFT JOIN User uClient ON uClient.id = t.clientId
            LEFT JOIN Service s ON s.id = t.serviceId
        WHERE status = 'ongoing' AND (t.clientId = ? OR t.vendorId = ?)
    `

	transactions := make([]*models.TransactionModel, 0)
	err := s.db.SelectContext(ctx, &transactions, query, id, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionStore) GetHistory(id string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            t.id,
            uVendor.name as vendor,
            uClient.name as client,
            t.createdAt as createdAt,
            t.status
        FROM
            Transaction t
            LEFT JOIN User uVendor ON uVendor.id = t.vendorId
            LEFT JOIN User uClient ON uClient.id = t.clientId
            LEFT JOIN Service s ON s.id = t.serviceId
        WHERE status = 'done' OR status = 'cancelled' AND (t.clientId = ? OR t.vendorId = ?)
    `

	transactions := make([]*models.TransactionModel, 0)
	err := s.db.SelectContext(ctx, &transactions, query, id, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionStore) MarkComplete(transactionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "UPDATE Transaction SET status = 'done' WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, transactionId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}
