package transaction_repo

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
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

	checkDuplicate := `
        SELECT CASE
            WHEN EXISTS
                (
                    SELECT
                        1
                    FROM
                        Transaction
                    WHERE
                        (clientId = ? AND serviceId = ?)
                        AND (status = 'confirmed' OR status = 'pending')
                )
            THEN 1
            ELSE 0
        END AS duplicate_booking
    `
	alreadyBooked := false
	if err := tx.GetContext(ctx, &alreadyBooked, checkDuplicate, data.ClientId, data.ServiceId); err != nil {
		return "", err
	}

	if alreadyBooked {
		return "", errors.New("You already have an confirmed or pending transaction for this service")
	}

	query := `
        INSERT INTO
            Transaction (id, vendorId, clientId, serviceId, cost)
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
            t.createdAt,
            t.scheduledAt,
            t.cancelReason,
            t.updatedAt,
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

	if isReviewed, err := s.IsReviewed(transaction.Id); err != nil {
		return nil, err
	} else {
		transaction.IsReviewed = isReviewed
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

func (s *MysqlTransactionRepository) GetConfirmedTransactionsOfVendor(vendorId string) ([]*models.TransactionModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT
            *
        FROM
            Transaction
        WHERE
            vendorId = ? AND status = 'confirmed'
    `

	transactions := make([]*models.TransactionModel, 0)
	if err := s.db.SelectContext(ctx, &transactions, query, vendorId); err != nil {
		return nil, err
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
		if isReviewed, err := s.IsReviewed(transaction.Id); err != nil {
			return nil, err
		} else {
			transaction.IsReviewed = isReviewed
		}

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
		if isReviewed, err := s.IsReviewed(transaction.Id); err != nil {
			return nil, err
		} else {
			transaction.IsReviewed = isReviewed
		}

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
		if isReviewed, err := s.IsReviewed(transaction.Id); err != nil {
			return nil, err
		} else {
			transaction.IsReviewed = isReviewed
		}

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
		if isReviewed, err := s.IsReviewed(transaction.Id); err != nil {
			return nil, err
		} else {
			transaction.IsReviewed = isReviewed
		}

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
		if isReviewed, err := s.IsReviewed(transaction.Id); err != nil {
			return nil, err
		} else {
			transaction.IsReviewed = isReviewed
		}

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
            uVendor.name AS vendor,
            uClient.name AS client
        FROM 
            Transaction t
            JOIN User uVendor ON uVendor.id = t.vendorId
            JOIN User uClient ON uClient.id = t.clientId
        WHERE
            t.clientId = ? AND t.status = 'done' AND NOT EXISTS (
                SELECT 1 FROM Review r WHERE r.transactionId = t.id
            )
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
		if isReviewed, err := s.IsReviewed(transaction.Id); err != nil {
			return nil, err
		} else {
			transaction.IsReviewed = isReviewed
		}

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

func (s *MysqlTransactionRepository) Cancel(transactionId, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE 
            Transaction 
        SET 
            cancelReason = ?, status = 'cancelled'
        WHERE 
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, reason, transactionId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlTransactionRepository) Accept(transactionId, schedule string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Transaction
        SET
            scheduledAt = ?, status = 'confirmed'
        WHERE
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, schedule, transactionId); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlTransactionRepository) Reject(transactionId, reason string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        UPDATE
            Transaction
        SET
            cancelReason = ?,
            status = 'rejected'
        WHERE 
            id = ?
    `
	if _, err := s.db.ExecContext(ctx, query, reason, transactionId); err != nil {
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

func (s *MysqlTransactionRepository) IsReviewed(transactionId string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
        SELECT CASE
            WHEN EXISTS (SELECT 1 FROM Review WHERE transactionId = ?)
            THEN 1
            ELSE 0
        END AS is_reviewed
    `

	isReviewed := false
	if err := s.db.GetContext(ctx, &isReviewed, query, transactionId); err != nil {
		return false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false, context.DeadlineExceeded
	}

	return isReviewed, nil
}
