package models

import (
	"context"
	"fmt"
	"nearbyassist/internal/id_generator"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type IdentityVerificationModel struct {
	Model
	UpdateableModel
	UserId   string `json:"userId" db:"userId" validate:"required"`
	Name     string `json:"name" db:"name" validate:"required"`
	Address  string `json:"address" db:"address" validate:"required"`
	IdType   string `json:"idType" db:"idType" validate:"required"`
	IdNumber string `json:"idNumber" db:"idNumber" validate:"required"`
	FrontId  string `json:"frontId" db:"frontId" validate:"required"`
	BackId   string `json:"backId" db:"backId" validate:"required"`
	Face     string `json:"face" db:"face" validate:"required"`
}

func NewIdentityVerificationModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *IdentityVerificationModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &IdentityVerificationModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func (i *IdentityVerificationModel) EncryptName(encryptFunc func(string) (string, error)) (*IdentityVerificationModel, error) {
	if encrypted, err := encryptFunc(i.Name); err != nil {
		return nil, err
	} else {
		i.Name = encrypted
	}
	return i, nil
}

func (i *IdentityVerificationModel) EncryptAddress(encryptFunc func(string) (string, error)) (*IdentityVerificationModel, error) {
	if encrypted, err := encryptFunc(i.Address); err != nil {
		return nil, err
	} else {
		i.Address = encrypted
	}
	return i, nil
}

func (i *IdentityVerificationModel) EncryptIdNumber(encryptFunc func(string) (string, error)) (*IdentityVerificationModel, error) {
	if encrypted, err := encryptFunc(i.IdNumber); err != nil {
		return nil, err
	} else {
		i.IdNumber = encrypted
	}
	return i, nil
}

func (i *IdentityVerificationModel) DecryptName(decryptFunc func(string) (string, error)) (*IdentityVerificationModel, error) {
	if decrypted, err := decryptFunc(i.Name); err != nil {
		return nil, err
	} else {
		i.Name = decrypted
	}
	return i, nil
}

func (i *IdentityVerificationModel) DecryptAddress(decryptFunc func(string) (string, error)) (*IdentityVerificationModel, error) {
	if decrypted, err := decryptFunc(i.Address); err != nil {
		return nil, err
	} else {
		i.Address = decrypted
	}
	return i, nil
}

func (i *IdentityVerificationModel) DecryptIdNumber(decryptFunc func(string) (string, error)) (*IdentityVerificationModel, error) {
	if decrypted, err := decryptFunc(i.IdNumber); err != nil {
		return nil, err
	} else {
		i.IdNumber = decrypted
	}
	return i, nil
}

func (i *IdentityVerificationModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO IdentityVerification (id, name, address, idType, idNumber, frontId, backId, face)
        VALUES ( :id, :name, :address, :idType, :idNumber, :frontId, :backId, :face)
    `

	if _, err := i.Conn.NamedExecContext(ctx, query, i); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return i.Id, nil
}

func (i *IdentityVerificationModel) FindById(id string) (*IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, address, idType, idNumber, frontId, backId, face FROM IdentityVerification WHERE id = ?"

	if err := i.Conn.GetContext(ctx, i, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return i, nil
}

func (i *IdentityVerificationModel) FindAll(params map[string]string) ([]IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, user as userId, createdAt FROM IdentityVerification"

	// Order by id and createdAt (deterministic order)
	query += " ORDER BY id, createdAt"

	// Paginate using limit and offset
	if page, ok := params["page"]; ok {
		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}

		pageSize := DEFAULT_LIMIT
		if limit, ok := params["limit"]; ok {
			if size, err := strconv.Atoi(limit); err != nil {
				return nil, err
			} else {
				pageSize = size
			}
		}

		offset := (pageNumber - 1) * pageSize
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", pageSize, offset)
	} else {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", DEFAULT_LIMIT, DEFAULT_OFFSET)
	}

	requests := make([]IdentityVerificationModel, 0)
	if err := i.Conn.SelectContext(ctx, &requests, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return requests, nil
}
