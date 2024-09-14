package management

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlManagementStore struct {
	db *sqlx.DB
}

func NewMysqlManagementStore(db *sqlx.DB) *MysqlManagementStore {
	return &MysqlManagementStore{db: db}
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

func (s *MysqlManagementStore) GetApplications(filter map[string]string) ([]*models.ApplicationModel, error) {
	return nil, nil
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

func (s *MysqlManagementStore) GetVerificationRequests(filter map[string]string) ([]*models.IdentityVerificationModel, error) {
	return nil, nil
}

func (s *MysqlManagementStore) GetVerificationRequestById(id string) (*models.IdentityVerificationModel, error) {
	return nil, nil
}
