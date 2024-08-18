package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"nearbyassist/internal/response"
	"time"
)

func (m *Mysql) FindAllIdentityVerification() ([]response.AllVerification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, user, createdAt FROM IdentityVerification"

	requests := make([]response.AllVerification, 0)
	if err := m.Conn.SelectContext(ctx, &requests, query); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return requests, nil
}

func (m *Mysql) NewIdentityVerification(model *models.IdentityVerificationModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO IdentityVerification (id, name, address, idType, idNumber, frontId, backId, face)
        VALUES ( :id, :name, :address, :idType, :idNumber, :frontId, :backId, :face)
    `

	if _, err := m.Conn.NamedExecContext(ctx, query, model); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return model.Id, nil
}

func (m *Mysql) FindIdentityVerificationById(identityVerificationId string) (*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, address, idType, idNumber, frontId, backId, face FROM IdentityVerification WHERE id = ?"

	model := &models.IdentityVerificationModel{}
	if err := m.Conn.GetContext(ctx, model, query, identityVerificationId); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return nil, nil
}

func (m *Mysql) NewFrontId(model *models.FrontIdModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO FrontId (id, url) VALUES (:id, :url)"

	if _, err := m.Conn.NamedExecContext(ctx, query, model); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return model.Id, nil
}

func (m *Mysql) NewBackId(model *models.BackIdModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO BackId (id, url) VALUES (:id, :url)"

	if _, err := m.Conn.NamedExecContext(ctx, query, model); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return model.Id, nil
}

func (m *Mysql) NewFace(model *models.FaceModel) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Face (id, url) VALUES (:id, :url)"

	if _, err := m.Conn.NamedExecContext(ctx, query, model); err != nil {
		return "", err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return "", context.DeadlineExceeded
	}

	return model.Id, nil
}
