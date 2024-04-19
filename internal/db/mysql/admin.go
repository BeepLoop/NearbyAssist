package mysql

import (
	"context"
	"nearbyassist/internal/models"
	"time"
)

func (m *Mysql) FindAdminByUsername(username string) (*models.AdminModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	query := "SELECT id, username, password, role FROM Admin WHERE username = ?"

	admin := models.NewAdminModel()
	err := m.Conn.GetContext(ctx, admin, query, username)
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
