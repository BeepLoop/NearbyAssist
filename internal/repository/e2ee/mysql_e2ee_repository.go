package e2ee_repo

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlE2EERepository struct {
	db *sqlx.DB
}

func NewMysqlE2EERepository(db *sqlx.DB) *MysqlE2EERepository {
	return &MysqlE2EERepository{
		db: db,
	}
}

func (s *MysqlE2EERepository) NewPublicPem(data *models.PublicKeyModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data.Id = utils.GenerateId()

	query := "INSERT INTO PublicKey (id, owner, pem) VALUES (:id, :owner, :pem)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlE2EERepository) GetPublicPem(owner string) (*models.PublicKeyModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publicKey := new(models.PublicKeyModel)
	query := "SELECT id, owner, pem FROM PublicKey WHERE owner = ?"
	if err := s.db.GetContext(ctx, publicKey, query, owner); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return publicKey, nil
}

func (s *MysqlE2EERepository) NewPrivatePem(data *models.PrivateKeyModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data.Id = utils.GenerateId()

	query := "INSERT INTO PrivateKey (id, owner, pem) VALUES (:id, :owner, :pem)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlE2EERepository) GetPrivatePem(owner string) (*models.PrivateKeyModel, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existsQuery := `
        SELECT
            CASE
                WHEN EXISTS (SELECT 1 FROM PrivateKey WHERE owner = ?)
                THEN 1
                ELSE 0
            END AS user_exists;
    `
	exists := false
	if err := s.db.GetContext(ctx, &exists, existsQuery, owner); err != nil {
		return nil, false, err
	}

	if !exists {
		return nil, exists, nil
	}

	privateKey := new(models.PrivateKeyModel)
	query := "SELECT id, owner, pem FROM PrivateKey WHERE owner = ?"
	if err := s.db.GetContext(ctx, privateKey, query, owner); err != nil {
		return nil, false, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, false, context.DeadlineExceeded
	}

	return privateKey, exists, nil
}
