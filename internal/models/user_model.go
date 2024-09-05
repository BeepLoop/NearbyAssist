package models

import (
	"context"
	"nearbyassist/internal/id_generator"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserStatusFilter string

const (
	USER_STATUS_VERIFIED   UserStatusFilter = "verified"
	USER_STATUS_UNVERIFIED UserStatusFilter = "unverified"
	USER_STATUS_ALL        UserStatusFilter = "all"
)

type UserModel struct {
	Model
	UpdateableModel
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	ImageUrl string `json:"imageUrl" db:"imageUrl"`
	Verified bool   `json:"verified" db:"verified"`
	Hash     string `json:"hash" db:"hash"`
}

func NewUserModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *UserModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &UserModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func NewUserModelWithId(id string, conn *sqlx.DB) *UserModel {
	return &UserModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func (u *UserModel) DecryptName(decryptFunc func(string) (string, error)) (*UserModel, error) {
	if decrypted, err := decryptFunc(u.Name); err != nil {
		return nil, err
	} else {
		u.Name = decrypted
	}
	return u, nil
}

func (u *UserModel) DecryptEmail(decryptFunc func(string) (string, error)) (*UserModel, error) {
	if decrypted, err := decryptFunc(u.Email); err != nil {
		return nil, err
	} else {
		u.Email = decrypted
	}
	return u, nil
}

func (u *UserModel) HashEmail(plainEmail string, hashFunc func([]byte) (string, error)) (*UserModel, error) {
	if hash, err := hashFunc([]byte(plainEmail)); err != nil {
		return nil, err
	} else {
		u.Hash = hash
	}
	return u, nil
}

func (u *UserModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO User (id, name, email, imageUrl, emailHash) VALUES (:id, :name, :email, :imageUrl, :hash)"
	if _, err := u.Conn.NamedExecContext(ctx, query, u); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return u.Id, nil
}

func (u *UserModel) Count(filter UserStatusFilter) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM User"

	switch filter {
	case USER_STATUS_VERIFIED:
		query += " WHERE verified = 1"
	case USER_STATUS_UNVERIFIED:
		query += " WHERE verified = 0"
	case USER_STATUS_ALL:
	}

	count := 0
	err := u.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (u *UserModel) FindById(id string) (*UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, email, imageUrl, verified FROM User WHERE id = ?"

	err := u.Conn.GetContext(ctx, u, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return u, nil
}

func (u *UserModel) FindByEmailHash() (*UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, email, imageUrl, verified FROM User WHERE emailHash = ?"

	if err := u.Conn.GetContext(ctx, u, query, u.Hash); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return u, nil
}

func (u *UserModel) IsVerified() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT verified FROM User WHERE id = ?"

	var verified bool
	if err := u.Conn.GetContext(ctx, &verified, query, u.Id); err != nil {
		return false
	}

	if ctx.Err() == context.DeadlineExceeded {
		return false
	}

	return verified
}

// User Conversations
func (u *UserModel) GetConversations() ([]*UserModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT DISTINCT
            u.id,
            u.name,
            u.imageUrl
        FROM
            User u
        JOIN
            Message m ON u.id = m.sender OR u.id = m.receiver
        WHERE
            u.id <> ?
    `

	conversations := make([]*UserModel, 0)
	if err := u.Conn.SelectContext(ctx, &conversations, query, u.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return conversations, nil
}

func (u *UserModel) GetMessages(otherUserId string) ([]MessageModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        SELECT
            id, sender, receiver, content, createdAt
        FROM
            Message
        WHERE
            sender = ? AND receiver = ?
        OR
            sender = ? AND receiver = ?
        ORDER BY
            createdAt
    `

	messages := make([]MessageModel, 0)
	if err := u.Conn.SelectContext(
		ctx,
		&messages,
		query,
		u.Id,
		otherUserId,
		otherUserId,
		u.Id,
	); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return messages, nil
}

// User Transactions
func (u *UserModel) GetTransactions() ([]*TransactionModel, error) {
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
        WHERE
            t.clientId = ? OR t.vendorId = ?
    `

	transactions := make([]*TransactionModel, 0)
	if err := u.Conn.SelectContext(ctx, &transactions, query, u.Id, u.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (u *UserModel) GetOngoingTransactions() ([]TransactionModel, error) {
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

	transactions := make([]TransactionModel, 0)
	err := u.Conn.SelectContext(ctx, &transactions, query, u.Id, u.Id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}

func (u *UserModel) GetTransactionHistory() ([]TransactionModel, error) {
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

	transactions := make([]TransactionModel, 0)
	err := u.Conn.SelectContext(ctx, &transactions, query, u.Id, u.Id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return transactions, nil
}
