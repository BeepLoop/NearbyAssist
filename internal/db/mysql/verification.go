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

func (m *Mysql) NewIdentityVerification(model *models.IdentityVerificationModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := `
        INSERT INTO IdentityVerification (name, address, idType, idNumber, frontId, backId, face)
        VALUES ( :name, :address, :idType, :idNumber, :frontId, :backId, :face)
    `

	res, err := m.Conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) FindIdentityVerificationById(id int) (*models.IdentityVerificationModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, name, address, idType, idNumber, frontId, backId, face FROM IdentityVerification WHERE id = ?"

	model := &models.IdentityVerificationModel{}
	if err := m.Conn.GetContext(ctx, model, query, id); err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return nil, nil
}

func (m *Mysql) NewFrontId(model *models.FrontIdModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO FrontId (url) VALUES (:url)"

	res, err := m.Conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) NewBackId(model *models.BackIdModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO BackId (url) VALUES (:url)"

	res, err := m.Conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}

func (m *Mysql) NewFace(model *models.FaceModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Face (url) VALUES (:url)"

	res, err := m.Conn.NamedExecContext(ctx, query, model)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(id), nil
}
