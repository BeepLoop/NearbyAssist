package e2ee

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/store"
	"time"

	"github.com/jmoiron/sqlx"
)

type MysqlE2EEStore struct {
	db *sqlx.DB
}

func NewMysqlE2EEStore(db *sqlx.DB) *MysqlE2EEStore {
	return &MysqlE2EEStore{
		db: db,
	}
}

func (s *MysqlE2EEStore) NewPublicPem(data *models.PublicKeyModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := "INSERT INTO PublicKey (id, owner, pem) VALUES (:id, :owner, :pem)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlE2EEStore) GetPublicPem(owner string) (*models.PublicKeyModel, error) {
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

func (s *MysqlE2EEStore) NewPrivatePem(data *models.PrivateKeyModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if id, err := store.GenerateNanoId(); err != nil {
		return "", err
	} else {
		data.Id = id
	}

	query := "INSERT INTO PrivateKey (id, owner, pem) VALUES (:id, :owner, :pem)"
	if _, err := s.db.NamedExecContext(ctx, query, data); err != nil {
		return "", nil
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return data.Id, nil
}

func (s *MysqlE2EEStore) GetPrivatePem(owner string) (*models.PrivateKeyModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	privateKey := new(models.PrivateKeyModel)
	query := "SELECT id, owner, pem FROM PrivateKey WHERE owner = ?"
	if err := s.db.GetContext(ctx, privateKey, query, owner); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return privateKey, nil
}
