package models

import (
	"context"
	"fmt"
	"nearbyassist/internal/id_generator"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type SystemComplaintModel struct {
	Model
	UpdateableModel
	Title  string `json:"title" db:"title" validate:"required"`
	Detail string `json:"detail" db:"detail" validate:"required"`
}

func NewSystemComplaintModel(idGenerator id_generator.IdGenerator, conn *sqlx.DB) *SystemComplaintModel {
	id, err := idGenerator.Generate()
	if err != nil {
		return nil
	}

	return &SystemComplaintModel{
		Model: Model{Id: id, Conn: conn},
	}
}

func (s *SystemComplaintModel) EncryptTitle(encryptFunc func(string) (string, error)) (*SystemComplaintModel, error) {
	encrypted, err := encryptFunc(s.Title)
	if err != nil {
		return nil, err
	} else {
		s.Title = encrypted
	}
	return s, nil
}

func (s *SystemComplaintModel) DecryptTitle(decryptFunc func(string) (string, error)) (*SystemComplaintModel, error) {
	decrypted, err := decryptFunc(s.Title)
	if err != nil {
		return nil, err
	} else {
		s.Title = decrypted
	}
	return s, nil
}

func (s *SystemComplaintModel) EncryptDetail(encryptFunc func(string) (string, error)) (*SystemComplaintModel, error) {
	encrypted, err := encryptFunc(s.Detail)
	if err != nil {
		return nil, err
	} else {
		s.Detail = encrypted
	}
	return s, nil
}

func (s *SystemComplaintModel) DecryptDetail(decryptFunc func(string) (string, error)) (*SystemComplaintModel, error) {
	decrypted, err := decryptFunc(s.Detail)
	if err != nil {
		return nil, err
	} else {
		s.Detail = decrypted
	}
	return s, nil
}

func (s *SystemComplaintModel) Create() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO 
            SystemComplaint (id, title, detail)
        VALUES
            (:id, :title, :detail)
    `

	if _, err := s.Conn.NamedExecContext(ctx, query, s); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return s.Id, nil
}

func (s *SystemComplaintModel) Count() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT COUNT(*) FROM SystemComplaint"

	count := 0
	err := s.Conn.GetContext(ctx, &count, query)
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return count, nil
}

func (s *SystemComplaintModel) FindById(id string) (*SystemComplaintModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT * FROM SystemComplaint WHERE id = ?"

	if err := s.Conn.GetContext(ctx, s, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return s, nil
}

func (s *SystemComplaintModel) FindAll(params map[string]string) ([]*SystemComplaintModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, title FROM SystemComplaint"

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

	complaints := make([]*SystemComplaintModel, 0)
	if err := s.Conn.SelectContext(ctx, &complaints, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return complaints, nil
}

func (s *SystemComplaintModel) GetPhotos() ([]SystemComplaintImageModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT * FROM SystemComplaintImage WHERE complaintId = ?"

	images := make([]SystemComplaintImageModel, 0)
	if err := s.Conn.SelectContext(ctx, &images, query, s.Id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return images, nil
}
