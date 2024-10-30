package transaction_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"time"

	"github.com/jmoiron/sqlx"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type MysqlTransactionRepository struct {
	db *sqlx.DB
}

func NewMysqlTransactionRepository(db *sqlx.DB) *MysqlTransactionRepository {
	return &MysqlTransactionRepository{
		db: db,
	}
}

func (s *MysqlTransactionRepository) Create(data *models.TransactionModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := gonanoid.New(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := `
        INSERT INTO
            Transaction (id, vendorId, clientId, serviceId, start, end, confirmCode)
        VALUES
            (:id, :vendorId, :clientId, :serviceId, :start, :end, :confirmCode)
    `

	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlTransactionRepository) FindById(id string) (*models.TransactionModel, error) {
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

func (s *MysqlTransactionRepository) GetSummary(transactionId string) (*response.TransactionSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	summary := new(response.TransactionSummary)
	query := `
        SELECT
            t.id,
            t.createdAt,
            vendor.Name AS vendor,
            client.Name AS client,
            s.title AS serviceTitle,
            t.price,
            t.startDate,
            t.endDate,
            vendor.Email AS vendorEmail,
            client.Email AS clientEmail
        FROM
            Transaction t
            JOIN User client ON t.clientId = client.id
            JOIN User vendor ON t.vendorId = vendor.id
            JOIN Service s ON t.serviceId = s.id
        WHERE
            id = ?
    `
	if err := s.db.GetContext(ctx, summary, query, transactionId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return summary, nil
}

func (s *MysqlTransactionRepository) GetAll() ([]*models.TransactionModel, error) {
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

func (s *MysqlTransactionRepository) GetMyTransactions(id string) ([]*models.TransactionModel, error) {
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

func (s *MysqlTransactionRepository) GetOngoing(id string) ([]*models.TransactionModel, error) {
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

func (s *MysqlTransactionRepository) GetHistory(id string) ([]*models.TransactionModel, error) {
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

func (s *MysqlTransactionRepository) MarkComplete(transactionId string) error {
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
