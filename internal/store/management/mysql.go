package management

import (
	"context"
	"errors"
	"fmt"
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"nearbyassist/internal/store"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlManagementStore struct {
	DEFAULT_LIMIT  int
	DEFAULT_OFFSET int
	db             *sqlx.DB
}

func NewMysqlManagementStore(db *sqlx.DB) *MysqlManagementStore {
	return &MysqlManagementStore{
		DEFAULT_LIMIT:  10,
		DEFAULT_OFFSET: 0,
		db:             db,
	}
}

func (s *MysqlManagementStore) CreateStaff(data *models.AdminModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := "INSERT INTO Admin (id, username, password, usernameHash, role) VALUES (:id, :username, :password, :usernameHash, :role)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlManagementStore) GetUserById(id string) (*models.UserModel, error) {
	return nil, nil
}

func (s *MysqlManagementStore) RestrictVendor(id string) error {
	return nil
}

func (s *MysqlManagementStore) UnrestrictVendor(id string) error {
	return nil
}

func (s *MysqlManagementStore) GetApplications(params map[string]string) ([]*models.ApplicationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Initial query
	query := "SELECT id, applicantId, status, createdAt FROM Application"

	// Apply filters
	if filter, ok := params["status"]; ok {
		switch filter {
		case "all":
			query += " WHERE status = 'approved' OR status = 'rejected' OR status = 'pending'"
		case "approved":
			query += " WHERE status = 'approved'"
		case "rejected":
			query += " WHERE status = 'rejected'"
		case "pending":
			query += " WHERE status = 'pending'"
		default:
			return nil, errors.New("Invalid parameter found")
		}
	} else {
		query += " WHERE status = 'approved' OR status = 'rejected' OR status = 'pending'"
	}

	// Order by id and createdAt (deterministic order)
	query += " ORDER BY id, createdAt"

	// Paginate using limit and offset
	if page, ok := params["page"]; ok {
		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			return nil, err
		}

		pageSize := s.DEFAULT_LIMIT
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
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", s.DEFAULT_LIMIT, s.DEFAULT_OFFSET)
	}

	// Execute query
	applications := make([]*models.ApplicationModel, 0)
	if err := s.db.SelectContext(ctx, &applications, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return applications, nil
}

func (s *MysqlManagementStore) GetApplicationById(id string) (*response.ApplicationDetail, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	application := new(response.ApplicationDetail)
	query := `
        SELECT 
            a.id AS id,
            a.applicantId AS applicantId,
            a.job AS job,
            a.status AS status,
            a.createdAt AS createdAt,
            p.url AS proofUrl
        FROM 
            Application a
            JOIN ApplicationProof p ON a.id = p.applicationId
        WHERE 
            a.id = ?
    `
	if err := s.db.GetContext(ctx, application, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return application, nil
}

func (s *MysqlManagementStore) ApproveApplication(id string) error {
	return nil
}

func (s *MysqlManagementStore) RejectApplication(id string) error {
	return nil
}

func (s *MysqlManagementStore) GetTransaction(id string) (*models.TransactionModel, error) {
	return nil, nil
}

func (s *MysqlManagementStore) GetSystemComplaints(filter map[string]string) ([]*models.ComplaintModel, error) {
	return nil, nil
}

func (s *MysqlManagementStore) GetSystemComplaintById(id string) (*models.ComplaintModel, error) {
	return nil, nil
}

func (s *MysqlManagementStore) ApproveIdentityVerification(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	updateStatus := "UPDATE IdentityVerification SET status = 'approved' WHERE id = ?"
	if _, err := tx.ExecContext(ctx, updateStatus, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	updateUser := "UPDATE User SET verified = 1 WHERE id = (SELECT user FROM IdentityVerification WHERE id = ?)"
	if _, err := tx.ExecContext(ctx, updateUser, id); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return context.DeadlineExceeded
	}

	return nil
}

func (s *MysqlManagementStore) GetVerificationRequests(filter map[string]string) ([]*response.IdentityVerification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	requests := make([]*response.IdentityVerification, 0)
	query := `
        SELECT 
            id, user as userId, status, createdAt
        FROM 
            IdentityVerification
    `
	if filter, ok := filter["status"]; ok {
		switch filter {
		case "all":
			query += " WHERE status = 'approved' OR status = 'rejected' OR status = 'pending'"
		case "approved":
			query += " WHERE status = 'approved'"
		case "rejected":
			query += " WHERE status = 'rejected'"
		case "pending":
			query += " WHERE status = 'pending'"
		default:
			return nil, errors.New("Invalid parameter found")
		}
	} else {
		query += " WHERE status = 'approved' OR status = 'rejected' OR status = 'pending'"
	}

	if err := s.db.SelectContext(ctx, &requests, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return requests, nil
}

func (s *MysqlManagementStore) GetVerificationRequestById(id string) (*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	request := new(models.IdentityVerificationModel)
	query := `
        SELECT 
            id,
            user as userId,
            name,
            address,
            idType,
            idNumber,
            status,
            frontIdImageUrl,
            backIdImageUrl,
            faceImageUrl,
            createdAt
        FROM 
            IdentityVerification
        WHERE
            id = ?
    `
	if err := s.db.GetContext(ctx, request, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return request, nil
}
