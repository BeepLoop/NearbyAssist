package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) FindAdminByUsernameHash(hash string) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, username, password, role FROM Admin WHERE usernameHash = ?"

	admin := models.NewAdminModel()
	err := m.Conn.GetContext(ctx, admin, query, hash)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}

func (m *Mysql) FindAdminById(id int) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, username, password, role FROM Admin WHERE id = ?"

	admin := models.NewAdminModel()
	err := m.Conn.GetContext(ctx, admin, query, id)
	if err != nil {
		return nil, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return nil, context.DeadlineExceeded
	}

	return admin, nil
}

func (m *Mysql) NewAdmin(admin *models.AdminModel) (int, error) {

	return 0, nil
}

func (m *Mysql) NewStaff(staff *models.AdminModel) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "INSERT INTO Admin (username, password, usernameHash) VALUES (:username, :password, :usernameHash)"

	res, err := m.Conn.NamedExecContext(ctx, query, staff)
	if err != nil {
		return 0, err
	}

	insertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if ctx.Err() == context.DeadlineExceeded {
		return 0, context.DeadlineExceeded
	}

	return int(insertId), nil
}
