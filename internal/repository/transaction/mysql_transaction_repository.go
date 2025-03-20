package transaction_repo

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
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

	data.Id = utils.GenerateId()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", err
	}

	count := 0
	checkDuplicate := "SELECT count(id) FROM Transaction WHERE (clientId = ? AND serviceId = ?) AND (status = 'confirmed' OR status = 'pending')"
	if err := tx.GetContext(ctx, &count, checkDuplicate, data.ClientId, data.ServiceId); err != nil {
		return "", err
	}

	if count > 0 {
		return "", errors.New("You already have an confirmed or pending transaction for this service")
	}

	query := `
        INSERT INTO
            Transaction (id, vendorId, clientId, serviceId, cost )
        VALUES
            (:id, :vendorId, :clientId, :serviceId, :cost)
    `

	if _, err := tx.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	insertTransactionExtras := `
        INSERT INTO
            TransactionExtra (transactionId, extraId)
        VALUES 
            (?, ?)
    `
	for _, extra := range data.Extras {
		if _, err := tx.ExecContext(ctx, insertTransactionExtras, data.Id, extra.Id); err != nil {
			fmt.Println("extra: ", extra)
			fmt.Println("error: ", err.Error())
			return "", err
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return "", err
		}

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

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.id = ?
    `

	transaction := new(models.TransactionModel)
	if err := s.db.GetContext(ctx, transaction, transactionQuery, id); err != nil {
		fmt.Println("error get transaction: ", err.Error())
		return nil, err
	}

	serviceQuery := "SELECT * FROM Service WHERE id = ?"

	service := new(models.ServiceModel)
	if err := s.db.GetContext(ctx, service, serviceQuery, transaction.ServiceId); err != nil {
		fmt.Println("error get service: ", err.Error())
		return nil, err
	}
	transaction.Service = service

	extrasQuery := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            TransactionExtra te
            JOIN Extra e ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	extras := make([]*models.ExtraModel, 0)
	if err := s.db.SelectContext(ctx, &extras, extrasQuery, transaction.Id); err != nil {
		fmt.Println("error get extra: ", err.Error())
		return nil, err
	}
	transaction.Extras = extras

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
            t.cost,
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
            t.id = ?
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

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? OR t.vendorId = ?
        ORDER BY
            t.updatedAt DESC
    `

	transactions := make([]*models.TransactionModel, 0)
	if err := s.db.SelectContext(ctx, &transactions, transactionQuery, id, id); err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate
        FROM
            Service s
            JOIN Transaction t ON s.id = t.serviceId
        WHERE
            t.id = ?
        ORDER BY
            t.updatedAt DESC
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN TransactionExtra te ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	for _, transaction := range transactions {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionRepository) GetTransactionSent(id string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND (t.status = 'pending' OR t.status = 'confirmed')
        ORDER BY
            t.updatedAt DESC
    `

	transactions := make([]*models.TransactionModel, 0)
	if err := s.db.SelectContext(ctx, &transactions, transactionQuery, id); err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Transaction t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN TransactionExtra te ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	for _, transaction := range transactions {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionRepository) GetTransactionReceived(id string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? AND t.status = 'pending'
        ORDER BY
            t.updatedAt DESC
    `

	transactions := make([]*models.TransactionModel, 0)
	if err := s.db.SelectContext(ctx, &transactions, transactionQuery, id); err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Transaction t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN TransactionExtra te ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	for _, transaction := range transactions {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionRepository) GetRecent(userId string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? OR t.clientId = ?
        ORDER BY
            t.updatedAt DESC
        LIMIT
            10
    `

	transactions := make([]*models.TransactionModel, 0)
	if err := s.db.SelectContext(ctx, &transactions, transactionQuery, userId, userId); err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Transaction t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN TransactionExtra te ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	for _, transaction := range transactions {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionRepository) GetConfirmed(id string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.vendorId = ? AND t.status = 'confirmed'
        ORDER BY
            t.updatedAt DESC
    `

	transactions := make([]*models.TransactionModel, 0)
	err := s.db.SelectContext(ctx, &transactions, transactionQuery, id)
	if err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Transaction t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN TransactionExtra te ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	for _, transaction := range transactions {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionRepository) GetHistory(id string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            (t.vendorId = ? OR t.clientId = ?) AND (t.status = 'done' OR t.status = 'cancelled')
        ORDER BY
            t.updatedAt DESC
    `

	transactions := make([]*models.TransactionModel, 0)
	err := s.db.SelectContext(ctx, &transactions, transactionQuery, id, id)
	if err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Transaction t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN TransactionExtra te ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	for _, transaction := range transactions {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionRepository) GetReviewableTransactions(userId string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	transactionQuery := `
        SELECT
            t.id,
            t.vendorId,
            t.clientId,
            t.serviceId,
            t.status,
            t.cost,
            t.isReviewed,
            t.isReported,
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND t.status = 'done' AND t.isReviewed = 0
        ORDER BY
            t.updatedAt DESC
    `

	transactions := make([]*models.TransactionModel, 0)
	err := s.db.SelectContext(ctx, &transactions, transactionQuery, userId)
	if err != nil {
		return nil, err
	}

	getService := `
        SELECT
            s.id,
            s.vendorId,
            s.title,
            s.description,
            s.rate,
            s.latitude,
            s.longitude
        FROM
            Service s
            JOIN Transaction t ON s.id = t.serviceId
        WHERE
            t.id = ?
    `

	getExtras := `
        SELECT
            e.id,
            e.title,
            e.description,
            e.price
        FROM 
            Extra e
            JOIN TransactionExtra te ON e.id = te.extraId
        WHERE
            te.transactionId = ?
    `

	for _, transaction := range transactions {
		extras := make([]*models.ExtraModel, 0)
		if err := s.db.SelectContext(ctx, &extras, getExtras, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Extras = extras

		service := new(models.ServiceModel)
		if err := s.db.GetContext(ctx, service, getService, transaction.Id); err != nil {
			return nil, err
		}
		transaction.Service = service
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (s *MysqlTransactionRepository) Cancel(transactionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Transaction SET status = 'cancelled' WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, transactionId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlTransactionRepository) Accept(transactionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Transaction SET status = 'confirmed' WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, transactionId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlTransactionRepository) Reject(transactionId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE Transaction SET status = 'rejected' WHERE id = ?"
	if _, err := s.db.ExecContext(ctx, query, transactionId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
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
